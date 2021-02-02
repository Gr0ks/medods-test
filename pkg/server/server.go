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

func (a *App) Run(port string) {
	r := mux.NewRouter()
	delivery.RegisterHTTPEndpoints(r, a.useCase)
	go func() {
		log.Fatal(http.ListenAndServe(":" + port, r))
	}()
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://auth-mongo:27017"))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("authorizer")
}