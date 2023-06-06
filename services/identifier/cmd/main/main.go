package main

import (
	"log"
	"net"

	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/Semyon981/nexus/services/identifier/config"
	"github.com/Semyon981/nexus/services/identifier/service/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	srv := server.NewServer(viper.GetString("auth.hash_salt"), []byte(viper.GetString("auth.signing_key")))

	s := grpc.NewServer()

	identifierpb.RegisterServiceServer(s, srv)

	lis, err := net.Listen("tcp", ":"+viper.GetString("port"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
