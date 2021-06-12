package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client  *redis.Client
	timeout time.Duration
}

func GetRedisClient() *RedisClient {
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

	return &RedisClient{
		client:  rdb,
		timeout: 20 * time.Second,
	}

}
func (r *RedisClient) Save(ctx context.Context, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}

	namespacedKey := getNamespacedKey(key, value)
	fmt.Printf("Saving data in cache %s\n", namespacedKey)
	return r.client.Set(ctx, namespacedKey, p, r.timeout).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) (interface{}, error) {
	namespacedKey := getNamespacedKey(key, dest)
	fmt.Printf("Query cache for %s\n", namespacedKey)
	p, err := r.client.Get(ctx, namespacedKey).Result()
	if err == redis.Nil {
		return nil, errors.New(fmt.Sprintf("Key %s does not exist\n", namespacedKey))
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(p), dest)
	if err != nil {
		return nil, err
	}
	fmt.Printf("key %s found cache value %v", key, p)
	return dest, err
}

func getNamespacedKey(key string, dest interface{}) string {
	return key
	//return fmt.Sprintf("%s_%s", key, fmt.Sprintf("%T", dest))
}
