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
	URL := os.Getenv("REDIS_URL")
	rdb := redis.NewClient(&redis.Options{
		Addr:         URL,
		Password:     "",
		DB:           0, // use default DB
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    &tls.Config{},
	})
	return &SessionRedis{
		client: rdb,
	}
}
