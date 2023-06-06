package service

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, number, password, name, lastname string) error
	GetUser(ctx context.Context, number, password string) (int64, error)
}
