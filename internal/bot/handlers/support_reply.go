package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) SupportReply(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ticket, err := h.tickets.FindBySupportMsgID(dbCtx, msg.ReplyToMessage.MessageId)
	if err != nil {
		return fmt.Errorf("find ticket: %w", err)
	}
	if ticket == nil {
		return nil // ответ не на отслеживаемое сообщение
	}

	_, err = b.CopyMessage(ticket.UserChatID, msg.Chat.Id, msg.MessageId, nil)
	if err != nil {
		return fmt.Errorf("copy reply to user: %w", err)
	}

	return nil
}
