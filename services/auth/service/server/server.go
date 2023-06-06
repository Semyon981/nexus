package server

import (
	"context"
	"log"
	"strconv"

	"github.com/Semyon981/nexus/proto/authpb"
	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/Semyon981/nexus/proto/userspb"
)

type server struct {
	authpb.UnimplementedServiceServer
	identclient identifierpb.ServiceClient
	usrclient   userspb.UserServiceClient
	EXPJWT      int64
}

func NewServer(identclient identifierpb.ServiceClient, usrclient userspb.UserServiceClient, EXPJWT int64) *server {
	return &server{identclient: identclient, usrclient: usrclient, EXPJWT: EXPJWT}
}

func (s *server) SignIn(ctx context.Context, in *authpb.SignInRequest) (*authpb.SignInResponse, error) {
	passhash, err := s.identclient.Hash(ctx, &identifierpb.HashRequest{Str: in.Password})

	if err != nil {
		log.Println("1:", err)
		return &authpb.SignInResponse{}, err
	}

	resp, err := s.usrclient.GetUserAuth(ctx, &userspb.GetUserAuthRequest{Number: in.Number, Password: passhash.Hash})

	if err != nil {
		log.Println("2:", err)
		return &authpb.SignInResponse{}, err
	}

	resJwt, err := s.identclient.JwtGen(ctx, &identifierpb.JwtGenRequest{Subject: strconv.FormatInt(resp.IdUsers, 10), Time: s.EXPJWT})
	return &authpb.SignInResponse{Token: resJwt.Token}, err
}

func (s *server) SignUp(ctx context.Context, in *authpb.SignUpRequest) (*authpb.SignUpResponse, error) {

	passhash, err := s.identclient.Hash(ctx, &identifierpb.HashRequest{Str: in.Password})

	if err != nil {
		log.Println("1:", err)
		return &authpb.SignUpResponse{}, err
	}

	_, err = s.usrclient.CreateUser(ctx, &userspb.CreateUserRequest{Number: in.Number, Password: passhash.Hash,
		Name: in.Name, Lastname: in.Lastname})

	return &authpb.SignUpResponse{}, err

}
