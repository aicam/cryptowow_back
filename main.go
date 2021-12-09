package main

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// migration
	s := server.NewServer()
	s.DB = database.DbSqlMigration("aicam:021021ali@tcp(127.0.0.1:3306)/messenger_api?charset=utf8mb4&parseTime=True")
	s.Routes()
	log.Println(time.Now())
	var user database.UsersData
	username := "aicam"
	key := os.Getenv("SERVER_KEY")
	log.Println(key)
	if err := s.DB.Where(database.UsersData{Username: username}).Find(&user).Error; err != nil {
		s.DB.Save(&database.UsersData{
			Username: username,
			Password: server.MD5("ali"),
		})
	} else {
		user.Password = server.MD5("ali")
		s.DB.Save(&user)
	}
	err := http.ListenAndServe("0.0.0.0:4300", s.Router)
	if err != nil {
		log.Print(err)
	}

}
