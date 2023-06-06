package server

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/golang-jwt/jwt/v5"
)

type server struct {
	identifierpb.UnimplementedServiceServer
	hashSalt   string
	signingKey []byte
}

func NewServer(hashSalt string, signingKey []byte) *server {
	return &server{hashSalt: hashSalt, signingKey: signingKey}
}

func (s *server) Hash(ctx context.Context, in *identifierpb.HashRequest) (*identifierpb.HashResponse, error) {
	pwd := sha1.New()
	pwd.Write([]byte(in.Str))
	pwd.Write([]byte(s.hashSalt))

	return &identifierpb.HashResponse{Hash: fmt.Sprintf("%x", pwd.Sum(nil))}, nil
}

type AuthClaims struct {
	jwt.RegisteredClaims
}

func (s *server) JwtGen(ctx context.Context, in *identifierpb.JwtGenRequest) (*identifierpb.JwtGenResponse, error) {

	claims := AuthClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(in.Time))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   in.Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString(s.signingKey)
	return &identifierpb.JwtGenResponse{Token: tok}, err
}

func (s *server) JwtParse(ctx context.Context, in *identifierpb.JwtParseRequest) (*identifierpb.JwtParseResponse, error) {
	token, err := jwt.ParseWithClaims(in.Token, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.signingKey, nil
	})

	if err != nil {
		return &identifierpb.JwtParseResponse{}, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {

		return &identifierpb.JwtParseResponse{Subject: claims.Subject}, nil
	}

	return &identifierpb.JwtParseResponse{}, fmt.Errorf("invalid access token")
}
