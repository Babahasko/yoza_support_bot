package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) Done(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	if msg.Chat.Id != h.cfg.SupportChatID {
		return nil
	}
	if msg.ReplyToMessage == nil {
		_, err := msg.Reply(b, "Используйте /done как ответ на сообщение с вопросом.", nil)
		return err
	}

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleted, err := h.tickets.Delete(dbCtx, msg.ReplyToMessage.MessageId)
	if err != nil {
		return fmt.Errorf("delete ticket: %w", err)
	}

	var text string
	if deleted {
		text = "✅ Тикет закрыт."
	} else {
		text = "Тикет не найден — возможно, уже закрыт."
	}

	_, err = msg.Reply(b, text, nil)
	return err
}
