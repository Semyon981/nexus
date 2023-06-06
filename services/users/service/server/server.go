package server

import (
	"context"

	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/users/service"
	"github.com/Semyon981/nexus/services/users/service/repository/postgresql"

	"github.com/jmoiron/sqlx"
)

type server struct {
	userspb.UnimplementedUserServiceServer
	repo service.UserRepository
}

func NewServer(db *sqlx.DB) *server {
	return &server{
		repo: *postgresql.NewUserRepository(db),
	}
}

func (s *server) GetUserAuth(ctx context.Context, in *userspb.AuthUserRequest) (*userspb.AuthUserResponse, error) {
	id_users, err := s.repo.GetUserAuth(ctx, in.Number, in.Password)
	return &userspb.AuthUserResponse{IdUsers: id_users}, err

}

func (s *server) CreateUser(ctx context.Context, in *userspb.CreateUserRequest) (*userspb.CreateUserResponse, error) {

	err := s.repo.CreateUser(ctx, in.Number, in.Password, in.Name, in.Lastname)
	return &userspb.CreateUserResponse{}, err
}

func (s *server) GetUserId(ctx context.Context, in *userspb.GetUserIdRequest) (*userspb.GetUserIdResponse, error) {

	resp, err := s.repo.GetUserId(ctx, in.IdUsers)

	return &userspb.CreateUserResponse{IdUsers: resp.Id_users, Number: resp.Number, Password: resp.Password, Name: resp.Name, Lastname: resp.Lastname}, err
}
