package server

import (
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func WalletCurrencies() []string {
	return []string{"Ethereum", "CWT"}
}

func SetBuyHeroTransaction(username, selectedCurrency, price string, DB *gorm.DB) error {
	var wallets []database.Wallet
	DB.Where(&database.Wallet{Name: username}).Find(&wallets)
	prices := strings.Split(price, "&")
	var priceVal float64
	for _, price := range prices {
		if strings.Split(price, "-")[1] == selectedCurrency {
			priceVal, _ = strconv.ParseFloat(strings.Split(price, "-")[0], 32)
		}
	}
	var connectedWallet database.Wallet
	for i := 0; i < len(wallets); i++ {
		log.Println(wallets[i].CurrencyID, selectedCurrency, wallets[i].CurrencyID == selectedCurrency)
		if wallets[i].CurrencyID == selectedCurrency {
			connectedWallet = wallets[i]
		}
	}
	log.Print(connectedWallet)
	if (database.Wallet{}) == connectedWallet {
		return errors.New("Not enough balance")
	}
	if connectedWallet.Amount < priceVal {
		return errors.New("Not enough balance")
	}
	connectedWallet.Amount -= priceVal
	// TODO: add amount to vendor
	DB.Save(&connectedWallet)
	return nil
}
