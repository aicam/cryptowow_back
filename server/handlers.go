package server

import (
	"encoding/hex"
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func (s *Server) ReturnHeroInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		heroName := context.Param("hero_name")
		var heroInfo HeroInfo
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT guid, name, race, gender, level, class, equipmentCache FROM characters WHERE name='" + heroName + "'").Scan(&heroInfo)
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT achievement FROM character_achievement WHERE guid=" + strconv.Itoa(heroInfo.ID)).Scan(&heroInfo)
		context.JSON(http.StatusOK, heroInfo)
	}
}

func (s *Server) AvailableWallets() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, struct {
			Wallets []string `json:"wallets"`
		}{Wallets: []string{"Trust Wallet", "Bitpay"}})
	}
}

func (s *Server) ReturnUserInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		username := context.GetHeader("username")
		var id int
		log.Println("SELECT id from account where username='" + strings.ToUpper(username) + "'")
		s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account WHERE username='" + strings.ToUpper(username) + "'").Scan(&id)
		var heros []Hero
		s.DB.Raw("SELECT name, race, gender, level, class FROM characters WHERE account=" + strconv.Itoa(id)).Scan(&heros)

		var gifts []database.Gifts
		s.DB.Where(&database.Gifts{
			Username: username,
		}).Find(&gifts)

		currencies := WalletCurrencies()
		var wallets []database.Wallet
		s.DB.Where(&database.Wallet{Name: username}).Find(&wallets)
		context.JSON(http.StatusOK, struct {
			Heros      []Hero            `json:"heros"`
			Gifts      []database.Gifts  `json:"gifts"`
			Wallets    []database.Wallet `json:"wallets"`
			Currencies []string          `json:"currencies"`
		}{Heros: heros, Gifts: gifts, Wallets: wallets, Currencies: currencies})
	}
}
func (s *Server) AddUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var newUser database.UsersData
		var existUser database.UsersData
		_ = context.BindJSON(&newUser)
		if err := s.DB.Where(&database.UsersData{Username: newUser.Username}).Find(&existUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusOK, Response{
				StatusCode: 0,
				Body:       "Username exist",
			})
			return
		}
		gift := database.Gifts{
			Username:     newUser.Username,
			Description:  "Level up first hero free!",
			Action:       "lvlup",
			Condition:    "Register",
			Used:         false,
			UsedHeroName: "",
		}
		s.DB.Save(&gift)
		newUser.Password = MD5(newUser.Password)
		log.Println(newUser.Password)
		s.DB.Save(&newUser)

		// Uncomment
		//GMReqs.CreateAccount(newUser.Username, newUser.Password)

		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Added",
		})
	}
}

func (s *Server) GetToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user database.UsersData
		err := context.BindJSON(&user)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}
		key := []byte("Ali@Kian")
		if err := s.DB.Where(database.UsersData{Username: user.Username,
			Password: MD5(user.Password)}).First(&user).Error; err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}

		token, err := DesEncrypt([]byte(user.Username), key)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       hex.EncodeToString(token),
		})
	}
}

func (s *Server) AddInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		var jsData database.WebData
		err := context.BindJSON(&jsData)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		username := context.GetHeader("username")
		_ = "2006-01-02T15:04:05Z07:00"
		if jsData.ArmeniaTime.Year() != 1 {
			go sendNotificationByPushOver(jsData.ArmeniaTxt, "Armenia Time found")
			go sendNotificationByIFTTT(jsData.ArmeniaTxt, "Armenia Time found")
			go SendNotificationByTelegram(jsData.ArmeniaTxt, "Armenia Time found")
		}
		if jsData.DubaiTime.Month() >= 8 {
			go sendNotificationByPushOver(jsData.DubaiTxt, "Dubai Time found")
			go sendNotificationByIFTTT(jsData.DubaiTxt, "Dubai Time found")
			go SendNotificationByTelegram(jsData.DubaiTxt, "Dubai Time found")
		}
		if jsData.TurkeyTime.Year() != 1 {
			go sendNotificationByPushOver(jsData.DubaiTxt, "Dubai Time found")
			go sendNotificationByIFTTT(jsData.DubaiTxt, "Dubai Time found")
			go SendNotificationByTelegram(jsData.DubaiTxt, "Dubai Time found")
		}
		if jsData.DubaiTime.Year() == 1 {
			jsData.DubaiTime = time.Now()
		}
		if jsData.TurkeyTime.Year() == 1 {
			jsData.TurkeyTime = time.Now()
		}
		if jsData.ArmeniaTime.Year() == 1 {
			jsData.ArmeniaTime = time.Now()
		}
		//if jsData.Priority >= 0 {
		//	timeFounded, err := time.Parse(layout, jsData.ClosestDate)
		//	if err != nil {
		//		context.JSON(http.StatusOK, Response{
		//			StatusCode: -1,
		//			Body:       err.Error(),
		//		})
		//		return
		//	}
		//	if jsData.Priority > 0 {
		//		log.Print(strconv.Itoa(int(timeFounded.Sub(time.Now()).Hours() / 24)))
		//		go sendNotificationByPushOver("In "+timeFounded.Month().String()+" "+strconv.Itoa(timeFounded.Day()), "Time found in "+
		//			strconv.Itoa(int(timeFounded.Sub(time.Now()).Hours()/24))+" days"+" in "+jsData.Country)
		//	}
		//}

		jsData.Username = username
		s.DB.Save(&jsData)
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Data saved successfully!",
		})
	}
}

func (s *Server) GetInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		var DBData []database.WebData
		offset, err := strconv.Atoi(context.Param("offset"))
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		s.DB.Find(&DBData)
		context.JSON(http.StatusOK, DBData[len(DBData)-offset:])
	}
}
