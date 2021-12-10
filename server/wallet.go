package server

import "github.com/gin-gonic/gin"

func WalletCurrencies() []string {
	return []string{"Bitcoin", "Etherium", "CWT"}
}

func (s *Server) GetWalletInfo() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
