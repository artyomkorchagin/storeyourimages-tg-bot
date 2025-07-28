package pgxusers

import "context"

func (r *Repository) CreateUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (tg_id) VALUES ($1)", userID)
	return err
}
