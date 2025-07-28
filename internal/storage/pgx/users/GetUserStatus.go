package pgxusers

import "context"

func (r *Repository) GetUserStatus(ctx context.Context, userID int64) (string, error) {
	row, err := r.db.QueryContext(ctx, "SELECT status FROM users WHERE tg_id = $1", userID)
	if err != nil {
		return "", err
	}
	defer row.Close()
	var status string
	err = row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil

}
