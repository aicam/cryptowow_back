package main

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server"
	"gorm.io/plugin/dbresolver"
	"log"
)

type Hero struct {
	Name   string
	Race   bool
	Gender bool
	Level  int
}

func main() {
	// migration
	s := server.NewServer()
	s.DB = database.DbSqlMigration("aicam:021021ali@tcp(127.0.0.1:3306)/messenger_api?charset=utf8mb4&parseTime=True")
	var hero Hero
	var id int
	s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account where username='ALI'").Scan(&id)
	log.Println(id)
	s.DB.Raw("SELECT name, race, gender, level FROM characters").Scan(&hero)
	log.Print(hero)
}
