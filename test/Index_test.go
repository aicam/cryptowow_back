package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func sendGetReq(t *testing.T, url string) interface{} {
	resp, err := http.Get(url)

	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	var jsResp interface{}
	err = json.Unmarshal(body, &jsResp)

	if err != nil {
		t.Error(err)
	}
	return jsResp
}

func TestIndexInfo(t *testing.T) {
	indexJS := sendGetReq(t, "http://localhost/server_status")

}
