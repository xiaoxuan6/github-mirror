package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/sirupsen/logrus"
    "github.com/thoas/go-funk"
    "github.com/xiaoxuan6/github-mirror/redis"
    "os"
    "strings"
)

func main() {
    _ = godotenv.Load()

    kv := redis.NewKvClient(&redis.Option{
        Token:  os.Getenv("KV_REST_API_TOKEN"),
        Key:    os.Getenv("key"),
        Action: "get",
    })

    res, err := kv.Get()
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    l := strings.Split(res.Result, ",")
    urls := funk.UniqString(l)
    fmt.Println(len(urls))

    uris := make([]string, 0, 100)
    funk.ForEach(urls, func(url string) {
        if ok := funk.ContainsString([]string{
            //"autodiscover.xzpan.top",
            //"bt.moran233.xyz",
            //"gh.xsage.cn",
            //"103.133.178.32",
            //"8.218.216.196",
            //"103.177.248.124:12345",
            //"42.192.251.73",
            //"149.104.15.159:3000",
            //"149.104.30.123:9000",
            //"154.19.242.10",
            //"43.132.131.30",
            //"43.154.99.97",
            //"123.58.210.249",
            //"42.192.251.73",
            //"43.132.131.30",
            //"154.19.242.10",
            "ghproxy.lktuchaung.buzz",
        }, url); ok == false {
            uris = append(uris, url)
        }
    })

    fmt.Println(len(uris))
    b := strings.Join(uris, ",")

    //b, _ := os.ReadFile("site.txt")
    //
    //urls := make([]string, 0, 1000)
    //br := bufio.NewReader(strings.NewReader(string(b)))
    //for {
    //    a, _, c := br.ReadLine()
    //    if c == io.EOF {
    //        break
    //    }
    //
    //    uri := strings.TrimRight(string(a), "/")
    //    uri = strings.TrimLeft(uri, "http://")
    //    uri = strings.TrimLeft(uri, "https://")
    //
    //    urls = append(urls, uri)
    //}
    //
    //uniqUrls := funk.UniqString(urls)
    //urlStr := strings.Join(uniqUrls, ",")
    //
    //newUrls := fmt.Sprintf("%s,%s", res.Result, urlStr)
    kv = redis.NewKvClient(&redis.Option{
       Token:  os.Getenv("KV_REST_API_TOKEN"),
       Key:    os.Getenv("key"),
       Action: "set",
    })

    response, err := kv.Set(b)
    if err != nil {
       logrus.Error("kv set value fail: ", err.Error())
       return
    }

    logrus.Info("kv set ", response.Result)
}
