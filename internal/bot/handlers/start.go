package handlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) Start(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b,
		"Привет! Напишите ваш вопрос, и служба поддержки ответит вам как можно скорее.",
		nil,
	)
	return err
}
