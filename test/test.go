package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func checkErr(e error) {
	if e != nil {
		log.Println(e)
	}
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func main() {
	router := gin.Default()

	// This handler will match all requests
	router.GET("/*url", func(c *gin.Context) {
		url := c.Param("url")
		url = url[1:]
		cashFile := MD5(url) + ".txt"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max, username")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if _, err := os.Stat("test/cash/" + cashFile); !errors.Is(err, os.ErrNotExist) {
			f, _ := os.ReadFile("test/cash/" + cashFile)
			c.String(http.StatusOK, string(f))
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
		//Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//Convert the body to type string
		sb := string(body)
		f, err := os.Create("test/cash/" + cashFile)
		checkErr(err)
		_, err = f.Write(body)
		checkErr(err)

		c.String(http.StatusOK, sb)

	})
	router.Run(":4999")
}
