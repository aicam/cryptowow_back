package ArenaService

import (
	"context"
	"github.com/aicam/cryptowow_back/Prometheus"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Service struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Context context.Context
	PP      Prometheus.PrometheusParams
}

var READYCHECKCOUNTER = 800
