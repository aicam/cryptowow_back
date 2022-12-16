package server

import (
	"encoding/hex"
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"github.com/aicam/cryptowow_back/server/LogService"
	"github.com/aicam/cryptowow_back/server/WalletService"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (s *Server) AvailableWallets() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       []string{"Trust Wallet", "Bitpay"},
		})
	}
}

func (s *Server) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var id int
		s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account WHERE username='" + strings.ToUpper(username) + "'").Scan(&id)
		var heros []Hero
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT guid, name, race, gender, level, class FROM characters WHERE account=" + strconv.Itoa(id)).Scan(&heros)

		var sellingHeros []database.SellingHeros
		s.DB.Where(&database.SellingHeros{Username: username}).Find(&sellingHeros)

		var gifts []database.Gifts
		s.DB.Where(&database.Gifts{
			Username: username,
		}).Find(&gifts)

		currencies := WalletService.WalletCurrencies()
		var wallets []database.Wallet
		s.DB.Where(&database.Wallet{Name: username}).Find(&wallets)

		for i, hero := range heros {
			arenaTeams := []HeroArenaTeams{}
			s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT `arenaTeamId`, " +
				"`name`, `type`, `rating`, `seasonGames`," +
				" `seasonWins`, `weekGames`, `weekWins`, `rank`, `backgroundColor` FROM arena_team WHERE `captainGuid`=" + strconv.Itoa(hero.HeroID)).Find(&arenaTeams)
			heros[i].ArenaTeamsCaptain = arenaTeams
		}

		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body: struct {
				Heros        []Hero                  `json:"heros"`
				Gifts        []database.Gifts        `json:"gifts"`
				Wallets      []database.Wallet       `json:"wallets"`
				Currencies   []string                `json:"currencies"`
				SellingHeros []database.SellingHeros `json:"selling_heros"`
				Username     string                  `json:"username"`
			}{Heros: heros, Gifts: gifts, Wallets: wallets, Currencies: currencies, SellingHeros: sellingHeros, Username: username},
		})
	}
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* check request time */
		start := time.Now()
		/* check request time */

		var newUser database.UsersData
		var existUser database.UsersData
		_ = c.BindJSON(&newUser)

		ip := c.ClientIP()
		csrfHeader := c.GetHeader("X-CSRF-Token")
		log.Println(ip)
		var ipTrack database.IPRecords
		iperr := s.DB.Where(&database.IPRecords{IPAddress: ip}).First(&ipTrack).Error
		if errors.Is(iperr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, Response{
				StatusCode: -13,
				Body:       "Malicious activity detected!",
			})
			return
		}
		if ipTrack.Info != "" {
			if time.Now().Add(-30*time.Minute).Before(ipTrack.UpdatedAt) && ipTrack.Checked == 1 {
				c.JSON(http.StatusOK, Response{
					StatusCode: -1,
					Body:       "Too many requests",
				})
				return
			}
			if ipTrack.Info != csrfHeader {
				c.JSON(http.StatusOK, Response{
					StatusCode: -8,
					Body:       "Malicious activity detected",
				})
				return
			}

		}
		newUser.Username = strings.ToUpper(newUser.Username)
		if err := s.DB.Where(&database.UsersData{Username: newUser.Username}).First(&existUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, Response{
				StatusCode: 0,
				Body:       "Username exist",
			})
			return
		}

		gift := database.Gifts{
			Username:     newUser.Username,
			Description:  "The hero you choose for this gift will get all bg set by mail and level ups to 80",
			Action:       "For Signing Up",
			Condition:    "Register",
			Used:         false,
			UsedHeroName: "",
			GiftID:       1,
		}
		s.DB.Save(&gift)
		log.Println(newUser.Password)

		// Uncomment
		//err := GMReqs.CreateAccount(newUser.Username, newUser.Password)
		var err error
		err = nil
		if err != nil {
			c.JSON(http.StatusBadGateway, Response{
				StatusCode: -1,
				Body:       "Game server is currently closed!, try again later",
			})
		} else {
			newUser.Password = MD5(newUser.Password)
			s.DB.Save(&newUser)
			ipTrack.Checked = 1
			s.DB.Save(&ipTrack)
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				Body:       "Added",
			})
		}
		/* check request time */
		end := time.Now()
		/* check request time */
		s.PP.Histograms["registration_process_time"].Observe(float64(end.Sub(start).Nanoseconds()))
	}
}

func (s *Server) BuyHero() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var buyingHero struct {
			HeroInfo         database.SellingHeros `json:"hero_info"`
			SelectedCurrency string                `json:"selected_currency"`
		}
		var testBuyingHero database.SellingHeros
		err := c.BindJSON(&buyingHero)
		if err != nil || buyingHero.HeroInfo.HeroID == 0 {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid request",
			})
			return
		}
		err = s.DB.Where(&buyingHero.HeroInfo).First(&testBuyingHero).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, Response{
				StatusCode: -3,
				Body:       "Malicious activity detected",
			})
			return
		}
		var id int
		s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account WHERE username='" + strings.ToUpper(username) + "'").Scan(&id)
		err = WalletService.SetBuyHeroTransaction(username, buyingHero.HeroInfo.Username, buyingHero.SelectedCurrency, testBuyingHero.Price, s.DB)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		err = s.DB.Clauses(dbresolver.Use("characters")).Exec("UPDATE characters SET account=" + strconv.Itoa(id) + " WHERE guid=" +
			strconv.Itoa(testBuyingHero.HeroID)).Error
		if err != nil {
			log.Println("Update account id err: ", err)
		}
		s.DB.Delete(&testBuyingHero)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Hero transferred to " + username + " account",
		})

		s.PP.Gauges["Number_Currently_Selling_Heros"].Add(-1)
	}
}

func (s *Server) ReturnEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var events []database.Events
		s.DB.Find(&events)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       events,
		})
	}
}

func (s *Server) AddTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("username")
		var tx database.TransactionLog
		err := c.BindJSON(&tx)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid Data",
			})
			return
		}
		txHash := strings.ToUpper(WalletService.HashTransactionToken(tx.TransactionHash))
		if txHash != tx.TXHash {
			c.JSON(http.StatusOK, Response{
				StatusCode: -20,
				Body:       "Malicious activity detected! take care of your IP and token",
			})
			LogService.LogPotentialCyberAttack(c, "Add_Transaction")
			return
		}
		tx.Username = username
		s.DB.Save(&tx)
		WalletService.AddBalance(username, tx.CurrencyID, tx.Amount, s.DB)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Transaction Done successfully",
		})

		s.PP.Counters["Total_Successfull_Transactions"].Inc()
	}
}

func (s *Server) GetToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user database.UsersData
		err := c.BindJSON(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}
		key := []byte("Ali@Kian")
		if err := s.DB.Where(database.UsersData{Username: user.Username,
			Password: MD5(user.Password)}).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}

		token, err := DesEncrypt([]byte(user.Username), key)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       hex.EncodeToString(token),
		})
	}
}

func (s *Server) GetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		log.Println(ip)
		// TODO: environment variables
		var ipTrack database.IPRecords
		err := s.DB.Where(&database.IPRecords{IPAddress: ip}).First(&ipTrack).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if time.Now().Add(-30*time.Minute).Before(ipTrack.UpdatedAt) && ipTrack.Checked == 1 {
				c.JSON(http.StatusOK, Response{
					StatusCode: -1,
					Body:       "Too many requests",
				})
				return
			}
			csrfToken := ipTrack.Info
			if ipTrack.Checked == 0 {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					Body:       "Base64 " + csrfToken,
				})
				return
			}
		}
		csrfToken := tokenize("Ali@Kian"+time.Now().String(), ip)
		ipTrack.IPAddress = ip
		ipTrack.Info = csrfToken
		ipTrack.Checked = 0
		ipTrack.Reason = "Registration"
		s.DB.Save(&ipTrack)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Base64 " + csrfToken,
		})
	}
}
