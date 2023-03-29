package main

import (
	"actuator/engine"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	MySQLConfig engine.SqlConfig
}

func main() {
	echo := echo.New()

	cfg, err := LoadConfig("dev", "./config")
	if err != nil {
		log.Fatal("[Service] failed to load config cause: ", err.Error())
	}

	sqlAdapter := engine.NewSqlAdapter(cfg.MySQLConfig)
	fmt.Println(sqlAdapter)

	echo.Logger.Fatal(echo.Start(":8080"))
}

func LoadConfig(configName, pathToConfigFile string) (config Config, err error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(pathToConfigFile)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
