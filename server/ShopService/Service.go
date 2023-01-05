package ShopService

import (
	"github.com/aicam/cryptowow_back/monitoring"
	"github.com/aicam/cryptowow_back/server/GlobalStructs"
	"gorm.io/gorm"
	"strconv"
)

type Service struct {
	DB       *gorm.DB
	PP       monitoring.PrometheusParams
	BagsInfo GlobalStructs.BagsInfo
	BagItems map[int]int
}

func NewService(DB *gorm.DB, PP monitoring.PrometheusParams, bagsInfo GlobalStructs.BagsInfo) Service {
	bagItems := make(map[int]int)
	for _, item := range bagsInfo.Data {
		id, _ := strconv.Atoi(item.ID)
		slots, _ := strconv.Atoi(item.Slots)
		bagItems[id] = slots
	}
	return Service{
		DB:       DB,
		PP:       PP,
		BagsInfo: bagsInfo,
		BagItems: bagItems,
	}
}
