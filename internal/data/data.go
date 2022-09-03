package data

import (
	"valuation/internal/conf"
	"valuation/pkg/storage"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v9"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGoodRepo, storage.NewDB, storage.NewRedis)

// Data .
type Data struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		redis.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, redis: redis}, cleanup, nil
}
