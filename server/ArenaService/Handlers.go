package ArenaService

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/LogService"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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

		err = CheckQueueTeam(s.DB, reqParams.Inviter, reqParams.Invited)
		if err == nil {
			c.JSON(http.StatusBadRequest, actionResult(-8, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Duplicate_Check")
			return
		}

		arenaType := CheckSameArenaType(s.DB, reqParams.Inviter, reqParams.Invited)
		if arenaType == 0 {
			c.JSON(http.StatusBadRequest, actionResult(-20, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Invite_Same_Type_Check")
			return
		}

		s.InviteOperation(reqParams.Inviter, reqParams.Invited, invitedUsername, reqParams.BetAmount, reqParams.BetCurrency)
		c.JSON(http.StatusOK, actionResult(1, "joined successfully"))
		LogService.LogSucceedJoinOperation(username, invitedUsername)
	}
}

func (s *Service) AcceptInvitation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams GeneralInvitationRequest
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
			LogService.LogCrash("MySql", "ArenaService_Accept_Invitation")
			return
		}

		c.JSON(http.StatusOK, actionResult(1, "Accepted successfully!"))
	}
}

func (s *Service) StartGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams GeneralInvitationRequest
		err := c.BindJSON(&reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "parsing error"))
			return
		}

		username := strings.ToUpper(c.GetHeader("username"))
		usernameInviter := getUsernameByArenaTeamID(s.DB, reqParams.Inviter)
		usernameInvited := getUsernameByArenaTeamID(s.DB, reqParams.Invited)
		if username != usernameInvited && username != usernameInviter {
			c.JSON(http.StatusBadRequest, actionResult(-6, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Start_Game_Username_Check")
			return
		}

		readyBucket, err := CheckTeamReady(s.DB, reqParams.Inviter, reqParams.Invited)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-6, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Start_Game_Exist_Check")
			return
		}

		err = s.StartGameOperation(readyBucket.ID)
		if err != nil {
			c.JSON(http.StatusBadGateway, actionResult(-1, "Service unavailable"))
			LogService.LogCrash("Redis", "ArenaService_Accept_Invitation")
			return
		}

		c.JSON(http.StatusOK, struct {
			Status   int  `json:"status"`
			BucketId uint `json:"bucket_id"`
		}{Status: 1, BucketId: readyBucket.ID})
	}
}

func (s *Service) AcceptStartGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams AcceptStartGameRequest
		err := c.BindJSON(&reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "parsing error"))
			return
		}

		username := strings.ToUpper(c.GetHeader("username"))

		if username != getUsernameByArenaTeamID(s.DB, reqParams.TeamID) {
			c.JSON(http.StatusBadRequest, actionResult(-6, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Start_Game_Username_Check")
			return
		}
		var readyTeam database.TeamReadyGames
		readyTeam, err = CheckTeamReady(s.DB, reqParams.TeamID, reqParams.OpponentID)
		if err != nil {
			readyTeam, err = CheckTeamReady(s.DB, reqParams.OpponentID, reqParams.TeamID)
			if err != nil {
				c.JSON(http.StatusBadRequest, actionResult(-10, "error in parsing"))
				LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Start_Game_Bucket_Check")
				return
			}
		}

		if readyTeam.IsPlayed {
			c.JSON(http.StatusBadRequest, actionResult(-21, "error in parsing"))
			LogService.LogPotentialCyberAttack(c, "ArenaService_Accept_Start_Played_Bucket_Check")
			return
		}

		err = s.StartGameAcceptHandler(s.DB, readyTeam.ID, reqParams.TeamID)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "Time is over!"))
			return
		}

		c.JSON(http.StatusOK, actionResult(1, "Accepted successfully! Be Ready...."))
	}
}

func (s *Service) GetGameStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var invReq struct {
			BucketID uint `json:"bucket_id"`
		}
		err := c.BindJSON(&invReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "Wrong parameters"))
		}
		gameStat := s.IsStarted(invReq.BucketID)
		if gameStat == -1 {
			c.JSON(http.StatusOK, GameStatResponse{
				Status: 1,
				Code:   -1,
				Result: "Time is over!",
			})
		}
		if gameStat == 0 {
			c.JSON(http.StatusOK, GameStatResponse{
				Status: 1,
				Code:   0,
				Result: "Waiting for opponent",
			})
		}
		if gameStat == 1 {
			c.JSON(http.StatusOK, GameStatResponse{
				Status: 1,
				Code:   1,
				Result: "Joined Successfully!",
			})
		}

	}
}

func (s *Service) GetResult() gin.HandlerFunc {
	return func(c *gin.Context) {
		teamId := c.Param("team_id")
		var inGameInfos []database.InGameTeamData
		s.DB.Raw("SELECT * FROM `in_game_team_data` WHERE `in_game_team_data`.`bucket_id` IN " +
			"(SELECT `team_ready_games`.`id` FROM `team_ready_games` WHERE " +
			"`team_ready_games`.`inviter_team` = " + teamId + " OR " +
			"`team_ready_games`.`invited_team` = " + teamId + ")").Find(&inGameInfos)

		for _, inGameInfo := range inGameInfos {
			if inGameInfo.Winner == 0 {
				s.ProcessGame(inGameInfo.BucketID)
			}
		}
		c.JSON(http.StatusOK, struct {
			Status       int                       `json:"status"`
			AllGamesInfo []database.InGameTeamData `json:"all_games_info"`
		}{Status: 1, AllGamesInfo: inGameInfos})
	}
}
