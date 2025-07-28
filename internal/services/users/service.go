package users

import "context"

type Service struct {
	repo ReadWriter
}

func NewService(repo ReadWriter) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUserStatus(ctx context.Context, userID int64) (string, error) {
	return s.repo.GetUserStatus(ctx, userID)
}

func (s *Service) CreateUser(ctx context.Context, userID int64) error {
	return s.repo.CreateUser(ctx, userID)
}
func (s *Service) ChangeUserStatus(ctx context.Context, userID int64, status string) error {
	return s.repo.ChangeUserStatus(ctx, userID, status)
}
