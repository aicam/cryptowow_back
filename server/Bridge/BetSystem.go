package Bridge

import (
	"fmt"
	"github.com/aicam/cryptowow_back/database"
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
