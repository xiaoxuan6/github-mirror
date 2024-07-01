package main

import (
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/handlers"
    "os"
    "strings"
)

func main() {
    _ = godotenv.Load()

    b, _ := os.ReadFile("out.md")
    failUrls := strings.Split(string(b), " ")
    if len(failUrls) < 1 {
        logrus.Error("not fail urls")
        return
    }

    response, err := handlers.RedisHandler.Get()
    if err != nil {
        logrus.Error("kv get value fail: ", err.Error())
        return
    }

    diffUrls, _ := funk.Difference(response.Data, failUrls)
    diffs := diffUrls.([]string)
    if len(diffs) > 0 {
        body := strings.Join(diffs, ",")
        handlers.RedisHandler.Set(body)
    }

    logrus.Info("resync done.")
}
