package ShopService

import "github.com/gin-gonic/gin"

func (s *Service) BuyItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")

	}
}
