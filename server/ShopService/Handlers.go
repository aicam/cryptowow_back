package ShopService

import (
	"github.com/aicam/cryptowow_back/GMReqs"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/LogService"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (s *Service) BuyItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var req BuyItemRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-1, "Wrong request"))
			return
		}
		_, errCheck := checkHeroIsAllowed(c, s.DB, req.HeroName, username)
		if !errCheck {
			return
		}
		var item database.ShopItems
		err = s.DB.Where(" id = " + req.ItemID).First(&item).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, actionResult(-3, "Wrong request!"))
			LogService.LogCrashinShop("Buy Item", "ItemID in database check")
		}

		err = WalletService.ReduceBalance(false, username, item.CurrencyID, item.Amount, s.DB)

		if err != nil {
			c.JSON(http.StatusOK, actionResult(-3, "Not enough balance!"))
			LogService.LogCrashinShop("Buy Item", "Not enough balance")
		}

		GMReqs.AddItems("Shop order", "Your shopping order is delivered", req.HeroName, item.ItemID)

		var boughtItems database.BoughtItems
		boughtItems.ItemID = req.ItemID
		boughtItems.HeroName = req.HeroName
		boughtItems.Username = username
		s.DB.Save(&boughtItems)

		c.JSON(http.StatusOK, actionResult(1, "Item sent by mail to hero successfully!"))

		itemId, _ := strconv.Atoi(item.ItemID)
		s.PP.Histograms["shop_service_item_sold_by_id"].Observe(float64(itemId))
		s.PP.Counters["shop_service_item_sold_count"].Inc()
	}
}

func (s *Service) AddItemToShop() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req database.ShopItems
		s.DB.Save(&req)
		c.JSON(http.StatusOK, actionResult(1, "Added successfully!"))
	}
}

func (s *Service) GetShopItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []database.ShopItems
		s.DB.Find(&items)
		c.JSON(http.StatusOK, actionResult(1, items))
	}
}

func (s *Service) TempTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, findFirstEmptySlotInBag(s.DB, 1))
	}
}
