package Bridge

import (
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
)

type BetEvents struct {
	EventType uint8
	DB        *gorm.DB
}

func (b *BetEvents) InviteEvent(TeamIdA, TeamIdB int) {
	newQueued := database.QueuedTeams{
		TeamId:        TeamIdB,
		InQueueTeamId: TeamIdA,
	}
	b.DB.Save(&newQueued)

}
