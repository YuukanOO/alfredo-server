package main

import (
	"github.com/YuukanOO/alfredo/env"
	"github.com/YuukanOO/alfredo/handlers"
	"github.com/YuukanOO/alfredo/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	err := env.LoadFromFile("./alfredo.toml")

	if err != nil {
		panic(err)
	}

	defer env.Cleanup()

	err = env.LoadAdapters("./adapters.json")

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(middlewares.CORS(&middlewares.CORSOptions{
		AllowedOrigins: env.Current().Server.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}), middlewares.DB())

	handlers.Register(r)

	panic(r.Run(env.Current().Server.Listen))
}
