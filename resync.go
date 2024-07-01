package main

import (
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
        logrus.Error("not fail urls")
        return
    }

    client := redis.NewClient()
    result, err := client.Get(os.Getenv("key"))
    if err != nil {
        logrus.Error("kv get value fail: ", err.Error())
        return
    }

    urls := strings.Split(result, ",")
    diffUrls, _ := funk.Difference(urls, failUrls)
    diffs := diffUrls.([]string)
    if len(diffs) > 0 {
        body := strings.Join(diffs, ",")
        _ = client.Set(os.Getenv("key"), body)
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
