package server

func (s *Server) Routes() {
	s.Router.POST("/login", s.GetToken())
	// Regiser
	s.Router.GET("/get_csrf", s.GetCSRFToken())
	s.Router.POST("/register", s.AddUser())
	// index
	s.Router.GET("/server_status", s.GetServerStatus())
	s.Router.GET("/events", s.ReturnEvents())
	s.Router.GET("/get_available_wallets", s.checkToken(), s.AvailableWallets())
	// Login required (WoW DB)
	s.Router.GET("/wow/get_info", s.checkToken(), s.ReturnUserInfo())
	s.Router.GET("/wow/hero_info/:hero_name", s.checkToken(), s.GetHeroInfo())
	s.Router.GET("/wow/restore_hero/:hero_name", s.checkToken(), s.RestoreHero())
	// Hero selling functionality
	s.Router.POST("/wow/sell_hero", s.checkToken(), s.SellHero())
	s.Router.GET("/wow/cancel_selling_hero/:hero_name", s.checkToken(), s.CancellSellingHero())
	s.Router.GET("/wow/selling_heros", s.checkToken(), s.ReturnSellingHeros())
	s.Router.POST("/wow/buy_hero", s.checkToken(), s.BuyHero())
	// Wallet
	s.Router.POST("/wallet/add_transaction", s.checkToken(), s.AddTransaction())
	s.Router.GET("/wallet/reference", s.checkToken(), s.GetWalletAddress())
	s.Router.GET("/wallet/transaction_log", s.checkToken(), s.GetUserTransactions())
	s.Router.POST("/wallet/request_withdraw", s.checkToken(), s.AddCashOut())
	s.Router.GET("/wallet/return_withdrawal", s.checkToken(), s.ReturnCashOut())
	// Gifts
	s.Router.GET("/gift/:gift_id/:hero_name", s.checkToken(), s.GiftHandler())
	// Bet arena join
	s.Router.GET("/bet/arena_join/:leader1/:leader2/:arena_type", s.checkToken(), s.TrinityCoreBridgeServer.JoinNewArena())
}
