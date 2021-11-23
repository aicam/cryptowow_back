package server

import (
	"encoding/hex"
	"github.com/aicam/AlarmServer/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		newUser := context.Param("username")
		s.DB.Save(&DB.UsersData{
			Username:   newUser,
			LastOnline: time.Now(),
		})
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Added",
		})
	}
}

func (s *Server) GetToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user DB.UsersData
		username := context.GetHeader("username")
		key := []byte("Ali@Kian")
		if err := s.DB.Where(DB.UsersData{Username: username}).First(&user).Error; err != nil {
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				Body:       "Invalid data",
			})
			return
		}
		user.LastOnline = time.Now()
		s.DB.Save(&user)
		token, err := DesEncrypt([]byte(username), key)
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
		var jsData DB.WebData
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
		var DBData []DB.WebData
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
