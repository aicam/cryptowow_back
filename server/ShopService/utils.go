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

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func (s *Service) findFirstEmptySlotInBag(heroId uint) (int, int) {

	var charInventories []CharacterInventory
	s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT bag, slot, item FROM character_inventory WHERE guid=" + strconv.Itoa(int(heroId))).Find(&charInventories)

	charInventoriesMap := make(map[BagRow]int)
	for _, inv := range charInventories {
		charInventoriesMap[BagRow{
			Bag:  inv.Bag,
			Slot: inv.Slot,
		}] = inv.ItemId
	}

	// check backpack slots
	for _, i := range makeRange(s.BagsInfo.BackPackStart, s.BagsInfo.BackPackEnd) {
		_, ok := charInventoriesMap[BagRow{
			Bag:  0,
			Slot: i,
		}]
		if !ok {
			return 0, i
		}
	}

	// fill bags info
	charBags := make(map[int]int)
	for _, i := range makeRange(s.BagsInfo.BagSlotsStart, s.BagsInfo.BagSlotsEnd) {
		item, ok := charInventoriesMap[BagRow{
			Bag:  0,
			Slot: i,
		}]
		if ok {
			charBags[item] = s.BagItems[item]
		}
	}

	// check bag slots
	for bag, maxSlot := range charBags {
		for _, i := range makeRange(0, maxSlot) {
			_, ok := charInventoriesMap[BagRow{
				Bag:  bag,
				Slot: i,
			}]
			if !ok {
				return bag, i
			}
		}
	}
	return 0, 0

}
