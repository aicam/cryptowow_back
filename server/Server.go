package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
)

type Server struct {
	DB               *gorm.DB
	Router           *gin.Engine
	SocketConnection websocket.Upgrader
	WoWInfo          struct {
		Mounts     MountsInfo
		Companions CompanionsInfo
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max, username")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

// Here we create our new server
func NewServer() *Server {
	router := gin.Default()
	// here we opened cors for all
	router.Use(CORS())
	jsonFile, err := os.Open("WoWUtils/mounts_info.json")
	if err != nil {
		log.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var mounts MountsInfo
	err = json.Unmarshal(byteValue, &mounts)

	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	var mounts MountsInfo
	err = json.Unmarshal(byteValue, &mounts)

	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	return &Server{
		DB:      nil,
		Router:  router,
		WoWInfo: struct{ Mounts MountsInfo }{Mounts: mounts},
	}
}
