package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/abhaybhu10/login/model"
	"github.com/go-redis/redis/v8"
)

type UserRedis struct {
	client  *redis.Client
	timeout time.Duration
}

func (r *UserRedis) Save(ctx context.Context, user model.User) error {

	data, err := json.Marshal(user)
	if err != nil {
		fmt.Errorf("Error while marshaling")
	}
	err = r.client.Set(ctx, user.ID, data, r.timeout).Err()
	return err
}

func (r *UserRedis) Get(ctx context.Context, key string) (*model.User, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New(fmt.Sprintf("Key %s does not exit", key))
	}
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserRedis() *UserRedis {

	URL := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:         URL,
		Password:     "",
		DB:           1,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    &tls.Config{},
	})
	return &UserRedis{
		client: rdb,
	}
}
