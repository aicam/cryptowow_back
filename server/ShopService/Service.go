package ShopService

import (
	"github.com/aicam/cryptowow_back/monitoring"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
	PP monitoring.PrometheusParams
}
