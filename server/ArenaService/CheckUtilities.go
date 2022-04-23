package ArenaService

import (
	"context"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"strconv"
	"strings"
)

func CheckArenaTeamUserAccount(DB *gorm.DB, teamID int, username string) bool {
	usernameByID := getUsernameByArenaTeamID(DB, teamID)
	if strings.ToUpper(username) != strings.ToUpper(usernameByID) {
		return false
	}
	return true
}

func CheckBalance(DB *gorm.DB, username string, value float64, currency string) bool {
	amount := WalletService.GetAccountBalance(username, currency, DB)
	if amount < value {
		return false
	}
	return true
}

func CheckQueueTeam(DB *gorm.DB, inviter, invited int) error {
	var queueTeam database.QueuedTeams
	err := DB.Where(database.QueuedTeams{
		TeamId:        invited,
		InQueueTeamId: inviter,
	}).First(&queueTeam).Error

	if err != nil {
		return err
	}
	return nil
}

func CheckTeamRequest(DB *gorm.DB, inviter, invited int) error {
	var requestTeam database.TeamRequests
	err := DB.Where(database.TeamRequests{
		TeamId:          inviter,
		RequestedTeamId: invited,
	}).First(&requestTeam).Error

	if err != nil {
		return err
	}
	return nil
}

func CheckTeamReady(DB *gorm.DB, inviter, invited int) (database.TeamReadyGames, error) {
	var requestTeam database.TeamReadyGames
	err := DB.Where(database.TeamReadyGames{
		InviterTeam: inviter,
		InvitedTeam: invited,
	}).First(&requestTeam).Error

	if err != nil {
		return database.TeamReadyGames{}, err
	}
	return requestTeam, nil
}

func CheckSameArenaType(DB *gorm.DB, inviter, invited int) uint8 {
	var arenaTypeDBStruct struct {
		ArenaType uint8 `gorm:"column:type"`
	}
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT `type` from arena_team WHERE arenaTeamId = " + strconv.Itoa(inviter)).First(&arenaTypeDBStruct)
	typeAlliance := arenaTypeDBStruct.ArenaType
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT `type` from arena_team WHERE arenaTeamId = " + strconv.Itoa(invited)).First(&arenaTypeDBStruct)
	typeHorde := arenaTypeDBStruct.ArenaType
	log.Println(typeAlliance, " ", typeHorde)
	if typeAlliance != typeHorde {
		return 0
	}

	return typeAlliance
}

func CheckIsAlreadyStarted(DB *gorm.DB, rdb *redis.Client, ctx context.Context, teamID uint) bool {
	var buckets []database.TeamReadyGames
	DB.Where("inviter_team = ? OR invited_team = ?", teamID, teamID).Find(&buckets)
	for _, bucket := range buckets {
		err := rdb.Get(ctx, strconv.Itoa(int(bucket.ID))).Err()
		if err != redis.Nil {
			return false
		}
	}
	return true
}

func CheckAlreadyInArena(DB *gorm.DB, teamID uint) bool {
	var buckets []database.TeamReadyGames
	DB.Where("inviter_team = ? OR invited_team = ?", teamID, teamID).Find(&buckets)
	var gameData database.InGameTeamData
	for _, bucket := range buckets {
		if err := DB.Where(&database.InGameTeamData{BucketID: bucket.ID}).First(&gameData).Error; err == nil {
			if gameData.Winner == 0 {
				return false
			}
		}
	}
	return true
}
