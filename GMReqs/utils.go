package GMReqs

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func makeRequest(cm string) {
	client := &http.Client{}
	file, err := ioutil.ReadFile("GMReqs/TC_SOAP")
	if err != nil {
		log.Println(err)
	}
	reqBody := strings.Replace(string(file), "GMCommand", cm, 1)
	//log.Println(reqBody)
	req, err := http.NewRequest("POST", "http://127.0.0.1:7878", bytes.NewBufferString(reqBody))
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth("test1", "test1")
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

func CreateAccount(username, password string) {
	cm := "account create " + username + " " + password
	makeRequest(cm)
}
