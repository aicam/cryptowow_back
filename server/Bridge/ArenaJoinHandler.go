package Bridge

import (
	"github.com/gin-gonic/gin"
	"gorm.io/plugin/dbresolver"
)

func (s *Server) JoinNewArena() gin.HandlerFunc {
	return func(c *gin.Context) {
		leader1 := c.Param("leader1")
		leader2 := c.Param("leader2")
		arenaType := c.Param("arena_type")
		arenaTeam1 := ArenaTeam{}
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT * FROM arena_team WHERE ")
	}
}
