package Bridge

import (
	"github.com/gin-gonic/gin"
)

type ArenaTeam struct {
}

func (s *Server) JoinNewArena() gin.HandlerFunc {
	return func(c *gin.Context) {
		leader1 := c.Param("leader1")
		leader2 := c.Param("leader2")
		arenaType := c.Param("arena_type")

	}
}
