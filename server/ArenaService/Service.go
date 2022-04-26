package ArenaService

import (
	"context"
	"github.com/aicam/cryptowow_back/monitoring"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	DB      *gorm.DB
	Rdb     *redis.Client
	Context context.Context
	PP      monitoring.PrometheusParams
}

var READYCHECKCOUNTER = 800
