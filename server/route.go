package server

func (s *Server) Routes() {
	s.Router.POST("/login", s.GetToken())
	s.Router.GET("/get_available_wallets", s.checkToken(), s.AvailableWallets())
	s.Router.POST("/register", s.AddUser())
	// Login required (WoW DB)
	s.Router.GET("/wow/get_info", s.checkToken(), s.ReturnUserInfo())
	s.Router.GET("/wow/hero_info/:hero_name", s.checkToken(), s.ReturnHeroInfo())
	s.Router.GET("/wow/restore_hero/:hero_name", s.checkToken(), s.RestoreHero())
	s.Router.POST("/wow/sell_hero", s.checkToken(), s.SellHero())

	// Wallet
	s.Router.GET("/wallet/get_info", s.checkToken(), s.GetWalletInfo())
}
