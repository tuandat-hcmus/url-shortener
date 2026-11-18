package main

import (
	"github.com/tuandat-hcmus/url-shortener/internal/api"
	"github.com/tuandat-hcmus/url-shortener/pkg/redis"
)

// @title URL Shortener API
// @version 1.0
// @description This is a URL Shortener API server.
// @BasePath /

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	app := api.NewEngine(cfg, redisClient)
	err = app.Start()
	if err != nil {
		panic(err)
	}
}
