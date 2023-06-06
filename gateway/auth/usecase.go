package auth

import (
	"context"
)

const CtxUserKey = "id_users"

type UseCase interface {
	SignUp(ctx context.Context, number, password, name, lastname string) error
	SignIn(ctx context.Context, number, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
}
