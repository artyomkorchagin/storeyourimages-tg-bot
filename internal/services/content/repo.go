package content

import (
	"context"

	"github.com/artyomkorchagin/storeyourimages/internal/types"
)

type Reader interface {
	ReadContent(ctx context.Context, userID string, offset int) ([]string, error)
}

type Writer interface {
	WriteContent(ctx context.Context, wdr *types.WriteDataRequest) error
}

type ReadWriter interface {
	Reader
	Writer
}
