package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// https://go.dev/tour/methods/9
type Redis interface {
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

// RedisClient defines a redis client.
type RedisClient struct {
	client *redis.Client
}

// New returns a new RedisClient
func New(port int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%d", port),
	})

	return &RedisClient{
		client: rdb,
	}
}

// Ping checks the connection to Redis
func (rc *RedisClient) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

// definition for set
func (rc *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := rc.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// definition for get
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := rc.client.Get(ctx, key).Result()

	return val, err
}

// RedisClient struct should implement the Redis interface
var _ Redis = &RedisClient{}
