package ShopService

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/LogService"
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
			LogService.LogBadItemIDinShop("Buy Item", "ItemID in database check")
		}

	}
}
