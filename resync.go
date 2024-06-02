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

    failUrls := fetchMdContent()
    if len(failUrls) < 1 {
        logrus.Info("not fail urls")
        return
    }

    urls, err := fetchKvData()
    if err != nil {
        logrus.Error(fmt.Sprintf("获取 kv data 失败：%s", err.Error()))
        return
    }

    diffUrls, _ := funk.Difference(urls, failUrls)
    diffs := diffUrls.([]string)
    if len(diffs) > 0 {
        body := strings.Join(diffs, ",")
        setKvData(body)
    }

    logrus.Info("resync done.")
}

func fetchMdContent() []string {
    b, _ := os.ReadFile("out.md")
    split := strings.Split(string(b), " ")

    var urls []string
    for _, val := range split {
        urls = append(urls, strings.TrimSpace(val))
    }

    return urls
}

func fetchKvData() ([]string, error) {
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

func setKvData(value string) {
    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "set",
    })

    res, err := kv.Set(value)
    if err != nil {
        logrus.Error(fmt.Sprintf("kv set fail: %s", err.Error()))
        return
    }

    logrus.Info(fmt.Sprintf("kv set %s", res.Result))
}
