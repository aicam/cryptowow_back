package ArenaService

import (
	"errors"
	"fmt"
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func (s *Service) InviteOperation(inviter, invited int, invitedName string, betAmount float64, currency string) {
	newQueued := database.QueuedTeams{
		TeamId:        invited,
		InQueueTeamId: inviter,
	}
	notif := database.BetNotification{
		TeamId:    invited,
		Title:     "New match request",
		Body:      fmt.Sprintf("Team %s invited you to a %f bet!", invitedName, betAmount),
		Seen:      false,
		NotifType: 0,
	}
	newRequest := database.TeamRequests{
		TeamId:          inviter,
		RequestedTeamId: invited,
	}
	s.DB.Save(&newRequest)
	s.DB.Save(&notif)
	s.DB.Save(&newQueued)

	// TODO: uncomment.prometheus
	//s.PP.Counters["bet_system_invite_operation_counter"].Inc()
	//s.PP.Gauges["bet_system_invite_operation_in_progress"].Inc()
}

func (s *Service) AcceptInvitationOperation(inviter, invited int, inviterName string) error {
	var queued database.QueuedTeams
	if err := s.DB.Where(&database.QueuedTeams{TeamId: invited,
		InQueueTeamId: inviter}).First(&queued).Error; err != nil {
		return err
	}
	s.DB.Delete(&queued)

	var request database.TeamRequests
	if err := s.DB.Where(&database.TeamRequests{TeamId: inviter,
		RequestedTeamId: invited}).First(&request).Error; err != nil {
		return err
	}
	s.DB.Delete(&request)

	s.DB.Save(&database.TeamReadyGames{
		InviterTeam: inviter,
		InvitedTeam: invited,
	})

	notif := database.BetNotification{
		TeamId:    inviter,
		Title:     "Your request accepted",
		Body:      fmt.Sprintf("Team %s accepted your battle request", inviterName),
		Seen:      false,
		NotifType: 0,
	}
	s.DB.Save(&notif)
	// TODO: uncomment.prometheus
	//s.PP.Counters["bet_system_accept_operation_counter"].Inc()
	//s.PP.Gauges["bet_system_accept_operation_in_progress"].Inc()
	//s.PP.Gauges["bet_system_invite_operation_in_progress"].Add(-1)
	return nil
}

func (s *Service) StartGameOperation(bucketID uint) error {
	err := s.Redis.Set(s.Context, strconv.Itoa(int(bucketID)), "00", time.Duration(READYCHECKCOUNTER)*time.Second).Err()
	if err != nil {
		return err
	}
	// TODO: uncomment.prometheus
	//s.PP.Gauges["bet_system_accept_operation_in_progress"].Add(-1)
	return nil
}

func (s *Service) StartGameAcceptHandler(DB *gorm.DB, bucketID uint, acceptedID int) error {
	stat, err := s.Redis.Get(s.Context, strconv.Itoa(int(bucketID))).Result()
	if err != nil {
		return err
	}
	var bucketTeam database.TeamReadyGames
	err = s.DB.Where("id = " + strconv.Itoa(int(bucketID))).First(&bucketTeam).Error
	if err != nil {
		return err
	}
	var acceptedTeam uint8
	if acceptedTeam = 1; acceptedID == bucketTeam.InviterTeam {
		acceptedTeam = 0
	}
	var newStat string
	if newStat = stat[:1] + "1"; acceptedTeam == 0 {
		newStat = "1" + stat[1:]
	}

	s.Redis.Set(s.Context, strconv.Itoa(int(bucketID)), newStat, time.Duration(READYCHECKCOUNTER)*time.Second)
	// check result
	if newStat == "11" {
		var readyGame database.TeamReadyGames
		DB.Where("id = " + strconv.Itoa(int(bucketID))).First(&readyGame)
		inviterArenaTeam := getArenaTeamById(DB, uint(readyGame.InviterTeam))
		invitedArenaTeam := getArenaTeamById(DB, uint(readyGame.InvitedTeam))
		AppendNewGame(NewGameParams{
			ArenaFilePath: os.Getenv("ARENAFILEPATH"),
			TeamID1:       readyGame.InviterTeam,
			TeamID2:       readyGame.InvitedTeam,
			LeaderName1:   getHeroLeaderNameByArenaId(DB, uint(readyGame.InviterTeam)),
			LeaderName2:   getHeroLeaderNameByArenaId(DB, uint(readyGame.InvitedTeam)),
			ArenaType:     strconv.Itoa(int(inviterArenaTeam.ArenaType)),
		})
		s.DB.Save(&database.InGameTeamData{
			BucketID:           bucketID,
			InviterSeasonGames: uint(inviterArenaTeam.SeasonGames),
			InviterSeasonWins:  uint(inviterArenaTeam.SeasonWins),
			InvitedSeasonGames: uint(invitedArenaTeam.SeasonGames),
			InvitedSeasonWins:  uint(invitedArenaTeam.SeasonWins),
		})
		readyGame.IsPlayed = true
		s.DB.Save(&readyGame)
	}
	return nil
}

func (s *Service) IsStarted(bucketID uint) int {
	stat, err := s.Redis.Get(s.Context, strconv.Itoa(int(bucketID))).Result()
	if err != nil {
		return -1
	}
	if stat == "11" {
		return 1
	}
	return 0
}

func (s *Service) ProcessGame(bucketID uint) GameStatNow {
	var gameStats database.InGameTeamData
	err := s.DB.Where(&database.InGameTeamData{BucketID: bucketID}).First(&gameStats).Error
	if err != nil {
		return GameStatNow{
			IsError:        err,
			IsFinishedPast: false,
			WinnerId:       0,
		}
	}
	if gameStats.Winner != 0 {
		return GameStatNow{
			IsError:        nil,
			IsFinishedPast: true,
			WinnerId:       gameStats.Winner,
		}
	}

	var readyGame database.TeamReadyGames
	s.DB.Where("id = " + strconv.Itoa(int(bucketID))).First(&readyGame)

	inviterArenaTeam := getArenaTeamById(s.DB, uint(readyGame.InviterTeam))
	invitedArenaTeam := getArenaTeamById(s.DB, uint(readyGame.InvitedTeam))

	if inviterArenaTeam.SeasonGames == gameStats.InviterSeasonGames {
		return GameStatNow{
			IsError:        nil,
			IsFinishedPast: false,
			WinnerId:       0,
		}
	} else {
		if inviterArenaTeam.SeasonWins > gameStats.InviterSeasonWins &&
			invitedArenaTeam.SeasonWins == gameStats.InvitedSeasonWins {
			return GameStatNow{
				IsError:        nil,
				IsFinishedPast: false,
				WinnerId:       uint(readyGame.InviterTeam),
			}
		} else if invitedArenaTeam.SeasonWins > gameStats.InvitedSeasonWins &&
			inviterArenaTeam.SeasonWins == gameStats.InviterSeasonWins {
			return GameStatNow{
				IsError:        nil,
				IsFinishedPast: false,
				WinnerId:       uint(readyGame.InvitedTeam),
			}
		}
	}

	return GameStatNow{
		IsError:        errors.New("Match remained unhandled"),
		IsFinishedPast: false,
		WinnerId:       0,
	}
}
