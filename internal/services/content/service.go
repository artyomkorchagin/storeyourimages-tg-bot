package content

import (
	"context"

	"github.com/artyomkorchagin/storeyourimages/internal/types"
)

type Service struct {
	repo ReadWriter
}

func NewService(repo ReadWriter) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ReadContent(ctx context.Context, userID string, offset int) ([]string, error) {
	return s.repo.ReadContent(ctx, userID, offset)
}

func (s *Service) WriteContent(ctx context.Context, wdr *types.WriteDataRequest) error {
	return s.repo.WriteContent(ctx, wdr)
}
