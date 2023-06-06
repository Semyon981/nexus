package server

import (
	"context"
	"users/proto/userspb"
	"users/service/repository/postgresql"

	"github.com/jmoiron/sqlx"
)

type server struct {
	userspb.UnimplementedUserServiceServer
	repo postgresql.UserRepository
}

func NewServer(db *sqlx.DB) *server {
	return &server{
		repo: *postgresql.NewUserRepository(db),
	}
}

func (s *server) AuthUser(ctx context.Context, in *userspb.AuthUserRequest) (*userspb.AuthUserResponse, error) {
	id_users, err := s.repo.GetUser(ctx, in.Number, in.Password)
	return &userspb.AuthUserResponse{IdUsers: id_users}, err

}

func (s *server) CreateUser(ctx context.Context, in *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	err := s.repo.CreateUser(ctx, in.Number, in.Password, in.Name, in.Lastname)
	return &userspb.CreateUserResponse{}, err
}
