package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
)

func main() {
    _ = godotenv.Load()

    filename := ".github/src/urls.txt"
    if _, err := os.Stat(filename); err != nil {
        logrus.Error(fmt.Sprintf("file [%s] not exits", filename))
        return
    }

    b, err := os.ReadFile(filename)
    if err != nil {
        logrus.Error(fmt.Sprintf("read fail %s", err.Error()))
        return
    }

    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "set",
    })

    res, err := kv.Set(string(b))
    if err != nil {
        logrus.Error(fmt.Sprintf("kv set fail: %s", err.Error()))
        return
    }

    logrus.Info(fmt.Sprintf("kv set %s", res.Result))
}
