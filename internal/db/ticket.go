package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Ticket struct {
	ID           int64
	UserID       int64
	UserChatID   int64
	SupportMsgID int64
}

type TicketRepo struct {
	pool *pgxpool.Pool
}

func NewTicketRepo(pool *pgxpool.Pool) *TicketRepo {
	return &TicketRepo{pool: pool}
}

func (r *TicketRepo) Save(ctx context.Context, userID, userChatID, supportMsgID int64) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO tickets (user_id, user_chat_id, support_msg_id) VALUES ($1, $2, $3)`,
		userID, userChatID, supportMsgID,
	)
	return err
}

func (r *TicketRepo) Delete(ctx context.Context, supportMsgID int64) (bool, error) {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM tickets WHERE support_msg_id = $1`,
		supportMsgID,
	)
	return tag.RowsAffected() > 0, err
}

func (r *TicketRepo) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM tickets WHERE created_at < NOW() - $1 * interval '1 day'`,
		days,
	)
	return tag.RowsAffected(), err
}

func (r *TicketRepo) FindBySupportMsgID(ctx context.Context, supportMsgID int64) (*Ticket, error) {
	t := &Ticket{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, user_chat_id, support_msg_id FROM tickets WHERE support_msg_id = $1`,
		supportMsgID,
	).Scan(&t.ID, &t.UserID, &t.UserChatID, &t.SupportMsgID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}
