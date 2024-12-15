package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

// Rdb redis client
var (
	Rdb *redis.Client
	Rs  *redsync.Redsync
)

// InitRedis 初始化redis连接
//
//	@author 鹿鹿鹿鲨
//	@update 2023-09-27 12:36:02
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.GetString("redis.host"), config.Config.GetInt("redis.port")),
		Password: config.Config.GetString("redis.password"),
		DB:       config.Config.GetInt("redis.db"),
		PoolSize: config.Config.GetInt("redis.poolSize"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lo.Must(Rdb.Ping(ctx).Result())

	log.Logger().Info("Redis Connected!")
	initRedsync()
}

func initRedsync() {
	pool := goredis.NewPool(Rdb)
	Rs = redsync.New(pool)
	log.Logger().Info("Init Redsync Success!")
}

func GetSettingNumber(key string) int {
	atoi, err := strconv.Atoi(Rdb.Get(context.Background(), key).Val())
	if err != nil {
		log.Logger().Error("[GetSettingNumber]", zap.Error(err))
		return 0
	}
	return atoi
}

func SetSettingNumber(key string, num int) error {
	return Rdb.Set(context.Background(), key, num, 0).Err()
}
