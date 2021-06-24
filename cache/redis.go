package cache

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/ptflp/infoblog-server/config"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg config.Redis) (*RedisCache, error) {
	r := &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: "",
			DB:       0,
		}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := r.client.Ping(ctx).Err()

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RedisCache) Get(key string, ptrValue interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return Deserialize(b, ptrValue)
}

func (r *RedisCache) Set(key string, ptrValue interface{}, expires time.Duration) {
	b, err := Serialize(ptrValue)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = r.client.Set(ctx, key, b, expires).Err(); err != nil {
		return
	}
}
