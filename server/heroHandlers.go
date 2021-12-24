package server

import (
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

func checkHeroIsOnline(c *gin.Context, DB *gorm.DB, heroName string, username string) (Hero, bool) {
	var hero Hero
	var accID int
	DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id FROM account WHERE username='" + username + "'").Scan(&accID)
	log.Println(accID)
	err := DB.Clauses(dbresolver.Use("characters")).Raw("SELECT account, race, online from characters WHERE name='" + heroName + "'").First(&hero).Error
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

func (s *Server) RestoreHero() gin.HandlerFunc {
	return func(context *gin.Context) {
		heroName := context.Param("hero_name")
		username := context.GetHeader("username")
		hero, err := checkHeroIsOnline(context, s.DB, heroName, username)
		if !err {
			return
		}
		var heroHomeLoc HeroPosition
		switch hero.Race {
		case 2, 5, 6, 10, 8:
			heroHomeLoc = Home.Horde
		default:
			heroHomeLoc = Home.Alliance
		}
		s.DB.Clauses(dbresolver.Use("characters")).Raw("UPDATE characters SET " +
			"map=" + strconv.Itoa(heroHomeLoc.Map) + "," +
			"position_x=" + strconv.FormatFloat(float64(heroHomeLoc.PositionX), 'f', 4, 64) + "," +
			"position_y=" + strconv.FormatFloat(float64(heroHomeLoc.PositionY), 'f', 4, 64) + "," +
			"position_z=" + strconv.FormatFloat(float64(heroHomeLoc.PositionZ), 'f', 4, 64) + "," +
			"WHERE name = '" + heroName + "'")
		context.JSON(http.StatusOK, actionResult(1, "Hero resotred successfully"))
	}
}
