package pgxcontent

import (
	"context"
	"fmt"
)

func (r *Repository) ReadContent(ctx context.Context, userID string, offset int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT filepath 
        FROM content_data 
        WHERE tg_id = $1 
        ORDER BY created_at DESC 
        LIMIT 10 
        OFFSET $2`, userID, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var filepaths []string
	for rows.Next() {
		var filepath string
		if err := rows.Scan(&filepath); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		filepaths = append(filepaths, filepath)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return filepaths, nil
}
