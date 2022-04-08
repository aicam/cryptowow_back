package Bridge

import "math/rand"

type NewGameParams struct {
	ArenaFilePath string
	TeamID1       int
	TeamID2       int
	LeaderName1   string
	LeaderName2   string
	ArenaType     string
}

var ArenaTypes = []string{"4", "5", "8", "10", "11"}

func appendNewGame(params NewGameParams) {
	// TODO: this part should replaced with remote WoW server endpoints
	if params.ArenaType == "" {
		params.ArenaType = ArenaTypes[rand.Intn(len(ArenaTypes))]
	}

}
