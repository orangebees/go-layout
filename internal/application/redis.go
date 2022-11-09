package application

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

var redisClient *redis.Client
var keyPrefix = "app"

func GetRedis() *redis.Client {
	if redisClient != nil {
		return redisClient
	}
	return createRedis()
}
func GetRedisKeyPrefix() string {
	return keyPrefix
}
func createRedis() *redis.Client {
	if v := os.Getenv("REDIS_KEYPREFIX"); v != "" {
		keyPrefix = v
	}
	rdsopt := &redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	}
	if v := os.Getenv("REDIS_NETWORK"); v != "" {
		rdsopt.Network = v
	}
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		rdsopt.Addr = v
	}
	if v := os.Getenv("REDIS_USERNAME"); v != "" {
		rdsopt.Username = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		rdsopt.Password = v
	}
	rc := redis.NewClient(rdsopt)
	if err := rc.Ping(context.TODO()).Err(); err != nil {
		panic(err)
		return nil
	}
	go func() {
		for {
			time.Sleep(time.Second * 2)
			if err := rc.Ping(context.TODO()).Err(); err != nil {
				panic(err)
			}
		}
	}()
	GetLogger().Info().Msg("redis connect success")
	//注册释放方法
	closes = append(closes, rc.Close)
	return rc
}
