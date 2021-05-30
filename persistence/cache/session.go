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

type SessionRedis struct {
	client  *redis.Client
	timeout time.Duration
}

func (r *SessionRedis) Save(ctx context.Context, session model.Session) error {

	data, err := json.Marshal(session)
	if err != nil {
		fmt.Errorf("Error while marshaling")
	}
	err = r.client.Set(ctx, session.SessionID, data, r.timeout).Err()
	return err
}

func (r *SessionRedis) Get(ctx context.Context, key string) (*model.Session, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New(fmt.Sprintf("Key %s does not exit", key))
	}
	if err != nil {
		return nil, err
	}

	var session model.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func GetSessionRedis() *SessionRedis {
	HOST, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		fmt.Println("REDIS_HOST not set")
	}
	PORT := os.Getenv("REDIS_PORT")
	url := fmt.Sprintf("%s:%s", HOST, PORT)
	rdb := redis.NewClient(&redis.Options{
		Addr:         url,
		Password:     "",
		DB:           0, // use default DB
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	pong, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		fmt.Printf("Error while connecting to Redis URL: %s, error: %s\n", url, err.Error())
	}

	fmt.Printf("Redis ping response %s\n", pong)

	return &SessionRedis{
		client:  rdb,
		timeout: 20 * time.Second,
	}
}
