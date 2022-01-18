package server

import (
	"errors"
	"github.com/aicam/cryptowow_back/GMReqs"
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
)

func CheckGiftIsAllowed(giftID uint, username string, DB *gorm.DB) (database.Gifts, error) {
	var gift database.Gifts
	err := DB.Where(&database.Gifts{Username: username, GiftID: giftID}).First(&gift).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gift, errors.New("No gift exist")
	}
	if gift.Used {
		return gift, errors.New("Gift is used for hero " + gift.UsedHeroName)
	}
	return gift, nil
}

func LevelUp(heroName string, heroID, classID, race int, DB *gorm.DB) {
	log.Print(classID)
	for _, items := range LevelUpGift[classID] {
		GMReqs.AddItems("First hero gift!", "Welcome to CryptoWoW server!!", heroName, items)
	}
	switch race {
	case 2, 5, 6, 10, 8:
		for _, items := range HordeWeapons {
			GMReqs.AddItems("Horde weapons!", "Welcome to CryptoWoW server!!", heroName, items)
		}
	default:
		GMReqs.AddItems("Alliance weapons!", "Welcome to CryptoWoW server!!", heroName, AllianceWeapons)
	}
	DB.Clauses(dbresolver.Use("characters")).Exec("UPDATE characters SET level=80, money=500000 WHERE name='" + heroName + "';")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '293', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '413', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '45', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '43', '1', '1');")
	DB.Clauses(dbresolver.Use("characters")).Exec("INSERT INTO `character_skills` (`guid`, `skill`, `value`, `max`) VALUES ('" + strconv.Itoa(heroID) + "', '415', '1', '1');")
}

func (s *Server) GiftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		heroName := c.Param("hero_name")
		giftID := c.Param("gift_id")
		gID, err := strconv.Atoi(giftID)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -2,
				Body:       "Problem!!!!",
			})
		}
		gift, err := CheckGiftIsAllowed(uint(gID), username, s.DB)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		hero, errH := CheckHeroIsAllowed(c, s.DB, heroName, username)
		if !errH {
			return
		}
		hero.Name = heroName
		switch gID {
		case 1:
			// Uncomment
			//LevelUp(heroName, hero.HeroID, hero.Class, int(hero.Race), s.DB)
			//TeleportHeroHome(hero, s.DB)
			gift.Used = true
			gift.UsedHeroName = heroName
			s.DB.Save(&gift)
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Gift used successfully!!",
		})
	}
}
