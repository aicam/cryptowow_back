package Bridge

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Server struct {
	DB      *gorm.DB
	Redis   *redis.Client
	Context context.Context
}
