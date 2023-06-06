package server

import (
	"context"
	"time"

	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/msg/models"
	"github.com/Semyon981/nexus/services/msg/service"
	"github.com/Semyon981/nexus/services/msg/service/repository/postgresql"
	"github.com/jmoiron/sqlx"
)

type server struct {
	msgpb.UnimplementedServiceServer
	repo        service.Repository
	usersclient userspb.UserServiceClient
}

func NewServer(db *sqlx.DB) *server {
	return &server{
		repo: *postgresql.NewUserRepository(db),
	}
}

func (s *server) GetMessages(ctx context.Context, in *msgpb.GetMessagesRequest) (*msgpb.GetMessagesResponse, error) {
	res := make([]models.Message, 0)
	s.repo.GetMessages(ctx)

}

func (s *server) SendMessage(ctx context.Context, in *msgpb.SendMessageRequest) (*msgpb.SendMessageResponse, error) {
	_, err := s.usersclient.GetUserId(ctx, &userspb.GetUserIdRequest{IdUsers: in.IdTo})

	if err != nil {
		return &msgpb.SendMessageResponse{}, err
	}

	err = s.repo.SendMessage(ctx, in.IdFrom, in.IdTo, in.Msg, time.Unix(in.Time, 0))
	if err != nil {
		return &msgpb.SendMessageResponse{}, err
	}

	return &msgpb.SendMessageResponse{}, nil

}
