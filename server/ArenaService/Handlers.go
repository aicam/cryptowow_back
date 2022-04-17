package ArenaService

import (
	"github.com/aicam/cryptowow_back/server/LogService"
	"github.com/gin-gonic/gin"
	"net/http"
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
			c.JSON(http.StatusBadRequest, actionResult(-1, "error in parsing"))
			return
		}

		if reqParams.Inviter == reqParams.Invited {
			c.JSON(http.StatusBadRequest, actionResult(-1, "error in parsing"))
			return
		}

		if !CheckArenaTeamUserAccount(s.DB, reqParams.Inviter, username) {
			c.JSON(http.StatusBadRequest, actionResult(-12, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Username_Check")
			return
		}

		invitedUsername := getUsernameByArenaTeamID(s.DB, reqParams.Invited)
		if invitedUsername == "" {
			c.JSON(http.StatusBadRequest, actionResult(-8, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Invited_Username_Check")
			return
		}

		if !CheckBalance(s.DB, username, reqParams.BetAmount, reqParams.BetCurrency) ||
			!CheckBalance(s.DB, invitedUsername, reqParams.BetAmount, reqParams.BetCurrency) {
			c.JSON(http.StatusBadRequest, actionResult(-8, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Balance_Check")
			return
		}

		_, err = CheckQueueTeam(s.DB, reqParams.Inviter, reqParams.Invited)
		if err == nil {
			c.JSON(http.StatusBadRequest, actionResult(-8, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Duplicate_Check")
			return
		}

		s.InviteOperation(reqParams.Inviter, reqParams.Invited, invitedUsername, reqParams.BetAmount, reqParams.BetCurrency)
		c.JSON(http.StatusOK, actionResult(1, "joined successfully"))
		LogService.LogSucceedJoinOperation(username, invitedUsername)
	}
}

func (s *Service) AcceptInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams AcceptInvitationRequest
		err := c.BindJSON(&reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "parsing error"))
			return
		}

		usernameInvited := c.GetHeader("username")
		if !CheckArenaTeamUserAccount(s.DB, reqParams.Invited, usernameInvited) {
			c.JSON(http.StatusBadRequest, actionResult(-8, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Invitation_Username_Check")
			return
		}

		err = CheckQueueTeam(s.DB, reqParams.Inviter, reqParams.Invited)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-6, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Invitation_Queue_Team_Check")
			return
		}

		err = CheckTeamRequest(s.DB, reqParams.Inviter, reqParams.Invited)
		// should never happen
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-6, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Invitation_Request_Team_Check")
			return
		}

		err = s.AcceptInvitationOperation(reqParams.Inviter, reqParams.Invited, getUsernameByArenaTeamID(s.DB, reqParams.Inviter))
		if err != nil {
			c.JSON(http.StatusBadGateway, actionResult(-1, "Service unavailable"))
			LogService.LogCrash("Redis", "ArenaService_Accept_Invitation")
			return
		}
		c.JSON(http.StatusOK, actionResult(1, "Accepted successfully!"))

	}
}
