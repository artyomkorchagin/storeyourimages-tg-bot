package users

import "context"

type Reader interface {
	GetUserStatus(ctx context.Context, userID string) (string, error)
}

type Writer interface {
	CreateUser(ctx context.Context, userID string) error
	ChangeUserStatus(ctx context.Context, userID, status string) error
}

type ReadWriter interface {
	Reader
	Writer
}
