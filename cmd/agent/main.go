package main

import (
	"github.com/katenester/Distributed_computing/internal/agent/app"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initalization config %s", err.Error())
	}
	app.Run()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
