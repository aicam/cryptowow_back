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
	s.Router.GET("/wow/get_info", s.checkToken(), s.GetUserInfo())
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
	s.Router.POST("/bet/arena_invite/", s.checkToken(), s.TrinityCoreBridgeService.InviteTeam())
	s.Router.POST("/bet/arena_accept/", s.checkToken(), s.TrinityCoreBridgeService.AcceptInvitation())
	s.Router.POST("/bet/start_game/", s.checkToken(), s.TrinityCoreBridgeService.StartGame())
	s.Router.POST("/bet/accept_start_game/", s.checkToken(), s.TrinityCoreBridgeService.AcceptStartGame())
	s.Router.GET("/bet/get_result/:team_id", s.checkToken(), s.TrinityCoreBridgeService.GetResult())
	// CashCors
	s.Router.GET("/cashcors/:url", s.checkToken(), s.CashCorsReq())
	// Shop routes
	s.Router.POST("/shop/purchase/", s.checkToken(), s.ShopService.BuyItem())
	s.Router.POST("/shop/add_item/", s.checkToken(), s.ShopService.AddItemToShop())
	s.Router.GET("/shop/list", s.checkToken(), s.ShopService.GetShopItems())
	// Admin routes
	// Level 1
	s.AdminRouter.AddRoutes(s.Router, s.checkToken)
}
