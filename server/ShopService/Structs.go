package ShopService

type BuyItemRequest struct {
	ItemType   string `json:"item_type"`
	ItemID     string `json:"mount_id"`
	HeroName   string `json:"hero_name"`
	CurrencyID string `json:"currency_id"`
}
