package main

import (
	"bytes"
	bs64 "encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	client := &http.Client{}
	file, err := ioutil.ReadFile("test/TC_SOAP")
	if err != nil {
		log.Println(err)
	}
	authorization := "Basic " + bs64.StdEncoding.EncodeToString([]byte("ali:ali"))
	log.Println(authorization)
	reqBody := strings.Replace(string(file), "GMCommand", "character level Testm 80", 1)
	//log.Println(reqBody)
	req, err := http.NewRequest("POST", "http://127.0.0.1:7878", bytes.NewBufferString(reqBody))
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth("ali", "ali")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(f))
}
