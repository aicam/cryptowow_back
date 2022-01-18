package server

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

type HeroInfo struct {
	ID             int                      `json:"id" gorm:"column:guid"`
	Name           string                   `json:"name"`
	Race           uint                     `json:"race"`
	Gender         bool                     `json:"gender"`
	Level          int                      `json:"level"`
	Class          int                      `json:"class"`
	EquipmentCache string                   `json:"equipment_cache" gorm:"column:equipmentCache"`
	Achievements   []string                 `json:"achievements" gorm:"column:achievement"`
	Reputations    []map[string]interface{} `json:"reputations"`
	Mounts         []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"mounts"`
	Pets []string `json:"pets"`
}

type MountsInfo struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Database string `json:"database"`
	SqlExpr  string `json:"sql_expr"`
	Data     []struct {
		ID      string `json:"entry"`
		Name    string `json:"name"`
		SpellID string `json:"spellid_2"`
	}
}

type CompanionsInfo struct {
	SqlExpr string `json:"sql_expr"`
	Data    []struct {
		ID      string `json:"entry"`
		SpellID string `json:"spellid_2"`
	}
}

type HeroPosition struct {
	Map       int     `json:"map"`
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
var DruidClassNumber int = 11

// Level Up equipped items needed for each class
var HordeWeapons = []string{"23469 23468 23466 18826 18844 18835 18840", "23467 18868 18871 18877 18874 18828 16345 18831 23464"}
var AllianceWeapons = "18873 18876 18873 23454 18825 18873 18873 23454 18843 18873 23451 18865 23455 18876 18825"
var LevelUpGift = map[int][]string{
	1:  {"40826 42041 40866 38 40789 35161 40847 35146 ", "51364 40807 35131 51358 42133 41588 42081 18830"},
	2:  {"40828 42043 40869 45 40788 40978 40849 40979", "40984 40808 42118 42119 42134 41587 42079 18876 42853"},
	3:  {"41157 51355 41217 148 41087 35151 41205 35136 51352", "41143 51358 37927 51377 42134 51346 18830 18833 2102"},
	4:  {"41672 51355 41683 49 41650 41833 41655 35137", "51370 41767 51358 51336 42133 41587 51354 23456 12584 18833"},
	5:  {"41854 51349 41869 53 41859 41882 41864", "49183 51339 41874 51336 35129 42137 41587 51346 18873"},
	6:  {"40827 51355 40868 40787 40883 40848 40884 40890 40809", "51358 33919 42137 41587 51354 18869 42621 38145 38145 38145 38145"},
	7:  {"41013 51349 41038 23345 40992 41052 41027 41056", "51373 41001 42119 42118 51377 42135 51348 23454 18825 42598"},
	8:  {"41943 42045 41962 6096 41949 16818 41956 41904", "41894 41968 42119 42118 42133 51377 42076 18873 17966"},
	9:  {"41993 42047 42011 6097 41998 41882 42005", "41904 51329 42017 42119 42118 42136 41590 42077 18873"},
	11: {"41321 51349 41275 41310 41631 41298 41622 41841 41287", "51336 51358 51377 42137 51354 18873 42589"},
}
