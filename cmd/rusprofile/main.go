package main

import (
	"log"
	"zentotem/config"
	"zentotem/internal/app"
)

func main() {
	cfg := config.GetConfig()
	app := app.NewApp()
	if err := app.Run(cfg); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
