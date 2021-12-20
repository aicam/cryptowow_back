package server

type Hero struct {
	Name   string `json:"name"`
	Race   uint   `json:"race"`
	Gender bool   `json:"gender"`
	Level  int    `json:"level"`
	Class  int    `json:"class"`
}

type HeroInfo struct {
	ID             int     `json:"id" gorm:"column:guid"`
	Name           string  `json:"name"`
	Race           uint    `json:"race"`
	Gender         bool    `json:"gender"`
	Level          int     `json:"level"`
	Class          int     `json:"class"`
	EquipmentCache string  `json:"equipment_cache" gorm:"column:equipmentCache"`
	Achievements   []uint8 `json:"achievements" gorm:"column:achievement"`
}

type HeroPosition struct {
	Map       uint    `json:"map"`
	PositionX float32 `json:"position_x"`
	PositionY float32 `json:"position_y"`
	PositionZ float32 `json:"position_z"`
}

var Home = struct {
	Alliance HeroPosition
	Horde    HeroPosition
}{
	Alliance: HeroPosition{
		Map:       0,
		PositionX: -8956.84,
		PositionY: 518.406,
		PositionZ: 96.3553,
	},
	Horde: HeroPosition{
		Map:       1,
		PositionX: 1502.71,
		PositionY: -4415.42,
		PositionZ: 21.5512,
	},
}

var WarriorClassNumber int = 1
var PaladinClassNumber int = 2
var HunterClassNumber int = 3
var RogueClassNumber int = 4
var PriestClassNumber int = 5
var DKClassNumber int = 6
var ShamanClassNumber int = 7
var MageClassNumber int = 8
var WarlockClassNumber int = 9
var DruidClassNumber int = 10

// Level Up equipped items needed for each class
var LevelUpGift = map[int]string{

	8: "41943 0 42045 0 41962 0 6096 0 41949 0 16818 0 41956 0 41904 0 41894 0 41968 0 42119 0 42118 0 42133 0 51377 0 42076 0 18873 0 0 0 0 0 0 0 17966 0 0 0 0 0 0 0 "}
