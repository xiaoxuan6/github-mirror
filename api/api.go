package api

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Code int      `json:"code"`
    Msg  string   `json:"msg"`
    Data []string `json:"data"`
}

func Api(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    response := &Response{
        Code: 200,
        Msg:  "ok",
        Data: []string{
            "http://101.35.42.207:1188",
            "http://101.35.42.207:11881",
            "http://101.35.42.207:11882",
            "http://101.35.42.207:11883",
            "http://101.35.42.207:11884",
        },
    }

    b, _ := json.Marshal(response)
    _, _ = w.Write(b)
}
