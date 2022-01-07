package server

import (
	"crypto/sha1"
	b64 "encoding/base64"
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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

func ReduceBalance(username, currencyID string, amount float64, DB *gorm.DB) error {
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
	wallet.Amount -= amount
	DB.Save(wallet)
	return nil
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
	err := ReduceBalance(username, selectedCurrency, priceVal, DB)
	if err != nil {
		return err
	}
	AddBalance(vendorUsername, selectedCurrency, priceVal, DB)
	return nil
}

func (s *Server) GetWalletAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "0xDE4C72835bcC0041Dd1B446BfD0D85bE346BC0A2",
		})
	}
}

func (s *Server) GetUserTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var txs []database.TransactionLog
		s.DB.Where(&database.TransactionLog{Username: username}).Find(&txs)
		c.JSON(http.StatusOK, txs)
	}
}

func (s *Server) AddCashOut() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var co database.CashOutRequest
		err := c.BindJSON(co)
		if err != nil {
			return
		}
		err = ReduceBalance(username, co.CurrencyID, co.Amount, s.DB)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		co.PendingStage = 0
		s.DB.Save(&co)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Cashout request added successfully",
		})
	}
}

func HashTransactionToken(txHash string) string {
	h := sha1.New()
	h.Write([]byte(txHash))
	bs := h.Sum(nil)
	return b64.StdEncoding.EncodeToString(bs)[2:8]
}
