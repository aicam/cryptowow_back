package database

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"log"
	"strings"
	"time"
)

type WebData struct {
	gorm.Model
	Status      string    `json:"status"`
	Error       string    `json:"error"`
	Username    string    `json:"username"`
	DubaiTxt    string    `json:"dubai_txt"`
	DubaiTime   time.Time `json:"dubai_time"`
	ArmeniaTxt  string    `json:"armenia_txt"`
	ArmeniaTime time.Time `json:"armenia_time"`
	TurkeyTxt   string    `json:"turkey_txt"`
	TurkeyTime  time.Time `json:"turkey_time"`
}

type UsersData struct {
	gorm.Model
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Wallet   string `json:"wallet"`
	WalletID string `json:"wallet_id"`
}

type Gifts struct {
	gorm.Model
	Username     string `json:"username"`
	Description  string `json:"description"`
	Action       string `json:"action"`
	Condition    string `json:"condition"`
	Used         bool   `json:"used"`
	UsedHeroName string `json:"used_hero_name"`
}

func DbSqlMigration(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&WebData{})
	db.AutoMigrate(&UsersData{})
	db.AutoMigrate(&Gifts{})
	db.AutoMigrate(&Wallet{})
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(strings.Replace(url, "messenger_api", "characters", 1))}}, "characters").
		Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open(strings.Replace(url, "messenger_api", "auth", 1))},
		}, "auth"))
	if err != nil {
		log.Println(err)
	}
	return db
}
