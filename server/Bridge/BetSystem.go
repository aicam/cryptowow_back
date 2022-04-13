package Bridge

import (
	"fmt"
	"github.com/aicam/cryptowow_back/database"
	"strconv"
	"time"
)

func (s *Server) InviteOperation(inviter, invited int, invitedName string, betAmount float32) {
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
}

func (s *Server) AcceptInvitation(inviter, invited int, inviterName string) error {
	var queued database.QueuedTeams
	if err := s.DB.Where(&database.QueuedTeams{TeamId: invited,
		InQueueTeamId: inviter}).First(&queued).Error; err != nil {
		return err
	}
	s.DB.Delete(&queued)

	var request database.TeamRequests
	if err := s.DB.Where(&database.TeamRequests{TeamId: inviter,
		RequestedTeamId: invited}).First(&queued).Error; err != nil {
		return err
	}
	s.DB.Delete(&request)

	s.DB.Save(&database.TeamReadyGames{
		TeamId:     inviter,
		OpponentId: invited,
	}).Save(database.TeamReadyGames{
		TeamId:     invited,
		OpponentId: inviter,
	})

	notif := database.BetNotification{
		TeamId:    inviter,
		Title:     "Your request accepted",
		Body:      fmt.Sprintf("Team %s accepted your battle request", inviterName),
		Seen:      false,
		NotifType: 0,
	}
	s.DB.Save(&notif)

	return nil
}

func (s *Server) StartGame(inviter, invited int) error {
	err := s.Redis.Set(s.Context, strconv.Itoa(inviter), strconv.Itoa(invited), 40*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}
