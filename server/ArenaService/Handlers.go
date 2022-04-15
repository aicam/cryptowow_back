package ArenaService

import (
	"github.com/gin-gonic/gin"
	"gorm.io/plugin/dbresolver"
	"net/http"
	"strconv"
)

func actionResult(statusCode int, body string) struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
} {
	return struct {
		Status int    `json:"status"`
		Body   string `json:"body"`
	}{Status: statusCode, Body: body}
}

func (s *Service) InviteTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var reqParams InviteRequest
		err := c.BindJSON(&reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(http.StatusBadRequest, "error in parsing"))
		}

		// find account id based on inviter arena team id
		var accountID struct {
			ID int `gorm:"column:account"`
		}
		err = s.DB.Clauses(dbresolver.Use("characters")).Raw(
			"SELECT `characters`.`account` FROM `characters` WHERE `characters`.`guid` = " +
				"(SELECT `arena_team`.`captainGuid` FROM `arena_team` " +
				"WHERE `arena_team`.`arenaTeamId` = " + strconv.Itoa(reqParams.Inviter) + ")").
			First(&accountID).Error
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(http.StatusBadRequest, "error in parsing"))
		}

		var checkUsername struct {
			Username string `gorm:"column:username"`
		}
		err = s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT username FROM account WHERE id=" + strconv.Itoa(accountID.ID)).
			First(&checkUsername).Error
		// should never happen
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(http.StatusBadRequest, "error in parsing"))
		}
	}
}
