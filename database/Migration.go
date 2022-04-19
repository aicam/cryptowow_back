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
	GiftID       uint   `json:"gift_id"`
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

type TransactionLog struct {
	gorm.Model
	Username        string  `json:"username"`
	Amount          float64 `json:"amount"`
	CurrencyID      string  `json:"currency_id"`
	Status          bool    `json:"status"`
	BlockHash       string  `json:"block_hash"`
	BlockNumber     int     `json:"block_number"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	TransactionHash string  `json:"transaction_hash"`
	TXHash          string  `json:"tx_hash"`
}

type CashOutRequest struct {
	gorm.Model
	Username      string  `json:"username"`
	Amount        float64 `json:"amount"`
	CurrencyID    string  `json:"currency_id"`
	WalletAddress string  `json:"wallet_address"`
	WalletApp     string  `json:"wallet_app"`
	Note          string  `json:"note"`
	PendingStage  int     `json:"pending_stage"`
	TX            string  `json:"tx"`
}

// Arena Bet system
type QueuedTeams struct {
	gorm.Model
	TeamId        int `json:"team_id"`
	InQueueTeamId int `json:"in_queue_team"`
}

type TeamRequests struct {
	gorm.Model
	TeamId          int `json:"team_id"`
	RequestedTeamId int `json:"requested_team_id"`
}

type TeamReadyGames struct {
	gorm.Model
	InviterTeam int  `json:"team_id"`
	InvitedTeam int  `json:"opponent_id"`
	IsPlayed    bool `json:"is_played"`
}

type BetNotification struct {
	gorm.Model
	TeamId    int    `json:"team_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Seen      bool   `json:"seen"`
	NotifType uint8  `json:"notif_type"`
}

type InGameTeamData struct {
	gorm.Model
	BucketID           uint `json:"bucket_id"`
	InviterSeasonGames uint `json:"inviter_season_games"`
	InviterSeasonWins  uint `json:"inviter_season_wins"`
	InvitedSeasonGames uint `json:"invited_season_games"`
	InvitedSeasonWins  uint `json:"invited_season_wins"`
	Winner             uint `json:"winner"`
}

func DbSqlMigration(url string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&WebData{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&UsersData{})
	db.AutoMigrate(&Gifts{})
	db.AutoMigrate(&Wallet{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&SellingHeros{})
	db.AutoMigrate(&Events{})
	db.AutoMigrate(&IPRecords{})
	db.AutoMigrate(&TransactionLog{})
	db.AutoMigrate(&CashOutRequest{})
	db.AutoMigrate(&QueuedTeams{})
	db.AutoMigrate(&TeamRequests{})
	db.AutoMigrate(&TeamReadyGames{})
	db.AutoMigrate(&BetNotification{})
	db.AutoMigrate(&InGameTeamData{})
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(strings.Replace(url, "messenger_api", "characters", 1))}}, "characters").
		Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open(strings.Replace(url, "messenger_api", "auth", 1))},
		}, "auth"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
