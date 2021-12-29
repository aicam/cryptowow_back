package server

import (
	"encoding/hex"
	"errors"
	"github.com/aicam/cryptowow_back/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode int    `json:"status"`
	Body       string `json:"body"`
}

func (s *Server) ReturnHeroInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		heroName := context.Param("hero_name")
		var heroInfo HeroInfo
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT guid, name, race, gender, level, class, equipmentCache FROM characters WHERE name='" + heroName + "'").Scan(&heroInfo)
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT achievement FROM character_achievement WHERE guid=" + strconv.Itoa(heroInfo.ID)).Find(&heroInfo.Achievements)
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT faction, standing FROM character_reputation WHERE guid='" + strconv.Itoa(heroInfo.ID) +
			"' AND faction in (1106, 1052, 1090, 1098, 1156, 1073, 1119, 1091)").Find(&heroInfo.Reputations)
		var heroSpells []string
		s.DB.Clauses(dbresolver.Use("characters")).Raw("SELECT spell FROM character_spell WHERE guid=" + strconv.Itoa(heroInfo.ID)).Find(&heroSpells)
		for _, mount := range s.WoWInfo.Mounts.Data {
			if stringInSlice(mount.SpellID, heroSpells) {
				heroInfo.Mounts = append(heroInfo.Mounts, struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}{ID: mount.ID, Name: mount.Name})
			}
		}
		for _, companion := range s.WoWInfo.Companions.Data {
			if stringInSlice(companion.SpellID, heroSpells) {
				heroInfo.Pets = append(heroInfo.Pets, companion.ID)
			}
		}
		//heroInfo.Achievements = string(achievements.Achievement)
		context.JSON(http.StatusOK, heroInfo)
	}
}

func (s *Server) AvailableWallets() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, struct {
			Wallets []string `json:"wallets"`
		}{Wallets: []string{"Trust Wallet", "Bitpay"}})
	}
}

func (s *Server) ReturnUserInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		username := context.GetHeader("username")
		var id int
		s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account WHERE username='" + strings.ToUpper(username) + "'").Scan(&id)
		var heros []Hero
		s.DB.Raw("SELECT name, race, gender, level, class FROM characters WHERE account=" + strconv.Itoa(id)).Scan(&heros)

		var sellingHeros []database.SellingHeros
		s.DB.Where(&database.SellingHeros{Username: username}).Find(&sellingHeros)

		var gifts []database.Gifts
		s.DB.Where(&database.Gifts{
			Username: username,
		}).Find(&gifts)

		currencies := WalletCurrencies()
		var wallets []database.Wallet
		s.DB.Where(&database.Wallet{Name: username}).Find(&wallets)
		context.JSON(http.StatusOK, struct {
			Heros        []Hero                  `json:"heros"`
			Gifts        []database.Gifts        `json:"gifts"`
			Wallets      []database.Wallet       `json:"wallets"`
			Currencies   []string                `json:"currencies"`
			SellingHeros []database.SellingHeros `json:"selling_heros"`
			Username     string                  `json:"username"`
		}{Heros: heros, Gifts: gifts, Wallets: wallets, Currencies: currencies, SellingHeros: sellingHeros, Username: username})
	}
}
func (s *Server) AddUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		var newUser database.UsersData
		var existUser database.UsersData
		_ = context.BindJSON(&newUser)
		if err := s.DB.Where(&database.UsersData{Username: newUser.Username}).Find(&existUser).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusOK, Response{
				StatusCode: 0,
				Body:       "Username exist",
			})
			return
		}
		gift := database.Gifts{
			Username:     newUser.Username,
			Description:  "Level up first hero free!",
			Action:       "lvlup",
			Condition:    "Register",
			Used:         false,
			UsedHeroName: "",
		}
		s.DB.Save(&gift)
		newUser.Password = MD5(newUser.Password)
		log.Println(newUser.Password)
		s.DB.Save(&newUser)

		// Uncomment
		//GMReqs.CreateAccount(newUser.Username, newUser.Password)

		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Added",
		})
	}
}

func (s *Server) BuyHero() gin.HandlerFunc {
	return func(context *gin.Context) {
		username := context.GetHeader("username")
		var buyingHero struct {
			HeroInfo         database.SellingHeros `json:"hero_info"`
			SelectedCurrency string                `json:"selected_currency"`
		}
		var testBuyingHero database.SellingHeros
		err := context.BindJSON(&buyingHero)
		if err != nil || buyingHero.HeroInfo.HeroID == 0 {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid request",
			})
			return
		}
		err = s.DB.Where(&buyingHero.HeroInfo).First(&testBuyingHero).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusOK, Response{
				StatusCode: -3,
				Body:       "Malicious activity detected",
			})
			return
		}
		var id int
		s.DB.Clauses(dbresolver.Use("auth")).Raw("SELECT id from account WHERE username='" + strings.ToUpper(username) + "'").Scan(&id)
		err = SetBuyHeroTransaction(username, buyingHero.SelectedCurrency, testBuyingHero.Price, s.DB)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		s.DB.Clauses(dbresolver.Use("characters")).Raw("UPDATE characters SET account=" + strconv.Itoa(id) + " WHERE guid=" +
			strconv.Itoa(testBuyingHero.HeroID))
		s.DB.Delete(testBuyingHero)
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Hero transferred to " + username + " account",
		})
	}
}

func (s *Server) GetToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user database.UsersData
		err := context.BindJSON(&user)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}
		key := []byte("Ali@Kian")
		if err := s.DB.Where(database.UsersData{Username: user.Username,
			Password: MD5(user.Password)}).First(&user).Error; err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       "Invalid credentials",
			})
			return
		}

		token, err := DesEncrypt([]byte(user.Username), key)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       hex.EncodeToString(token),
		})
	}
}
