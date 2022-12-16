package AdminRouter

import (
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Service) AddBalanceToAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqParams AddBalanceToAccountRequest
		err := c.BindJSON(&reqParams)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				StatusCode: -1,
				Body:       "Error in parsing inputs",
			})
			return
		}
		WalletService.AddBalance(reqParams.Username, reqParams.CurrencyID, reqParams.Amount, s.DB)
		s.PP.Counters["admin_service_balance_added_count"].Inc()
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Added successfully",
		})
	}
}
