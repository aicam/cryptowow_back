package main

import (
	"github.com/aicam/cryptowow_back/server/Bridge"
	"os"
)

func main() {
	arenaFile := os.Getenv("CUSTOMARENAFILELOC")
	params := Bridge.NewGameParams{
		ArenaFilePath: arenaFile,
		TeamID1:       1,
		TeamID2:       2,
		LeaderName1:   "Matarsak",
		LeaderName2:   "Koonde",
		ArenaType:     "2",
		//MapType:       "4",
	}
	Bridge.AppendNewGame(params)
}
