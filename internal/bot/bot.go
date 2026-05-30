package bot

import (
	"context"
	"log"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	gotghandlers "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"

	"yozatune_support_bot/internal/config"
	"yozatune_support_bot/internal/db"
	"yozatune_support_bot/internal/bot/handlers"
)

func Run(cfg *config.Config) error {
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		return err
	}

	b, err := gotgbot.NewBot(cfg.BotToken, nil)
	if err != nil {
		return err
	}
	log.Printf("Bot started: @%s (ID: %d)", b.Username, b.Id)

	tickets := db.NewTicketRepo(pool)
	h := handlers.New(cfg, tickets)

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Printf("handler error: %v", err)
			return ext.DispatcherActionNoop
		},
	})

	dispatcher.AddHandler(gotghandlers.NewCommand("start", h.Start))
	dispatcher.AddHandler(gotghandlers.NewCommand("done", h.Done))
	dispatcher.AddHandler(gotghandlers.NewMessage(filters.Message(h.IsPrivateMessage), h.UserMessage))
	dispatcher.AddHandler(gotghandlers.NewMessage(filters.Message(h.IsSupportReply), h.SupportReply))

	go runCleanup(tickets)

	updater := ext.NewUpdater(dispatcher, nil)

	if cfg.WebhookURL != "" {
		if err := updater.StartWebhook(b, cfg.WebhookPath, ext.WebhookOpts{
			ListenAddr:  cfg.WebhookListenAddr,
			SecretToken: cfg.WebhookSecretToken,
		}); err != nil {
			return err
		}
		if _, err := b.SetWebhook(cfg.WebhookURL+cfg.WebhookPath, &gotgbot.SetWebhookOpts{
			SecretToken:        cfg.WebhookSecretToken,
			DropPendingUpdates: true,
		}); err != nil {
			return err
		}
		log.Printf("Webhook: %s%s (listen %s)", cfg.WebhookURL, cfg.WebhookPath, cfg.WebhookListenAddr)
	} else {
		if _, err := b.DeleteWebhook(&gotgbot.DeleteWebhookOpts{DropPendingUpdates: true}); err != nil {
			log.Printf("delete webhook: %v", err)
		}
		if err := updater.StartPolling(b, nil); err != nil {
			return err
		}
		log.Printf("Polling mode started")
	}

	updater.Idle()
	return nil
}

func runCleanup(tickets *db.TicketRepo) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		n, err := tickets.DeleteOlderThan(context.Background(), 30)
		if err != nil {
			log.Printf("cleanup: %v", err)
		} else if n > 0 {
			log.Printf("cleanup: deleted %d old tickets", n)
		}
	}
}
