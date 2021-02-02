package server

import (
	"medods-test/pkg/auth"
	"medods-test/pkg/auth/delivery"
	"medods-test/pkg/auth/usecase"
	mongoRepo "medods-test/pkg/auth/repository/mongo"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"log"
	"context"
	"time"
	"github.com/gorilla/mux"
	"os"
	"os/signal"
)

type App struct {
	httpServer *http.Server

	useCase auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := mongoRepo.NewSessionRepository(db, "sessions")
	useCase := usecase.NewAccessPairCreator(
		userRepo,
		"auth.hash_salt",
		[]byte("auth.signing_key"),
		15*60*time.Second,
	)

	return &App{
		useCase: useCase,
	}
}

func (a *App) Run(port string) error {
	r := mux.NewRouter()
	delivery.RegisterHTTPEndpoints(r, a.useCase)
	go func() {
		log.Fatal(http.ListenAndServe(":8090", r))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatal(err)
    }
	
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("authorizer")
}