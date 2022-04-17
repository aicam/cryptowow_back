package server

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetWalletAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "0xDE4C72835bcC0041Dd1B446BfD0D85bE346BC0A2",
		})
	}
}

func (s *Server) GetUserTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var txs []database.TransactionLog
		s.DB.Where(&database.TransactionLog{Username: username}).Find(&txs)
		c.JSON(http.StatusOK, txs)
	}
}

func (s *Server) AddCashOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var co database.CashOutRequest
		err := c.BindJSON(&co)
		if err != nil {
			return
		}
		err = WalletService.ReduceBalance(username, co.CurrencyID, co.Amount, s.DB)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		co.PendingStage = 0
		co.Username = username
		s.DB.Save(&co)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Cashout request added successfully",
		})
	}
}

func (s *Server) ReturnCashOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		var cos []database.CashOutRequest
		username := c.GetHeader("username")
		s.DB.Where(&database.CashOutRequest{Username: username}).Find(&cos)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       cos,
		})
	}
}
