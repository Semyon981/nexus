package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"

	"github.com/Semyon981/nexus/gateway/auth"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	c              userspb.UserServiceClient
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

type AuthClaims struct {
	jwt.RegisteredClaims
}

func NewAuthUseCase(
	c userspb.UserServiceClient,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		c:              c,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: tokenTTLSeconds,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, number, password, name, lastname string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	_, err := a.c.CreateUser(ctx, &userspb.CreateUserRequest{Number: number, Password: password, Name: name})
	return err
}

func (a *AuthUseCase) SignIn(ctx context.Context, number, password string) (string, error) {

	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	resp, err := a.c.AuthUser(ctx, &userspb.AuthUserRequest{Number: number, Password: password})

	if err != nil {
		return "", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(resp.IdUsers, 10),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (int64, error) {

	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		Id_users, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			return 0, err
		}
		return Id_users, nil
	}

	return 0, auth.ErrInvalidAccessToken
}
