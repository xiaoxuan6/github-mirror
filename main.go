package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
)

func main() {
    _ = godotenv.Load()

    key := "prefix"

    // kv set
    key = "kv"
    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    key,
        Action: "set",
    })

    res, err := kv.Set("mirror")
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println(res.Result)

    // kv get
    kv = redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    key,
        Action: "get",
    })

    res, err = kv.Get()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Println(res.Result)
}
