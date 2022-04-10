package Bridge

import (
	"log"
	"math/rand"
	"os"
	"strconv"
)

type NewGameParams struct {
	ArenaFilePath string
	TeamID1       int
	TeamID2       int
	LeaderName1   string
	LeaderName2   string
	ArenaType     string
	MapType       string
}

var MapTypes = []string{"4", "5", "8", "10", "11"}

func AppendNewGame(params NewGameParams) {

	// TODO: this part should replaced with remote WoW server endpoints
	if params.MapType == "" {
		params.MapType = MapTypes[rand.Intn(len(MapTypes))]
	}
	newArenaStr := strconv.Itoa(params.TeamID1) + "," +
		params.LeaderName1 + "," +
		strconv.Itoa(params.TeamID2) + "," +
		params.LeaderName2 + "," +
		params.ArenaType + "," +
		params.MapType + "\n"

	file, err := os.OpenFile(params.ArenaFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if _, err := file.WriteString(newArenaStr); err != nil {
		log.Fatal(err)
	}

}
