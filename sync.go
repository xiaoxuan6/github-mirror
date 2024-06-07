package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
    "strings"
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

    s, _ := fetchKvItem()
    if len(s) > 0 {
        urls := strings.Split(string(b), ",")

        funk.ForEach(urls, func(url string) {
            s = append(s, url)
        })

        newUrls := funk.UniqString(s)
        urlStr := strings.Join(newUrls, ",")
        b = []byte(urlStr)
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

func fetchKvItem() ([]string, error) {
    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "get",
    })

    res, err := kv.Get()
    if err != nil {
        return nil, err
    }

    result := strings.Split(res.Result, ",")
    return result, nil
}
