package server

import (
	"context"
	"encoding/json"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/monitoring"
	"github.com/aicam/cryptowow_back/server/ArenaService"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
)

type Response struct {
	StatusCode int         `json:"status"`
	Body       interface{} `json:"body"`
}

type Server struct {
	DB               *gorm.DB
	Router           *gin.Engine
	SocketConnection websocket.Upgrader
	WoWInfo          struct {
		Mounts     MountsInfo
		Companions CompanionsInfo
	}
	TrinityCoreBridgeVars   map[string]string
	TrinityCoreBridgeServer ArenaService.Service
	PP                      monitoring.PrometheusParams
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

func SetupLogger() {
	loggerFile, e := os.OpenFile("Server_Logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if e != nil {
		log.Fatal(e)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(loggerFile)
}

// Here we create our new server
func NewServer() *Server {
	SetupLogger()
	router := gin.Default()
	// here we opened cors for all
	router.Use(CORS())
	jsonFile, err := os.Open("WoWUtils/mounts_info.json")
	if err != nil {
		log.Println(err)
	}
	// TODO: clean up reading files
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var mounts MountsInfo
	err = json.Unmarshal(byteValue, &mounts)

	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	jsonFile, err = os.Open("WoWUtils/companions_info.json")
	if err != nil {
		log.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)

	var companions CompanionsInfo
	err = json.Unmarshal(byteValue, &companions)

	if err != nil {
		log.Print(err)
		os.Exit(-1)
	}

	// generate database gorm structure
	log.Println(os.Getenv("MYSQLCONNECTION"))
	DBStruct := database.DbSqlMigration(os.Getenv("MYSQLCONNECTION"))

	// redis server
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: os.Getenv("REDISPASS"), // no password set
		DB:       0,                      // use default DB
	})

	return &Server{
		DB:     DBStruct,
		Router: router,
		WoWInfo: struct {
			Mounts     MountsInfo
			Companions CompanionsInfo
		}{Mounts: mounts, Companions: companions},
		PP:                    monitoring.GetGlobalPrometheusParams(),
		TrinityCoreBridgeVars: make(map[string]string),
		TrinityCoreBridgeServer: ArenaService.Service{DB: DBStruct, Rdb: rdb, Context: context.Background(),
			PP: monitoring.GetArenaBetServicePrometheusParams()},
	}
}
