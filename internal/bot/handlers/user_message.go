package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func (h *Handler) UserMessage(b *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	user := ctx.EffectiveUser

	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var supportMsgID int64

	if msg.Text != "" {
		text := fmt.Sprintf("📨 <b>%s</b> (ID: <code>%d</code>):\n\n%s",
			userMention(user), user.Id, msg.Text)
		sent, err := b.SendMessage(h.cfg.SupportChatID, text, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		if err != nil {
			return fmt.Errorf("send to support: %w", err)
		}
		supportMsgID = sent.MessageId
	} else {
		header := fmt.Sprintf("📨 Сообщение от <b>%s</b> (ID: <code>%d</code>):",
			userMention(user), user.Id)
		_, err := b.SendMessage(h.cfg.SupportChatID, header, &gotgbot.SendMessageOpts{
			ParseMode: "HTML",
		})
		if err != nil {
			return fmt.Errorf("send header to support: %w", err)
		}
		fwd, err := b.ForwardMessage(h.cfg.SupportChatID, msg.Chat.Id, msg.MessageId, nil)
		if err != nil {
			return fmt.Errorf("forward to support: %w", err)
		}
		supportMsgID = fwd.MessageId
	}

	if err := h.tickets.Save(dbCtx, user.Id, msg.Chat.Id, supportMsgID); err != nil {
		return fmt.Errorf("save ticket: %w", err)
	}

	return nil
}

func userMention(u *gotgbot.User) string {
	if u.Username != "" {
		return fmt.Sprintf(`<a href="tg://user?id=%d">@%s</a>`, u.Id, u.Username)
	}
	name := u.FirstName
	if u.LastName != "" {
		name += " " + u.LastName
	}
	return fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, u.Id, name)
}
