package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/xiaoxuan6/github-mirror/redis"
    "io/ioutil"
    "os"
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
        var bodys []string
        for _, val := range result {
            uri := fmt.Sprintf("http://%s", val)
            bodys = append(bodys, fmt.Sprintf("[%s](%s)", uri, uri))
        }

        body := strings.Join(bodys, "\n")
        err := ioutil.WriteFile("./Links.md", []byte(body), os.ModePerm)
        if err != nil {
            fmt.Printf("Error writing to file: %s\n", err.Error())
        }
    }

    fmt.Println("generate md done.")
}
