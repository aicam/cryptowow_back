package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestWoWDB(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("wowdb_route.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var testReqArgs TestReqArgs
	err = json.Unmarshal(byteValue, &testReqArgs)
	if err != nil {
		t.Error(err)
	}
	for _, req := range testReqArgs.Requests {
		if req.ReqType == "GET" {
			var auth = ""
			if req.Auth {
				auth = testReqArgs.AuthorizationToken
			}
			respJS := sendGETReq(t, req.URL, auth)
			if respJS.StatusCode != 1 {
				log.Print(req.URL + " has wrong status code")
				t.Error(respJS.Body)
			}
		}
	}
}
