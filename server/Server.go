package server

import (
	"context"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/monitoring"
	"github.com/aicam/cryptowow_back/server/AdminRouter"
	"github.com/aicam/cryptowow_back/server/ArenaService"
	"github.com/aicam/cryptowow_back/server/GlobalStructs"
	"github.com/aicam/cryptowow_back/server/ShopService"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		Bags       GlobalStructs.BagsInfo
	}
	TrinityCoreBridgeVars    map[string]string
	TrinityCoreBridgeService ArenaService.Service
	ShopService              ShopService.Service
	AdminRouter              AdminRouter.Service
	PP                       monitoring.PrometheusParams
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
		os.Create("Server_Logs.txt")
		loggerFile, _ = os.OpenFile("Server_Logs.txt", os.O_APPEND|os.O_WRONLY, 0644)
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
	companions := parseCompanions("./WoWUtils/companions_info.json")
	bags := parseBags("./WoWUtils/bags_info.json")
	mounts := parseMounts("./WoWUtils/mounts_info.json")

	// generate database gorm structure
	log.Println(os.Getenv("MAINMYSQLCONNECTION"))
	DBStruct := database.DbSqlMigration(os.Getenv("MAINMYSQLCONNECTION"))
	// redis server
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISCONNECTION"),
		Password: os.Getenv("REDISPASS"), // no password set
		DB:       0,                      // use default DB
	})
	err := rdb.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		DB:     DBStruct,
		Router: router,
		WoWInfo: struct {
			Mounts     MountsInfo
			Companions CompanionsInfo
			Bags       GlobalStructs.BagsInfo
		}{Mounts: mounts, Companions: companions, Bags: bags},
		PP:                    monitoring.GetGlobalPrometheusParams(),
		TrinityCoreBridgeVars: make(map[string]string),
		TrinityCoreBridgeService: ArenaService.Service{DB: DBStruct,
			Rdb:     rdb,
			Context: context.Background(),
			PP:      monitoring.GetArenaBetServicePrometheusParams()},
		ShopService: ShopService.NewService(
			DBStruct,
			monitoring.GetShopPrometheusParams(),
			bags,
		),
		AdminRouter: AdminRouter.Service{
			DB: DBStruct,
			PP: monitoring.GetAdminPrometheusParams(),
		},
	}
}
