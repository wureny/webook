package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

/*
	func InitRedis() redis.Cmdable {
		redisClient := redis.NewClient(&redis.Options{
			Addr: config.Config.Redis.Addr,
		})
		return redisClient
	}
*/
func InitRedis() redis.Cmdable {
	addr := viper.GetString("redis.addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return redisClient
}
