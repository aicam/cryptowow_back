package GMReqs

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func makeRequest(cm string) error {
	client := &http.Client{}
	file, err := ioutil.ReadFile("GMReqs/TC_SOAP")
	if err != nil {
		log.Println(err)
		return err
	}
	reqBody := strings.Replace(string(file), "GMCommand", cm, 1)
	//log.Println(reqBody)
	req, err := http.NewRequest("POST", "http://127.0.0.1:7878", bytes.NewBufferString(reqBody))
	if err != nil {
		log.Println(err)
		return err
	}
	req.SetBasicAuth("test1", "test1")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
	//log.Println(string(f))
}

func CreateAccount(username, password string) error {
	cm := "account create " + username + " " + password
	return makeRequest(cm)
}

func AddItems(title, body, heroName, items string) {
	cm := "send items " + heroName + " \"" + title + "\" \"" + body + "\" " + items
	makeRequest(cm)
}

func LevelUpHero(heroName string) {
	cm := "character level " + heroName + " 80"
	makeRequest(cm)
}
