package cache

import (
	"context"
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
	fmt.Printf("User %v saved to redis\n", user)
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

	HOST := os.Getenv("REDIS_HOST")
	PORT := os.Getenv("REDIS_PORT")
	url := fmt.Sprintf("%s:%s", HOST, PORT)
	rdb := redis.NewClient(&redis.Options{
		Addr:         url,
		Password:     "",
		DB:           1,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	pong, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		fmt.Printf("Error while connecting to Redis URL: %s, error: %s\n", url, err.Error())
	}

	fmt.Printf("Redis ping response %s\n", pong)

	return &UserRedis{
		client:  rdb,
		timeout: 20 * time.Second,
	}
}
