package main

import (
	"account/impl"
	"github.com/Alvs0/actuator/engine"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
)

const (
	ConfigPath = "./config"
)

type Config struct {
	StartAtPort string
	MySQLConfig engine.SqlConfig
}

func main() {
	e := echo.New()

	var cfg Config
	err := engine.LoadConfig("dev", ConfigPath, &cfg)
	if err != nil {
		log.Fatal("[Service] failed to load config cause: ", err.Error())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	sqlAdapter := engine.NewSqlAdapter(cfg.MySQLConfig)
	accountQuery := impl.NewAccountQuery(sqlAdapter)

	authHandler := impl.NewAuthHandler(accountQuery)

	e.POST("/login", authHandler.Login)

	r := e.Group("/validate")
	r.Use(impl.CreateMiddleware())
	r.POST("", authHandler.Validate)

	e.Logger.Fatal(e.Start(cfg.StartAtPort))
}
