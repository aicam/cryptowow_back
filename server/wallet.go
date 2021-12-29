package server

import (
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func WalletCurrencies() []string {
	return []string{"Ethereum", "CWT"}
}

func SetBuyHeroTransaction(username, vendorUsername, selectedCurrency, price string, DB *gorm.DB) error {
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
	var vendorWallet database.Wallet
	err := DB.Where(&database.Wallet{Name: vendorUsername, CurrencyID: connectedWallet.CurrencyID}).First(&vendorWallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		vendorWallet.CurrencyID = connectedWallet.CurrencyID
		vendorWallet.Name = vendorUsername
		vendorWallet.Amount = 0.0
	}
	vendorWallet.Amount += priceVal
	DB.Save(&vendorWallet)
	DB.Save(&connectedWallet)
	return nil
}
