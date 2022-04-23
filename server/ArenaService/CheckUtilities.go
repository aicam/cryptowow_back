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

func CheckArenaTeamUserAccount(DB *gorm.DB, teamID uint, username string) bool {
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

func CheckQueueTeam(DB *gorm.DB, inviter, invited uint) error {
	var queueTeam database.BetInfo
	return DB.Where("(inviter_team = ? AND invited_team = ?) AND (step = 0 OR step = 1)",
		inviter, invited).First(&queueTeam).Error
}

func CheckTeamReady(DB *gorm.DB, inviter, invited uint) (database.BetInfo, error) {
	var betInfo database.BetInfo
	err := DB.Where(&database.BetInfo{
		InviterTeam: inviter,
		InvitedTeam: invited,
		Step:        InvitationAccepted,
	}).First(&betInfo).Error

	if err != nil {
		return database.BetInfo{}, err
	}
	return betInfo, nil
}

func CheckSameArenaType(DB *gorm.DB, inviter, invited uint) uint8 {
	var arenaTypeDBStruct struct {
		ArenaType uint8 `gorm:"column:type"`
	}
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT `type` from arena_team WHERE arenaTeamId = " + strconv.Itoa(int(inviter))).First(&arenaTypeDBStruct)
	typeAlliance := arenaTypeDBStruct.ArenaType
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT `type` from arena_team WHERE arenaTeamId = " + strconv.Itoa(int(invited))).First(&arenaTypeDBStruct)
	typeHorde := arenaTypeDBStruct.ArenaType
	log.Println(typeAlliance, " ", typeHorde)
	if typeAlliance != typeHorde {
		return 0
	}

	return typeAlliance
}

func CheckIsAlreadyStarted(DB *gorm.DB, rdb *redis.Client, ctx context.Context, teamID uint) bool {
	var buckets []database.BetInfo
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
	var bucket database.BetInfo
	if err := DB.Where("(inviter_team = ? OR invited_team = ?) AND step = "+strconv.Itoa(int(GameStarted)),
		teamID, teamID).First(&bucket).Error; err != gorm.ErrRecordNotFound {
		return false
	}
	return true
}
