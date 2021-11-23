package main

import (
	"github.com/aicam/AlarmServer/DB"
	"github.com/aicam/AlarmServer/server"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// migration
	s := server.NewServer()
	s.DB = DB.DbSqlMigration("aicam:021021ali@tcp(127.0.0.1:3306)/messenger_api?charset=utf8mb4&parseTime=True")
	s.Routes()
	log.Println(time.Now())
	var user DB.UsersData
	username := "aicam"
	key := os.Getenv("SERVER_KEY")
	log.Println(key)
	if err := s.DB.Where(DB.UsersData{Username: username}).Find(&user).Error; err != nil {
		s.DB.Save(&DB.UsersData{
			Username:   username,
			LastOnline: time.Now(),
		})
	}
	err := http.ListenAndServe("0.0.0.0:4300", s.Router)
	if err != nil {
		log.Print(err)
	}

}
