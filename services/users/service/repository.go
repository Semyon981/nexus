package service

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, number, password, name, lastname string) error
	GetUserAuth(ctx context.Context, number, password string) (int64, error)
	GetUserId(ctx context.Context, Id_users int64) (int64, error)
}
