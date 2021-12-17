package server

func (s *Server) Routes() {
	s.Router.POST("/login", s.GetToken())
	s.Router.GET("/get_available_wallets", s.checkToken(), s.AvailableWallets())
	s.Router.POST("/add_info", s.checkToken(), s.AddInfo())
	s.Router.GET("/get_info/:offset", s.checkToken(), s.GetInfo())
	s.Router.POST("/register", s.AddUser())
	// Login required (WoW DB)
	s.Router.GET("/wow/get_info", s.checkToken(), s.ReturnUserInfo())
	// Wallet
	s.Router.GET("/wallet/get_info", s.checkToken(), s.GetWalletInfo())
}
