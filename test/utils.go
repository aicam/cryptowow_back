package test

type Response struct {
	StatusCode int         `json:"status"`
	Body       interface{} `json:"body"`
}
