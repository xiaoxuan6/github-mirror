package redis

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

type Option struct {
    Token  string
    Key    string
    Action string
}

type kv struct {
    uri     string
    request func(uri string) (*http.Request, error)
    result  func(req *http.Request) (response Response, err error)
}

func NewKvClient(option *Option) *kv {
    request := func(uri string) (*http.Request, error) {
        request, err := http.NewRequest(http.MethodGet, uri, nil)
        if err != nil {
            return nil, fmt.Errorf("create http request fail: %s", err.Error())
        }

        request.Header.Set("Authorization", "Bearer "+option.Token)

        return request, nil
    }

    response := func(req *http.Request) (response Response, err error) {
        res, err := http.DefaultClient.Do(req)
        defer res.Body.Close()
        if err != nil {
            return response, fmt.Errorf("请求失败：%s", err.Error())
        }

        b, _ := ioutil.ReadAll(res.Body)
        err = json.Unmarshal(b, &response)
        if err != nil {
            return response, fmt.Errorf("json 解析错误：%s", err.Error())
        }

        return
    }

    return &kv{
        uri:     fmt.Sprintf("https://prime-anchovy-41693.upstash.io/%s/%s", option.Action, option.Key),
        request: request,
        result:  response,
    }
}

type Response struct {
    Result string `json:"result"`
}

func (k kv) Set(value string) (response Response, errors error) {
    // url q去除前缀和后缀
    value = strings.Trim(
        strings.ReplaceAll(
            strings.ReplaceAll(value, "http://", ""),
            "http://",
            "",
        ),
        "/",
    )

    req, err := k.request(fmt.Sprintf("%s/%s", k.uri, value))
    if err != nil {
        return response, err
    }

    return k.result(req)
}

func (k kv) Get() (response Response, errors error) {
    req, err := k.request(k.uri)
    if err != nil {
        return response, err
    }

    return k.result(req)
}
