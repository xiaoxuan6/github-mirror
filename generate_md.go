package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/xiaoxuan6/github-mirror/handlers"
    "io/ioutil"
    "os"
    "strings"
)

func main() {
    _ = godotenv.Load()

    res, err := handlers.RedisHandler.Get()
    if err != nil {
        fmt.Printf("kv 获取失败：%s", err.Error())
        return
    }

    result := res.Data
    fmt.Println("urls count: ", len(res.Data))

    var content strings.Builder
    for _, val := range result {
        uri := fmt.Sprintf("http://%s", val)
        content.WriteString(fmt.Sprintf("[%s](%s)\n", uri, uri))
    }

    err = ioutil.WriteFile("./Links.md", []byte(content.String()), os.ModePerm)
    if err != nil {
        fmt.Printf("Error writing to file: %s\n", err.Error())
    }

    fmt.Println("generate md done.")
}
