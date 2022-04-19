package ArenaService

type ArenaTeam struct {
	ArenaTeamID int    `json:"arena_team_id" gorm:"column:arenaTeamId"`
	LeaderName  string `json:"leader_name" gorm:"column:name"`
	LeaderGUID  int    `json:"leader_guid" gorm:"column:captainGuid"`
	ArenaType   int8   `json:"arena_type" gorm:"column:type"`
	SeasonGames uint   `json:"season_games" gorm:"column:seasonGames"`
	SeasonWins  uint   `json:"season_wins" gorm:"column:seasonWins"`
}

type InviteRequest struct {
	Inviter     int     `json:"inviter"`
	Invited     int     `json:"invited"`
	BetAmount   float64 `json:"bet_amount"`
	BetCurrency string  `json:"bet_currency"`
}

type GeneralInvitationRequest struct {
	Inviter int `json:"inviter"`
	Invited int `json:"invited"`
}

type AcceptStartGameRequest struct {
	TeamID     int `json:"team_id"`
	OpponentID int `json:"opponent_id"`
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
