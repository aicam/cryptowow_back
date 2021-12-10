package database

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	// Name of username is same as name of wallet
	Name       string `json:"name"`
	CurrencyID string `json:"currency_id"`
	Amount     int    `json:"amount"`
}
