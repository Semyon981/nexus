package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Semyon981/nexus/gateway/auth"
	"github.com/Semyon981/nexus/proto/userspb"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authhttp "github.com/Semyon981/nexus/gateway/auth/http"
	authusecase "github.com/Semyon981/nexus/gateway/auth/usecase"
)

type App struct {
	httpServer *http.Server
	authUC     auth.UseCase
}

func NewApp() *App {

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := userspb.NewUserServiceClient(conn)

	return &App{
		authUC: authusecase.NewAuthUseCase(
			c,
			"salt",
			[]byte("key"),
			time.Minute*60,
		),
	}
}

func (a *App) Run(port string) error {

	router := gin.Default()
	router.Use(
		gin.Recovery(),
		//gin.Logger(),
	)

	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	//authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	//api := router.Group("/api", authMiddleware)

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
