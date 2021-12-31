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
	Title        string `json:"title"`
	Description  string `json:"description"`
	Action       string `json:"action"`
	Condition    string `json:"condition"`
	Used         bool   `json:"used"`
	UsedHeroName string `json:"used_hero_name"`
}

type SellingHeros struct {
	gorm.Model
	Username   string `json:"username"`
	HeroName   string `json:"hero_name"`
	HeroID     int    `json:"hero_id"`
	Price      string `json:"price"`
	Class      int    `json:"class"`
	Race       int    `json:"race"`
	Gender     bool   `json:"gender"`
	Level      int    `json:"level"`
	Money      int    `json:"money"`
	TotalTime  int    `json:"total_time"`
	TotalKills int    `json:"total_kills"`
	Note       string `json:"note"`
}

type Events struct {
	gorm.Model
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Gift        string    `json:"gift"`
	Conditions  string    `json:"conditions"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
}

type IPRecords struct {
	gorm.Model
	IPAddress string `json:"ip_address"`
	TrackID   int    `json:"track_id"`
	Reason    string `json:"reason"`
	Info      string `json:"info"`
	Checked   int    `json:"checked"`
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
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&SellingHeros{})
	db.AutoMigrate(&Events{})
	db.AutoMigrate(&IPRecords{})
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
