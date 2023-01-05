package ShopService

type BuyItemRequest struct {
	ItemType   string `json:"item_type"`
	ItemID     string `json:"mount_id"`
	HeroName   string `json:"hero_name"`
	CurrencyID string `json:"currency_id"`
}

type Hero struct {
	AccountID  int    `json:"account_id" gorm:"column:account"`
	HeroID     int    `json:"hero_id" gorm:"column:guid"`
	Name       string `json:"name"`
	Race       uint   `json:"race"`
	Gender     bool   `json:"gender"`
	Level      int    `json:"level"`
	Class      int    `json:"class"`
	Online     bool   `json:"online"`
	Money      int    `json:"money"`
	TotalTime  int    `json:"total_time" gorm:"column:totaltime"`
	TotalKills int    `json:"total_kills" gorm:"column:totalKills"`
}

type CharacterInventory struct {
	Bag    int `json:"bag" gorm:"column:bag"`
	Slot   int `json:"slot" gorm:"column:slot"`
	ItemId int `json:"itemId" gorm:"column:item"`
}

type BagRow struct {
	Bag  int
	Slot int
}
