package ArenaService

type ArenaTeam struct {
	ArenaTeamID int    `json:"arena_team_id" gorm:"column:arenaTeamId"`
	LeaderName  string `json:"leader_name" gorm:"column:name"`
	LeaderGUID  int    `json:"leader_guid" gorm:"column:captainGuid"`
	ArenaType   uint8  `json:"arena_type" gorm:"column:type"`
	SeasonGames uint   `json:"season_games" gorm:"column:seasonGames"`
	SeasonWins  uint   `json:"season_wins" gorm:"column:seasonWins"`
}

type InviteRequest struct {
	Inviter     uint    `json:"inviter"`
	Invited     uint    `json:"invited"`
	BetAmount   float64 `json:"bet_amount"`
	BetCurrency string  `json:"bet_currency"`
}

type GeneralInvitationRequest struct {
	Inviter uint `json:"inviter"`
	Invited uint `json:"invited"`
}

type AcceptStartGameRequest struct {
	TeamID     uint `json:"team_id"`
	OpponentID uint `json:"opponent_id"`
}

type GameStatResponse struct {
	Status int    `json:"status"`
	Code   int    `json:"code"`
	Result string `json:"result"`
}

type GameStatNow struct {
	IsError        error
	IsFinishedPast bool
	WinnerId       uint
}

var InvitationSent = uint8(0)
var InvitationAccepted = uint8(1)
var GameStarted = uint8(2)
var GameFinished = uint8(3)
