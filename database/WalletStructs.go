package database

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	// Name of username is same as name of wallet
	Name       string  `json:"name"`
	CurrencyID string  `json:"currency_id"`
	Amount     float64 `json:"amount"`
}

type Transaction struct {
	gorm.Model
	// same as username
	WalletName string `json:"wallet_name"`
	CurrencyID string `json:"currency_id"`
	Amount     int    `json:"amount"`
	HashCode   string `json:"hash_code"`
	Used       int    `json:"used"`
}
