package server

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func (s *Server) GetServerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Open our jsonFile
		jsonFile, _ := os.Open("database/index.json")
		defer jsonFile.Close()
		b, _ := ioutil.ReadAll(jsonFile)
		c.String(http.StatusOK, string(b))
	}
}
