package main

import (
	"log"
	"net"

	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/Semyon981/nexus/services/users/service/server"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

func main() {
	db := InitDB()
	srv := server.NewServer(db)

	s := grpc.NewServer()

	userspb.RegisterUserServiceServer(s, srv)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func InitDB() *sqlx.DB {
	url := "postgres://postgres:password@dbusers/postgres"
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
