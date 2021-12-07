package server

func (s *Server) Routes() {
	s.Router.POST("/login", s.checkToken(), s.GetToken())
	s.Router.GET("/get_available_wallets", s.checkToken(), s.AvailableWallets())
	s.Router.POST("/add_info", s.checkToken(), s.AddInfo())
	s.Router.GET("/get_info/:offset", s.checkToken(), s.GetInfo())
	s.Router.POST("/register", s.checkToken(), s.AddUser())
}
