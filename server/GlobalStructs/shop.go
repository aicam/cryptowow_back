package GlobalStructs

type BagsInfo struct {
	Data []struct {
		ID    string `json:"entry"`
		Slots string `json:"ContainerSlots"`
	} `json:"data"`
	BackPackStart int `json:"backPackStart"`
	BackPackEnd   int `json:"backPackEnd"`
	BagSlotsStart int `json:"bagSlotsStart"`
	BagSlotsEnd   int `json:"bagSlotsEnd"`
}
