package ArenaService

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strconv"
	"strings"
)

func getUsernameByArenaTeamID(DB *gorm.DB, teamID int) string {
	var accountID struct {
		ID int `gorm:"column:account"`
	}
	err := DB.Clauses(dbresolver.Use("characters")).Raw(
		"SELECT `characters`.`account` FROM `characters` WHERE `characters`.`guid` = " +
			"(SELECT `arena_team`.`captainGuid` FROM `arena_team` " +
			"WHERE `arena_team`.`arenaTeamId` = " + strconv.Itoa(teamID) + ")").
		First(&accountID).Error
	if err != nil {
		return ""
	}

	var checkUsername struct {
		Username string `gorm:"column:username"`
	}
	err = DB.Clauses(dbresolver.Use("auth")).Raw("SELECT username FROM account WHERE id=" + strconv.Itoa(accountID.ID)).
		First(&checkUsername).Error
	// should never happen
	if err != nil {
		return ""
	}
	return checkUsername.Username
}

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
