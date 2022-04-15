package ArenaService

import (
	"fmt"
	"github.com/aicam/cryptowow_back/database"
	"strconv"
	"time"
)

func (s *Service) InviteOperation(inviter, invited int, invitedName string, betAmount float32) {
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

func (s *Service) AcceptInvitation(inviter, invited int, inviterName string) error {
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

func (s *Service) StartGame(bucketID uint) error {
	err := s.Redis.Set(s.Context, strconv.Itoa(int(bucketID)), "00", time.Duration(READYCHECKCOUNTER)*time.Second).Err()
	if err != nil {
		return err
	}
	// TODO: uncomment.prometheus
	//s.PP.Gauges["bet_system_accept_operation_in_progress"].Add(-1)
	return nil
}

func (s *Service) StartGameAcceptHandler(bucketID uint, acceptedID int) (error, int) {
	stat, err := s.Redis.Get(s.Context, strconv.Itoa(int(bucketID))).Result()
	if err != nil {
		return err, 0
	}
	var bucketTeam database.TeamReadyGames
	err = s.DB.Where("id = " + strconv.Itoa(int(bucketID))).First(&bucketTeam).Error
	if err != nil {
		return err, -1
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
		return nil, 1
	}
	return nil, 0
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
