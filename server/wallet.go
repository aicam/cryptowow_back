package server

import (
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WalletCurrencies() []string {
	return []string{"Bitcoin", "Etherium", "CWT"}
}

func (s *Server) GetWalletInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		username := context.GetHeader("username")
		currencies := WalletCurrencies()
		var wallets []database.Wallet
		s.DB.Where(&database.Wallet{Name: username}).Find(&wallets)
		context.JSON(http.StatusOK, struct {
			Currencies []string          `json:"currencies"`
			Wallets    []database.Wallet `json:"wallets"`
		}{
			Currencies: currencies,
			Wallets:    wallets,
		})
	}
}
