package api

import (
    "encoding/json"
    errors2 "errors"
    "fmt"
    "github.com/sirupsen/logrus"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/redis"
    "io/ioutil"
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
    //_ = godotenv.Load()

    uri := r.RequestURI
    uri = strings.Trim(uri, "/")
    if ok := strings.Compare(uri, "api/urls"); ok == 0 {
        w.Header().Set("Content-Type", "application/json;charset=utf-8")

        var response *Response
        res, err := fetchData()
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
        return
    }

    if ok := strings.Compare(uri, "api/url/save"); ok == 0 {
        w.Header().Set("Content-Type", "application/json;charset=utf-8")

        type RequestBody struct {
            Url string `json:"url"`
        }
        var requestBody RequestBody
        _ = json.NewDecoder(r.Body).Decode(&requestBody)
        logrus.Info("请求参数 requestBody.url：", requestBody.Url)

        res, err := http.Get(requestBody.Url)
        if err != nil {
            response := errors(errors2.New(fmt.Sprintf("url fail: %s", err.Error())))
            b, _ := json.Marshal(response)
            _, _ = w.Write(b)
            return
        }

        if res.StatusCode != 200 {
            response := errors(errors2.New(fmt.Sprintf("url fail status code: %s", res.Status)))
            b, _ := json.Marshal(response)
            _, _ = w.Write(b)
            return
        }

        url := strings.TrimLeft(strings.TrimLeft(requestBody.Url, "http://"), "https://")
        url = strings.Trim(url, "/")

        response := save(url)
        b, _ := json.Marshal(response)
        _, _ = w.Write(b)
        return
    }

    fmt.Println(uri)
    if ok := strings.HasPrefix(uri, "https:/github.com"); ok == false {
        _, _ = w.Write([]byte("The URL prefix must be https://github.com"))
        return
    }

    res, err := fetchData()
    if err != nil {
        _, _ = w.Write([]byte(err.Error()))
        return
    }

    result := strings.Split(res.Result, ",")
    i := funk.RandomInt(0, len(result)-1)
    proxy := result[i]
    if len(proxy) == 0 {
        proxy = "https://mirror.ghproxy.com"
    }

    newUri := fmt.Sprintf("http://%s/%s", proxy, strings.ReplaceAll(uri, "https:/", "https://"))
    response, err := http.Get(newUri)
    defer response.Body.Close()
    if err != nil {
        _, _ = w.Write([]byte(err.Error()))
        return
    }

    b, _ := ioutil.ReadAll(response.Body)
    _, _ = w.Write(b)
}

func save(uri string) *Response {
    res, err := fetchData()
    if err != nil {
        return errors(err)
    }

    l := strings.Split(res.Result, ",")
    urls := funk.UniqString(l)

    stat := true
    funk.ForEach(urls, func(url string) {
        if ok := strings.Compare(url, uri); ok == 0 {
            stat = false
            return
        }
    })

    if stat == true {
        urls = append(urls, uri)
    }

    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "set",
    })

    _, err = kv.Set(strings.Join(urls, ","))
    if err != nil {
        return errors(err)
    }

    data := make([]string, 0)
    return success(data)
}

func fetchData() (redis.Response, error) {
    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "get",
    })

    return kv.Get()
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
