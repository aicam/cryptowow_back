package Bridge

import (
	"github.com/aicam/cryptowow_back/server"
	"github.com/gin-gonic/gin"
)

type Server server.Server

func (s *Server) JoinNewArena() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
