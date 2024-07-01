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

    client := redis.NewClient()

    uri := r.RequestURI
    uri = strings.Trim(uri, "/")
    if ok := strings.Compare(uri, "api/urls"); ok == 0 {
        w.Header().Set("Content-Type", "application/json;charset=utf-8")

        var response *Response
        res, err := client.Get(os.Getenv("key"))
        if err != nil {
            response = errors(err)
        } else {
            result := strings.Split(res, ",")

            var urls []string
            for _, val := range result {
                urls = append(urls, fmt.Sprintf("http://%s", val))
            }

            response = success(urls)
        }

        output(response, w)
        return
    }

    if ok := strings.Compare(uri, "api/url/save"); ok == 0 {
        w.Header().Set("Content-Type", "application/json;charset=utf-8")

        type RequestBody struct {
            Url      string `json:"url"`
            Response string `json:"response"`
        }
        var requestBody RequestBody
        _ = json.NewDecoder(r.Body).Decode(&requestBody)
        logrus.Info("请求参数 requestBody.url：", requestBody)

        if len(requestBody.Response) < 1 {
            response := errors(errors2.New("验证失败"))
            output(response, w)
            return
        }

        if ok := turnstileVerify(requestBody.Response); !ok {
            response := errors(errors2.New("验证失败"))
            output(response, w)
            return
        }

        res, err := http.Get(requestBody.Url)
        defer res.Body.Close()
        if err != nil {
            response := errors(errors2.New(fmt.Sprintf("url fail: %s", err.Error())))
            output(response, w)
            return
        }

        if res.StatusCode != 200 {
            response := errors(errors2.New(fmt.Sprintf("url fail status code: %s", res.Status)))
            output(response, w)
            return
        }

        body, _ := ioutil.ReadAll(res.Body)
        if stat := strings.Contains(string(body), "gh-proxy"); !stat {
            response := errors(errors2.New(fmt.Sprintf("url [%s] not support github proxy", requestBody.Url)))
            output(response, w)
            return
        }

        url := strings.TrimLeft(strings.TrimLeft(requestBody.Url, "http://"), "https://")
        url = strings.Trim(url, "/")

        result, err := client.Get(os.Getenv("key"))
        if err != nil {
            output(errors(errors2.New("kv get value fail")), w)
            return
        }

        urls := funk.UniqString(strings.Split(result, ","))
        if !funk.ContainsString(urls, url) {
            urls = append(urls, url)
        } else {
            output(errors(errors2.New(fmt.Sprintf("url [%s] is exists", url))), w)
            return
        }

        _ = client.Set(os.Getenv("key"), strings.Join(urls, ","))
        output(success(nil), w)
        return
    }

    if ok := strings.HasPrefix(uri, "https:/github.com"); ok == false {
        _, _ = w.Write([]byte("The URL prefix must be https://github.com"))
        return
    }

    res, err := client.Get(os.Getenv("key"))
    if err != nil {
        _, _ = w.Write([]byte(err.Error()))
        return
    }

    result := strings.Split(res, ",")
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

func turnstileVerify(response string) bool {
    body := `{"secret":"` + os.Getenv("secret") + `", "response":"` + response + `"}`
    re, _ := http.NewRequest(http.MethodPost, "https://challenges.cloudflare.com/turnstile/v0/siteverify", strings.NewReader(body))
    re.Header.Set("Content-Type", "application/json")
    res, err := http.DefaultClient.Do(re)
    defer res.Body.Close()
    if err != nil {
        logrus.Error("http client fail: " + err.Error())
        return false
    }

    result := struct {
        Success     bool     `json:"success"`
        ErrorCodes  []string `json:"error-codes"`
        ChallengeTs string   `json:"challenge_ts"`
        Hostname    string   `json:"hostname"`
    }{}

    b, _ := ioutil.ReadAll(res.Body)
    _ = json.Unmarshal(b, &result)
    if result.Success != true {
        return false
    }

    return true
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

func output(response *Response, w http.ResponseWriter) {
    b, _ := json.Marshal(response)
    _, _ = w.Write(b)
}
