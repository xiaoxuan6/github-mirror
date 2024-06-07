package main

import (
    "bufio"
    "fmt"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/redis"
    "io"
    "os"
    "strings"
)

func main() {
    _ = godotenv.Load()

    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "get",
    })

    res, err := kv.Get()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    b, _ := os.ReadFile("site.txt")

    urls := make([]string, 0, 1000)
    br := bufio.NewReader(strings.NewReader(string(b)))
    for {
        a, _, c := br.ReadLine()
        if c == io.EOF {
            break
        }

        uri := strings.TrimRight(string(a), "/")
        uri = strings.TrimLeft(uri, "http://")
        uri = strings.TrimLeft(uri, "https://")

        urls = append(urls, uri)
    }

    uniqUrls := funk.UniqString(urls)
    urlStr := strings.Join(uniqUrls, ",")

    newUrls := fmt.Sprintf("%s,%s", res.Result, urlStr)
    kv = redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "set",
    })

    response, err := kv.Set(newUrls)
    if err != nil {
        logrus.Error("kv set value fail: ", err.Error())
        return
    }

    logrus.Info("kv set ", response.Result)
}
