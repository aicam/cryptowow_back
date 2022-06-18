package ShopService

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"net/http"
)

type Hero struct {
	AccountID  int    `json:"account_id" gorm:"column:account"`
	HeroID     int    `json:"hero_id" gorm:"column:guid"`
	Name       string `json:"name"`
	Race       uint   `json:"race"`
	Gender     bool   `json:"gender"`
	Level      int    `json:"level"`
	Class      int    `json:"class"`
	Online     bool   `json:"online"`
	Money      int    `json:"money"`
	TotalTime  int    `json:"total_time" gorm:"column:totaltime"`
	TotalKills int    `json:"total_kills" gorm:"column:totalKills"`
}

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
