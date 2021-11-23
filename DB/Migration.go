package DB

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
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
	Username string `json:"username"`
	Password string `json:"password"`
}

func DbSqlMigration(url string) *gorm.DB {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&WebData{})
	db.AutoMigrate(&UsersData{})
	return db
}
