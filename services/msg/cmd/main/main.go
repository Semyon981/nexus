package main

import (
	"log"
	"net"

	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/msg/service/server"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := userspb.NewUserServiceClient(conn)

	db := InitDB()
	srv := server.NewServer(db, c)

	s := grpc.NewServer()
	msgpb.RegisterServiceServer(s, srv)

	lis, err := net.Listen("tcp", ":50054")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func InitDB() *sqlx.DB {
	url := "postgres://postgres:password@localhost"
	database, err := sqlx.Open("pgx", url)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return database
}
