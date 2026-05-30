package main

import (
	"log"

	"github.com/joho/godotenv"

	"yozatune_support_bot/internal/bot"
	"yozatune_support_bot/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if err := bot.Run(cfg); err != nil {
		log.Fatalf("bot: %v", err)
	}
}
