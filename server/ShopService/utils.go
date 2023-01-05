package ShopService

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"net/http"
	"strconv"
)

func actionResult(statusCode int, body interface{}) struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
} {
	return struct {
		Status int         `json:"status"`
		Body   interface{} `json:"body"`
	}{Status: statusCode, Body: body}
}

func checkHeroIsAllowed(c *gin.Context, DB *gorm.DB, heroName string, username string) (Hero, bool) {
	var hero Hero
	var accID int
	DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id FROM account WHERE username='" + username + "'").Scan(&accID)

	err := DB.Clauses(dbresolver.Use("characters")).Raw("SELECT account, guid, race, online, gender, level, class, money, totaltime, totalKills from characters WHERE name='" + heroName + "'").First(&hero).Error
	if err != nil {
		c.JSON(http.StatusOK, actionResult(-3, "Malicious activity detected"))
		return hero, false
	}
	if hero.AccountID != accID {
		c.JSON(http.StatusOK, actionResult(-8, "Malicious activity detected"))
		return hero, false
	}
	return hero, true
}

func findFirstEmptySlotInBag(DB *gorm.DB, heroId uint) []CharacterInventory {
	var charInventories []CharacterInventory
	DB.Clauses(dbresolver.Use("characters")).Raw("SELECT bag, slot, item FROM character_inventory WHERE guid=" + strconv.Itoa(int(heroId))).Find(&charInventories)
	return charInventories

}
