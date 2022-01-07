package server

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := []byte("Ali@Kian")
		token, err := hex.DecodeString(c.GetHeader("Authorization"))
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		username, err := DesDecrypt(token, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		c.Request.Header.Set("username", string(username))
		//c.Header()
		c.Next()
	}
}
