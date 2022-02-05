package test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestReqArgs struct {
	AuthorizationToken string `json:"authorization_token"`
	Requests           []struct {
		ReqType string `json:"req_type"`
		URL     string `json:"url"`
		Body    string `json:"body"`
		Auth    bool   `json:"auth"`
	}
}

type Response struct {
	StatusCode int         `json:"status"`
	Body       interface{} `json:"body"`
}

func sendGetReq(t *testing.T, url string, auth string) Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	var jsResp Response
	err = json.Unmarshal(body, &jsResp)

	if err != nil {
		t.Error(err)
	}
	return jsResp
}
