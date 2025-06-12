package main

import (
    "fmt"
    "github.com/tidwall/gjson"
    "github.com/xiaoxuan6/github-mirror/handlers"
    "io/ioutil"
    "net/http"
    "strings"
    "sync"
)

func main() {
    response, err := http.DefaultClient.Get("https://git.mxg.pub/api/github/list")
    if err != nil {
        panic(err)
    }
    defer response.Body.Close()

    var (
        wg   sync.WaitGroup
        mu   sync.Mutex
        urls []string
    )
    b, _ := ioutil.ReadAll(response.Body)
    data := gjson.GetBytes(b, "data").Array()
    for _, item := range data {
        url := gjson.Get(item.String(), "url").String()

        wg.Add(1)
        go func(url string) {
            defer wg.Done()
            re, err := http.DefaultClient.Get(url)
            if err != nil {
                println(fmt.Sprintf("url [%s] fetch failed", url))
                return
            }

            defer re.Body.Close()
            if re.StatusCode != 200 {
                println(fmt.Sprintf("url [%s] status not 200", url))
                return
            }

            b, _ = ioutil.ReadAll(re.Body)
            if ok := strings.Contains(string(b), "href=\"https://github.com/hunshcn/gh-proxy\">hunshcn/gh-proxy</a>"); !ok {
                println(fmt.Sprintf("url [%s] not support github proxy", url))
                return
            }

            mu.Lock()
            urls = append(urls, strings.TrimLeft(strings.TrimLeft(url, "http://"), "https://"))
            mu.Unlock()
        }(url)
    }

    wg.Wait()
    handlers.ActionHandler.Save(urls)
}
