package main

import (
	"log"
	"net"

	"github.com/Semyon981/nexus/proto/authpb"
	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/auth/config"
	"github.com/Semyon981/nexus/services/auth/service/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c1 := identifierpb.NewServiceClient(conn)
	c2 := userspb.NewUserServiceClient(conn)

	srv := server.NewServer(c1, c2, viper.GetInt64("EXPJWT"))

	s := grpc.NewServer()

	authpb.RegisterServiceServer(s, srv)

	lis, err := net.Listen("tcp", ":"+viper.GetString("port"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
