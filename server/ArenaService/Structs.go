package ArenaService

type ArenaTeam struct {
	ArenaTeamID int    `json:"arena_team_id" gorm:"column:arenaTeamId"`
	LeaderName  string `json:"leader_name" gorm:"column:name"`
	LeaderGUID  int    `json:"leader_guid" gorm:"column:captainGuid"`
	ArenaType   int16  `json:"arena_type" gorm:"column:type"`
	SeasonGames int16  `json:"season_games" gorm:"column:seasonGames"`
	SeasonWins  int16  `json:"season_wins" gorm:"column:seasonWins"`
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
