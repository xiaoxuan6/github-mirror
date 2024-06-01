package api

import "net/http"

func Api(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json;charset=utf-8")

    _, _ = w.Write([]byte(`{"code":200, "msg": "ok"}`))
}
