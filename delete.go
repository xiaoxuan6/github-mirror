package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/handlers"
    "strings"
)

var input string

func main() {
    println("请输入需要删除的链接地址：")
    _, _ = fmt.Scanln(&input)

    _ = godotenv.Load()
    url := handlers.RedisHandler.GetWithBody()
    urls := strings.Split(url, ",")

    newUrls := funk.FilterString(urls, func(s string) bool {
        return strings.Compare(s, strings.TrimSpace(input)) != 0
    })
    newUrl := strings.Join(newUrls, ",")
    handlers.RedisHandler.Set(newUrl)
    println("delete rul successfully")
}
