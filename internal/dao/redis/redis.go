package redis

import (
	"github.com/go-redis/redis/v8"
)

type Client struct {
	db        *redis.Client
	keyPrefix string
}

func New(rdb *redis.Client, keyPrefix string) *Client {
	return &Client{
		db:        rdb,
		keyPrefix: keyPrefix,
	}
}
