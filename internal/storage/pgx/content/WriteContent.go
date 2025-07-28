package pgxcontent

import (
	"context"

	"github.com/artyomkorchagin/storeyourimages/internal/types"
)

func (r *Repository) WriteContent(ctx context.Context, wdr *types.WriteDataRequest) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO content_data (tg_id, filepath, type) VALUES ($1, $2, $3)`,
		wdr.UserID, wdr.Filepath, wdr.Datatype)
	return err
}
