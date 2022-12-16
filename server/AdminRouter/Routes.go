package AdminRouter

import "github.com/gin-gonic/gin"

func (s *Service) AddRoutes(router *gin.Engine, middleware func() gin.HandlerFunc) {
	router.POST("/admin/add_balance/", middleware(), s.AddBalanceToAccount())
}
