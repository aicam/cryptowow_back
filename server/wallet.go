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

func SetBuyHeroTransaction(username, heroName, selectedCurrency, price string, DB *gorm.DB) error {
	var wallets []database.Wallet
	DB.Where(&database.Wallet{Name: username}).Find(&wallets)
	prices := strings.Split(price, "&")
	log.Print(prices)
	var priceVal float64
	for _, price := range prices {
		if strings.Split(price, "-")[1] == selectedCurrency {
			priceVal, _ = strconv.ParseFloat(strings.Split(price, "-")[0], 32)
		}
	}
	log.Println(priceVal)
	var connectedWallet database.Wallet
	for i := 0; i < len(wallets); i++ {
		if wallets[i].CurrencyID == selectedCurrency {
			connectedWallet = wallets[i]
		}
	}
	if (database.Wallet{}) == connectedWallet {
		return errors.New("Not enough balance")
	}
	if connectedWallet.Amount < priceVal {
		return errors.New("Not enough balance")
	}
	connectedWallet.Amount -= priceVal
	DB.Save(&connectedWallet)
	return nil
}
