package WalletService

import (
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strconv"
)

func GetArenaBetTotalDebt(DB *gorm.DB, username string) map[string]float64 {
	var accountID uint
	totalDebt := make(map[string]float64)
	DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id FROM account WHERE username='" + username + "'").Scan(&accountID)

	// get user all arena team ids
	var arenaTeams []uint
	DB.Clauses(dbresolver.Use("characters")).Raw("SELECT `arenaTeamId` FROM `arena_team` WHERE `captainGuid` " +
		"IN (SELECT `guid` FROM `characters` WHERE `account` = " + strconv.Itoa(int(accountID)) + ")").
		Find(&arenaTeams)

	for _, arenaTeam := range arenaTeams {
		var betInfos []database.BetInfo
		DB.Where("(inviter_team = ? OR invited_team = ?) AND ( step != 3 )", arenaTeam, arenaTeam).Find(&betInfos)
		for _, betInfo := range betInfos {
			if _, ok := totalDebt[betInfo.Currency]; ok {
				totalDebt[betInfo.Currency] += betInfo.Amount
			} else {
				totalDebt[betInfo.Currency] = betInfo.Amount
			}
		}
	}

	return totalDebt
}
