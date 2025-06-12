package handlers

import (
    "errors"
    "fmt"
    "github.com/joho/godotenv"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
    "strings"
)

var (
    client       *redis.Redis
    RedisHandler = new(redisHandler)
)

type redisHandler struct {
}

func init() {
    _ = godotenv.Load()
    client = redis.NewClient()
}

func (kh redisHandler) Get() (*Response, error) {
    result, err := client.Get(os.Getenv("key"))
    if err != nil {
        return nil, err
    }

    return Success(strings.Split(result, ",")), nil
}

func (kh redisHandler) GetWithBody() string {
    result, err := client.Get(os.Getenv("key"))
    if err != nil {
        return ""
    }

    return result
}

func (kh redisHandler) Set(value string) string {
    return client.Set(os.Getenv("key"), value)
}

func (kh redisHandler) SetWithUrl(url string) error {
    result, err := client.Get(os.Getenv("key"))
    if err != nil {
        return errors.New("kv get value fail")
    }

    urls := funk.UniqString(strings.Split(result, ","))
    if funk.ContainsString(urls, url) {
        return errors.New(fmt.Sprintf("url [%s] is exists", url))
    }

    urls = append(urls, url)
    kh.Set(strings.Join(urls, ","))
    return nil
}
