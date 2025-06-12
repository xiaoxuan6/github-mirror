package handlers

import (
    "github.com/joho/godotenv"
    "github.com/thoas/go-funk"
    "strings"
)

var ActionHandler = new(actionHandler)

type actionHandler struct{}

func (actionHandler) Save(urls []string) {
    _ = godotenv.Load()

    response := RedisHandler.GetWithBody()
    if response == "" {
        return
    }

    responseUrls := strings.Split(response, ",")
    for _, url := range urls {
        responseUrls = append(responseUrls, url)
    }

    uniqUrls := funk.UniqString(responseUrls)
    uniqUrl := strings.Join(uniqUrls, ",")
    RedisHandler.Set(uniqUrl)
    println("upgrade successfully")
}
