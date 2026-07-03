package repository

import (
	"context"
	"encoding/json"

	"github.com/Xenios7/Trade-executor/internal/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
    client *redis.Client //object that connects to and talks to Redis.
}

func NewRedisRepository(c *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: c,
	}
}

func (r *RedisRepository) Cache(order domain.Order) error {
    key := "order:" + order.ID
    ctx := context.Background()
    
    value, err := json.Marshal(order)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, value, 0).Err()
}
func (r *RedisRepository) GetCached(id string) (domain.Order, error) {
    key := "order:" + id
    ctx := context.Background()
    
    value, err := r.client.Get(ctx, key).Result()
    if err != nil {
        return domain.Order{}, err
    }
    
    var order domain.Order
    err = json.Unmarshal([]byte(value), &order)
    return order, err
}


