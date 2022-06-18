package ShopService

import (
	"github.com/aicam/cryptowow_back/GMReqs"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/LogService"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusOK, actionResult(1, "Item sent by mail to hero successfully!"))
	}
}
