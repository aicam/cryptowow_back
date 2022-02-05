package server

import (
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
)

func actionResult(statusCode int, body string) struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
} {
	return struct {
		Status int    `json:"status"`
		Body   string `json:"body"`
	}{Status: statusCode, Body: body}
}

func CheckHeroIsAllowed(c *gin.Context, DB *gorm.DB, heroName string, username string) (Hero, bool) {
	var hero Hero
	var accID int
	DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id FROM account WHERE username='" + username + "'").Scan(&accID)
	log.Println(accID)
	err := DB.Clauses(dbresolver.Use("characters")).Raw("SELECT account, guid, race, online, gender, level, class, money, totaltime, totalKills from characters WHERE name='" + heroName + "'").First(&hero).Error
	if err != nil {
		c.JSON(http.StatusOK, actionResult(-3, "Malicious activity detected"))
		return hero, false
	}
	if hero.AccountID != accID {
		c.JSON(http.StatusOK, actionResult(-8, "Malicious activity detected"))
		return hero, false
	}
	if hero.Online {
		c.JSON(http.StatusOK, actionResult(-1, "Hero is currently online"))
		return hero, false
	}
	return hero, true
}

func TeleportHeroHome(hero Hero, DB *gorm.DB) {
	var heroHomeLoc HeroPosition
	switch hero.Race {
	case 2, 5, 6, 10, 8:
		heroHomeLoc = Home.Horde
	default:
		heroHomeLoc = Home.Alliance
	}
	DB.Clauses(dbresolver.Use("characters")).Exec("UPDATE characters SET " +
		"map=" + strconv.Itoa(heroHomeLoc.Map) + "," +
		"position_x=" + strconv.FormatFloat(float64(heroHomeLoc.PositionX), 'f', 4, 64) + "," +
		"position_y=" + strconv.FormatFloat(float64(heroHomeLoc.PositionY), 'f', 4, 64) + "," +
		"position_z=" + strconv.FormatFloat(float64(heroHomeLoc.PositionZ), 'f', 4, 64) +
		"WHERE name = '" + hero.Name + "'")
}

func (s *Server) RestoreHero() gin.HandlerFunc {
	return func(c *gin.Context) {
		heroName := c.Param("hero_name")
		username := c.GetHeader("username")
		hero, err := CheckHeroIsAllowed(c, s.DB, heroName, username)
		if !err {
			return
		}
		hero.Name = heroName
		TeleportHeroHome(hero, s.DB)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Hero resotred successfully",
		})
	}
}

func (s *Server) SellHero() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody struct {
			HeroName  string `json:"hero_name"`
			HeroPrice string `json:"hero_price"`
			Note      string `json:"note"`
		}
		errBody := c.BindJSON(&reqBody)
		if errBody != nil {
			c.JSON(http.StatusOK, actionResult(-1, "Invalid request"))
			return
		}
		username := c.GetHeader("username")
		hero, err := CheckHeroIsAllowed(c, s.DB, reqBody.HeroName, username)
		if !err {
			return
		}
		var sellingHero database.SellingHeros
		if err := s.DB.Where(&database.SellingHeros{HeroName: reqBody.HeroName}).First(&sellingHero).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var newSellingHero database.SellingHeros
			newSellingHero.Price = reqBody.HeroPrice
			newSellingHero.HeroName = reqBody.HeroName
			newSellingHero.Note = reqBody.Note
			newSellingHero.Race = int(hero.Race)
			newSellingHero.Username = username
			newSellingHero.HeroID = hero.HeroID
			newSellingHero.Class = hero.Class
			newSellingHero.Level = hero.Level
			newSellingHero.Money = hero.Money
			newSellingHero.TotalKills = hero.TotalKills
			newSellingHero.TotalTime = hero.TotalTime
			newSellingHero.Gender = hero.Gender
			s.DB.Save(&newSellingHero)
			c.JSON(http.StatusOK, actionResult(1, "Added successfully"))
		} else {
			c.JSON(http.StatusOK, actionResult(-1, "Hero is already for sale!"))
			return
		}
	}
}

func (s *Server) CancellSellingHero() gin.HandlerFunc {
	return func(c *gin.Context) {
		heroName := c.Param("hero_name")
		username := c.GetHeader("username")
		_, err := CheckHeroIsAllowed(c, s.DB, heroName, username)
		if !err {
			return
		}
		var sHero database.SellingHeros
		findErr := s.DB.Where(database.SellingHeros{Username: username, HeroName: heroName}).First(&sHero).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, actionResult(-1, "Hero is not for sale already!!"))
			return
		}
		s.DB.Delete(&sHero)
		c.JSON(http.StatusOK, actionResult(1, "Selling canceled"))
	}
}

func (s *Server) ReturnSellingHeros() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sellingheros []database.SellingHeros
		s.DB.Find(&sellingheros)
		c.JSON(http.StatusOK, sellingheros)
	}
}
