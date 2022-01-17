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
	"strconv"
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

// Cash directory to store webpages
var CASHPATH = "test/cash/"
var PORT = 44297

func main() {
	router := gin.Default()

	// This handler will match all requests
	router.GET("/*url", func(c *gin.Context) {
		url := c.Param("url")
		url = url[1:]
		// File name can not contain specific characters so we hash the url to provide appropriate file name
		cashFile := MD5(url) + ".txt"

		// Bypass cors
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max, username")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if _, err := os.Stat(CASHPATH + cashFile); !errors.Is(err, os.ErrNotExist) {
			f, _ := os.ReadFile(CASHPATH + cashFile)
			c.String(http.StatusOK, string(f))
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
			c.String(http.StatusBadGateway, "Error in sending request")
			return
		}
		//Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			c.String(http.StatusBadGateway, "Error in sending request")
			return
		}
		//Convert the body to type string
		sb := string(body)
		f, err := os.Create(CASHPATH + cashFile)
		checkErr(err)
		_, err = f.Write(body)
		checkErr(err)

		c.String(http.StatusOK, sb)

	})
	router.Run(":" + strconv.Itoa(PORT))
}
