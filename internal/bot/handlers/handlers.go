package handlers

import (
	"yozatune_support_bot/internal/config"
	"yozatune_support_bot/internal/db"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Handler struct {
	cfg     *config.Config
	tickets *db.TicketRepo
}

func New(cfg *config.Config, tickets *db.TicketRepo) *Handler {
	return &Handler{cfg: cfg, tickets: tickets}
}

func (h *Handler) IsPrivateMessage(msg *gotgbot.Message) bool {
	return msg.Chat.Type == "private"
}

func (h *Handler) IsSupportReply(msg *gotgbot.Message) bool {
	return msg.Chat.Id == h.cfg.SupportChatID &&
		msg.ReplyToMessage != nil &&
		msg.From != nil &&
		!msg.From.IsBot
}
