package ArenaService

import (
	"context"
	"github.com/aicam/cryptowow_back/prometheus"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	DB      *gorm.DB
	Rdb     *redis.Client
	Context context.Context
	PP      prometheus.PrometheusParams
}

var READYCHECKCOUNTER = 800
