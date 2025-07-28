package pgxusers

import "context"

func (r *Repository) ChangeUserStatus(ctx context.Context, userID, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET status = $1 WHERE tg_id = $2`, status, userID)
	return err
}
