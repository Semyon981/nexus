package service

import (
	"context"

	"github.com/Semyon981/nexus/services/users/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, number, password, name, lastname string) error
	GetUserAuth(ctx context.Context, number, password string) (int64, error)
	GetUserId(ctx context.Context, Id_users int64) (models.User, error)
}
