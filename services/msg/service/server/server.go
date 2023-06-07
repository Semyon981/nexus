package server

import (
	"context"
	"time"

	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/msg/service"
	"github.com/Semyon981/nexus/services/msg/service/repository/postgresql"
	"github.com/jmoiron/sqlx"
)

type server struct {
	msgpb.UnimplementedServiceServer
	repo        service.Repository
	usersclient userspb.UserServiceClient
}

func NewServer(db *sqlx.DB, usersclient userspb.UserServiceClient) *server {
	return &server{
		repo:        *postgresql.NewUserRepository(db),
		usersclient: usersclient,
	}
}

func (s *server) GetMessages(ctx context.Context, in *msgpb.GetMessagesRequest) (*msgpb.GetMessagesResponse, error) {
	res, err := s.repo.GetMessages(ctx, in.IdFrom, in.IdTo, in.Limit, in.Offset)
	if err != nil {
		return &msgpb.GetMessagesResponse{}, err
	}

	response := make([]*msgpb.Message, 0)
	for i := range res {
		message := msgpb.Message{
			IdMessages: res[i].Id_Messages,
			IdFrom:     res[i].Id_from,
			IdTo:       res[i].Id_to,
			Msg:        res[i].Msg,
			Time:       res[i].Time.Unix(),
		}
		response = append(response, &message)
	}

	return &msgpb.GetMessagesResponse{Messages: response}, nil

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
