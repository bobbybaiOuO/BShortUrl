package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bobbybaiOuO/BShortUrl/config"
	"github.com/bobbybaiOuO/BShortUrl/internal/repo"
	"github.com/redis/go-redis/v9"
)

// SetURL(ctx context.Context, url repo.Url) error

// RedisCache .
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache .
func NewRedisCache(cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Address,
		Password: cfg.PassWord,
		DB: cfg.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

// SetURL .
func (c *RedisCache) SetURL(ctx context.Context, url repo.Url) error {
	data, err := json.Marshal(url)
	if err != nil {
		return err
	}
	
	if err := c.client.Set(ctx, url.ShortCode, data, time.Until(url.ExpiredAt)).Err(); err != nil {
		return err
	}

	return nil
}

// GetURL .
func (c *RedisCache) GetURL(ctx context.Context, shortCode string) (*repo.Url, error) {
	data, err := c.client.Get(ctx, shortCode).Bytes()
	// 查询不出来不算err
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var url repo.Url
	if err := json.Unmarshal(data, &url); err != nil {
		return nil, err
	}
	return &url, nil
}

// Close .
func (c *RedisCache) Close() error {
	return c.client.Close()
}