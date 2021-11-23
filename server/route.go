package server

func (s *Server) Routes() {
	s.Router.POST("/login", s.checkToken(), s.GetToken())
	s.Router.POST("/add_info", s.checkToken(), s.AddInfo())
	s.Router.GET("/get_info/:offset", s.checkToken(), s.GetInfo())
	s.Router.GET("/add_user/:username", s.checkToken(), s.AddUser())
}
