package redis

import (
    "context"
    "crypto/tls"
    "github.com/go-redis/redis/v8"
    "os"
)

type Redis struct {
    client *redis.Client
}

func NewClient() *Redis {
    url := os.Getenv("KV_URL")
    options, err := redis.ParseURL(url)
    if err != nil {
        panic(err.Error())
    }

    options.TLSConfig = &tls.Config{InsecureSkipVerify: true}
    redisClient := redis.NewClient(options)

    return &Redis{client: redisClient}
}

func (r *Redis) Get(key string) (string, error) {
    res, err := r.client.Get(context.Background(), key).Result()
    if err != nil {
        return "", err
    }

    return res, nil
}

func (r *Redis) Set(key, value string) string {
    return r.client.Set(context.Background(), key, value, 0).String()
}
