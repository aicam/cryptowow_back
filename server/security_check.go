package server

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) checkToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		key := []byte("Ali@Kian")
		token, err := hex.DecodeString(context.GetHeader("Authorization"))
		username := []byte(context.GetHeader("username"))
		if len(token) == 0 || len(username) == 0 {
			context.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		if err != nil {
			context.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		destext, err := DesDecrypt(token, key)
		if err != nil || string(destext) != string(username) {
			context.AbortWithStatusJSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Authorization failed",
			})
			return
		}
		context.Next()
	}
}
