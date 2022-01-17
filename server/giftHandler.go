package server

import (
	"github.com/aicam/cryptowow_back/GMReqs"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
)

func LevelUp(heroName string, heroID, classID int, DB *gorm.DB) {
	log.Print(classID)
	for _, items := range LevelUpGift[classID] {
		GMReqs.AddItems("First hero gift!", "Welcome to CryptoWoW server!!", heroName, items)
	}
	DB.Clauses(dbresolver.Use("characters")).Exec("UPDATE characters SET level=80 WHERE name='" + heroName + "';")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '293', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '413', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '45', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '43', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '415', '1', '1');")
}

func (s *Server) LevelUpGiftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		heroName := c.Param("hero_name")
		hero, err := CheckHeroIsAllowed(c, s.DB, heroName, username)
		if !err {
			return
		}
		LevelUp(heroName, hero.HeroID, hero.Class, s.DB)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Gift used successfully!!",
		})
	}
}
