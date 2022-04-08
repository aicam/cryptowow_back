package main

import (
	"log"
	"os"
)

type NewGameParams struct {
	ArenaFilePath string
	TeamID1       int
	TeamID2       int
	LeaderName1   string
	LeaderName2   string
	ArenaType     string `default:"random"`
}

func main() {
	arenaFile := os.Getenv("CUSTOMARENAFILELOC")

	file, err := os.OpenFile(arenaFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	if _, err := file.WriteString("second line"); err != nil {
		log.Fatal(err)
	}
}
