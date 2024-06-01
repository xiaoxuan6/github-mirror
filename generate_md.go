package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    _ = godotenv.Load()
    kv := redis.NewKvClient(&redis.Option{
        Key:    os.Getenv("key"),
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Action: "get",
    })

    res, err := kv.Get()
    if err != nil {
        fmt.Printf("kv 获取失败：%s", err.Error())
        return
    }

    result := strings.Split(res.Result, ",")
    fmt.Println("urls count: ", len(result))

    if len(result) > 0 {
        dir, _ := os.Getwd()
        filename := filepath.Join(dir, "Links.md")
        f, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModePerm)
        for _, val := range result {
            uri := fmt.Sprintf("https://%s", val)
            _, _ = f.WriteString(fmt.Sprintf("[%s](%s)\n", uri, uri))
        }
    }

    fmt.Println("generate md done.")
}
