package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Semyon981/nexus/gateway/http/auth"
	"github.com/Semyon981/nexus/gateway/http/msg"
	"github.com/Semyon981/nexus/proto/authpb"
	"github.com/Semyon981/nexus/proto/identifierpb"
	"github.com/Semyon981/nexus/proto/msgpb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	httpServer  *http.Server
	authclient  authpb.ServiceClient
	identclient identifierpb.ServiceClient
	msgclient   msgpb.ServiceClient
}

func NewApp() *App {

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c1 := authpb.NewServiceClient(conn)

	conn, err = grpc.Dial("identifier:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c2 := identifierpb.NewServiceClient(conn)

	conn, err = grpc.Dial("msgpb:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c3 := msgpb.NewServiceClient(conn)

	return &App{
		authclient:  c1,
		identclient: c2,
		msgclient:   c3,
	}
}

func (a *App) Run(port string) error {

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		//gin.Logger(),
	)

	auth.RegisterHTTPEndpoints(router, a.authclient)

	authMiddleware := auth.NewAuthMiddleware(a.identclient)

	api := router.Group("/api", authMiddleware)

	api.GET("ping", func(c *gin.Context) { c.JSON(200, "pong") })

	msg.RegisterHTTPEndpoints(api, a.msgclient)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
