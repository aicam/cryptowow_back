package ArenaService

import (
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strconv"
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

func getArenaTeamById(DB *gorm.DB, arenaTeamID uint) ArenaTeam {
	var arenaTeam ArenaTeam
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT * from arena_team WHERE arenaTeamId = " + strconv.Itoa(int(arenaTeamID))).Scan(&arenaTeam)
	return arenaTeam
}

func getHeroLeaderNameByArenaId(DB *gorm.DB, teamID uint) string {
	var HeroNameS struct {
		HeroName string `gorm:"column:name"`
	}
	DB.Clauses(dbresolver.Use("characters")).
		Raw("SELECT `name`, `guid` FROM `characters` WHERE " +
			"`guid` = (SELECT `captainGuid` FROM `arena_team` WHERE `arenaTeamId` = " + strconv.Itoa(int(teamID)) + ")").
		First(&HeroNameS)
	return HeroNameS.HeroName
}
