package server

import (
	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/Semyon981/nexus/services/msg/service"
	"github.com/Semyon981/nexus/services/msg/service/repository/postgresql"
	"github.com/jmoiron/sqlx"
)

type server struct {
	msgpb.UnimplementedUserServiceServer
	repo service.Repository
}

func NewServer(db *sqlx.DB) *server {
	return &server{
		repo: *postgresql.NewUserRepository(db),
	}
}
