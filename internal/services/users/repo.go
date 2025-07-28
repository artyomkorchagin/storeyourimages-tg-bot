package users

import "context"

type Reader interface {
	GetUserStatus(ctx context.Context, userID int64) (string, error)
}

type Writer interface {
	CreateUser(ctx context.Context, userID int64) error
	ChangeUserStatus(ctx context.Context, userID int64, status string) error
}

type ReadWriter interface {
	Reader
	Writer
}
