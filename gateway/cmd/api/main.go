package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Semyon981/nexus/proto/userspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := userspb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AuthUser(ctx, &userspb.AuthUserRequest{Number: "+71253264367", Password: "qwerty123"})
	fmt.Println(r, err)
	req, err := c.CreateUser(ctx, &userspb.CreateUserRequest{Number: "+71253264367", Password: "qwerty123", Name: "Sema", Lastname: "pisa"})
	fmt.Println(req, err)

}
