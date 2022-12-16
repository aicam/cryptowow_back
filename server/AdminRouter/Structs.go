package AdminRouter

type Response struct {
	StatusCode int         `json:"status"`
	Body       interface{} `json:"body"`
}

type AddBalanceToAccountRequest struct {
	CurrencyID string  `json:"currency_id"`
	Amount     float64 `json:"amount"`
	Username   string  `json:"username"`
}
