package ArenaService

import (
	"github.com/aicam/cryptowow_back/database"
	"os"
	"strconv"
	"time"
)

func (s *Service) InviteOperation(inviter, invited uint, invitedName string, betAmount float64, currency string) {
	inviterArenaTeam := getArenaTeamById(s.DB, inviter)
	invitedArenaTeam := getArenaTeamById(s.DB, invited)
	s.DB.Save(&database.BetInfo{
		InviterTeam:        inviter,
		InvitedTeam:        invited,
		Amount:             betAmount,
		Currency:           currency,
		InviterUsername:    getUsernameByArenaTeamID(s.DB, inviter),
		InvitedUsername:    getUsernameByArenaTeamID(s.DB, invited),
		Step:               0,
		InviterSeasonGames: inviterArenaTeam.SeasonGames,
		InviterSeasonWins:  inviterArenaTeam.SeasonWins,
		InvitedSeasonGames: invitedArenaTeam.SeasonGames,
		InvitedSeasonWins:  invitedArenaTeam.SeasonWins,
		ArenaType:          inviterArenaTeam.ArenaType,
		Winner:             0,
	})
	// prometheus
	s.PP.Counters["bet_system_invite_operation_counter"].Inc()
	s.PP.Gauges["bet_system_invite_operation_in_progress"].Inc()
}

func (s *Service) AcceptInvitationOperation(inviter, invited uint) error {
	var betInfo database.BetInfo
	err := s.DB.Where(&database.BetInfo{
		InviterTeam: inviter,
		InvitedTeam: invited,
		Step:        0,
	}).First(&betInfo).Error
	betInfo.Step = 1
	s.DB.Save(&betInfo)

	// prometheus
	s.PP.Counters["bet_system_accept_operation_counter"].Inc()
	s.PP.Gauges["bet_system_accept_operation_in_progress"].Inc()
	s.PP.Gauges["bet_system_invite_operation_in_progress"].Add(-1)
	return err
}

func (s *Service) StartGameOperation(bucketID uint) error {
	err := s.Rdb.Set(s.Context, strconv.Itoa(int(bucketID)), "00", time.Duration(READYCHECKCOUNTER)*time.Second).Err()
	if err != nil {
		return err
	}
	// prometheus
	s.PP.Counters["bet_system_start_game_operation_counter"].Inc()
	s.PP.Gauges["bet_system_accept_operation_in_progress"].Add(-1)
	return nil
}

func (s *Service) StartGameAcceptHandler(bucketID uint, acceptedID uint) error {
	stat, err := s.Rdb.Get(s.Context, strconv.Itoa(int(bucketID))).Result()
	if err != nil {
		return err
	}
	var betInfo database.BetInfo
	err = s.DB.Where("id = " + strconv.Itoa(int(bucketID))).First(&betInfo).Error
	if err != nil {
		return err
	}
	var acceptedTeam uint8
	if acceptedTeam = 1; acceptedID == betInfo.InviterTeam {
		acceptedTeam = 0
	}
	var newStat string
	if newStat = stat[:1] + "1"; acceptedTeam == 0 {
		newStat = "1" + stat[1:]
	}

	s.Rdb.Set(s.Context, strconv.Itoa(int(bucketID)), newStat, time.Duration(READYCHECKCOUNTER)*time.Second)
	// check result
	if newStat == "11" {
		AppendNewGame(NewGameParams{
			ArenaFilePath: os.Getenv("ARENAFILEPATH"),
			TeamID1:       betInfo.InviterTeam,
			TeamID2:       betInfo.InvitedTeam,
			LeaderName1:   getHeroLeaderNameByArenaId(s.DB, betInfo.InviterTeam),
			LeaderName2:   getHeroLeaderNameByArenaId(s.DB, betInfo.InvitedTeam),
			ArenaType:     strconv.Itoa(int(betInfo.ArenaType)),
		})
		betInfo.Step = GameStarted
		s.DB.Save(&betInfo)

		// prometheus
		s.PP.Counters["bet_system_match_counter"].Inc()
		s.PP.Gauges["bet_system_match_in_progress"].Inc()
	}
	return nil
}

func (s *Service) IsStarted(bucketID uint) int {
	stat, err := s.Rdb.Get(s.Context, strconv.Itoa(int(bucketID))).Result()
	if err != nil {
		return -1
	}
	if stat == "11" {
		return 1
	}
	return 0
}

func (s *Service) ProcessGame(betInfo database.BetInfo) database.BetInfo {
	inviterArenaTeam := getArenaTeamById(s.DB, betInfo.InviterTeam)
	invitedArenaTeam := getArenaTeamById(s.DB, betInfo.InvitedTeam)
	winnerId := uint(0)
	if inviterArenaTeam.SeasonGames == betInfo.InviterSeasonGames {
		return betInfo
	} else {
		if inviterArenaTeam.SeasonWins > betInfo.InviterSeasonWins &&
			invitedArenaTeam.SeasonWins == betInfo.InvitedSeasonWins {
			winnerId = betInfo.InviterTeam
		} else if invitedArenaTeam.SeasonWins > betInfo.InvitedSeasonWins &&
			inviterArenaTeam.SeasonWins == betInfo.InviterSeasonWins {
			winnerId = uint(betInfo.InvitedTeam)
		}
	}
	if winnerId != 0 {
		betInfo.Winner = winnerId
		betInfo.Step = GameFinished
		s.DB.Save(&betInfo)
		// prometheus
		s.PP.Counters["bet_system_match_finished"].Inc()
		s.PP.Gauges["bet_system_match_in_progress"].Add(-1)
	}
	return betInfo
}
