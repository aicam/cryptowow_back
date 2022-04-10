package Bridge

type ArenaTeam struct {
	ArenaTeamID int    `json:"arena_team_id" gorm:"column:arenaTeamId"`
	LeaderName  string `json:"leader_name" gorm:"column:name"`
	LeaderGUID  int    `json:"leader_guid" gorm:"column:captainGuid"`
	ArenaType   int16  `json:"arena_type" gorm:"column:type"`
	SeasonGames int16  `json:"season_games" gorm:"column:seasonGames"`
	SeasonWins  int16  `json:"season_wins" gorm:"column:seasonWins"`
}
