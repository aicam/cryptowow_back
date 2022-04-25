package WalletService

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func WalletCurrencies() []string {
	return []string{"Ethereum", "CWT"}
}

func AddBalance(username, currencyID string, amount float64, DB *gorm.DB) {
	var wallet database.Wallet
	err := DB.Where(&database.Wallet{Name: username, CurrencyID: currencyID}).First(&wallet).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		wallet.CurrencyID = currencyID
		wallet.Name = username
		wallet.Amount = 0.0
	}
	wallet.Amount += amount
	DB.Save(&wallet)
}

func ReduceBalance(isBet bool, username, currencyID string, amount float64, DB *gorm.DB) error {
	var wallets []database.Wallet
	DB.Where(&database.Wallet{Name: username}).Find(&wallets)
	var wallet database.Wallet
	for i := 0; i < len(wallets); i++ {
		if wallets[i].CurrencyID == currencyID {
			wallet = wallets[i]
		}
	}

	if (database.Wallet{}) == wallet {
		return errors.New("Not enough balance")
	}
	if wallet.Amount < amount {
		return errors.New("Not enough balance")
	}

	// check arena debt before allow reducing
	if !isBet {
		totalBetDebt := GetArenaBetTotalDebt(DB, username)
		if wallet.Amount < amount+totalBetDebt[wallet.CurrencyID] {
			return errors.New("You can not use arena bet money")
		}
	}

	wallet.Amount -= amount
	DB.Save(wallet)
	return nil
}

func GetAccountBalance(username, selectedCurrency string, DB *gorm.DB) float64 {
	var wallets []database.Wallet
	DB.Where(&database.Wallet{Name: username}).Find(&wallets)
	for _, wallet := range wallets {
		if wallet.CurrencyID == selectedCurrency {
			return wallet.Amount
		}
	}
	return 0
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
	err := ReduceBalance(false, username, selectedCurrency, priceVal, DB)
	if err != nil {
		return err
	}
	AddBalance(vendorUsername, selectedCurrency, priceVal, DB)
	return nil
}

func HashTransactionToken(txHash string) string {
	h := sha1.New()
	h.Write([]byte(txHash))
	bs := h.Sum(nil)
	return b64.StdEncoding.EncodeToString(bs)[2:8]
}
