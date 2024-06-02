package api

import (
    "encoding/json"
    "fmt"
    "github.com/xiaoxuan6/github-mirror/redis"
    "net/http"
    "os"
    "strings"
)

type Response struct {
    Code int      `json:"code"`
    Msg  string   `json:"msg"`
    Data []string `json:"data"`
}

func Api(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "get",
    })

    var response *Response
    res, err := kv.Get()
    if err != nil {
        response = errors(err)
    } else {
        result := strings.Split(res.Result, ",")

        var urls []string
        for _, val := range result {
            urls = append(urls, fmt.Sprintf("http://%s", val))
        }

        response = success(urls)
    }

    b, _ := json.Marshal(response)
    _, _ = w.Write(b)
}

func success(data []string) *Response {
    return &Response{
        Code: 200,
        Msg:  "ok",
        Data: data,
    }
}

func errors(err error) *Response {
    return &Response{
        Code: 500,
        Msg:  err.Error(),
        Data: nil,
    }
}
