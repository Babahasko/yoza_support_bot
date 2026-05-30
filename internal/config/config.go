package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	BotToken      string
	DatabaseURL   string
	SupportChatID int64

	// Webhook-режим; если WEBHOOK_URL не задан — используется polling
	WebhookURL         string
	WebhookPath        string
	WebhookListenAddr  string
	WebhookSecretToken string
}

func Load() (*Config, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN is not set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	supportChatStr := os.Getenv("SUPPORT_CHAT_ID")
	if supportChatStr == "" {
		return nil, fmt.Errorf("SUPPORT_CHAT_ID is not set")
	}
	supportChatID, err := strconv.ParseInt(supportChatStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("SUPPORT_CHAT_ID must be an integer: %w", err)
	}

	webhookPath := os.Getenv("WEBHOOK_PATH")
	if webhookPath == "" {
		webhookPath = "/webhook"
	}
	webhookListenAddr := os.Getenv("WEBHOOK_LISTEN_ADDR")
	if webhookListenAddr == "" {
		webhookListenAddr = ":8080"
	}

	return &Config{
		BotToken:           token,
		DatabaseURL:        dbURL,
		SupportChatID:      supportChatID,
		WebhookURL:         os.Getenv("WEBHOOK_URL"),
		WebhookPath:        webhookPath,
		WebhookListenAddr:  webhookListenAddr,
		WebhookSecretToken: os.Getenv("WEBHOOK_SECRET_TOKEN"),
	}, nil
}
