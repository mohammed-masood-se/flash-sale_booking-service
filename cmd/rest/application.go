package main

import (
	"booking-service/internal/adapters/rest"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type ApplicationConfig struct {
	MongoURI       string
	RestServerPort string
}

type Application struct {
	config      *ApplicationConfig
	MongoClient *mongo.Client

	RestServer *rest.RestServer
}

func NewApplication(config ApplicationConfig) (*Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mongoClient, err := GetMongoClient(ctx, config.MongoURI)

	if err != nil {
		return nil, fmt.Errorf("failed getting mongo client: %w", err)
	}
	log.Println("[mongo-client] connected successfully")

	restServer, err := rest.NewRestServer(rest.RestServerConfig{
		Port: config.RestServerPort,
	})

	if err != nil {
		_ = mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed creating NewRestServer: %w", err)
	}

	return &Application{
		config:      &config,
		MongoClient: mongoClient,

		RestServer: restServer,
	}, nil
}

func (app *Application) Run() error {
	log.Printf("[rest-server] started on port %v\n", app.config.RestServerPort)
	err := app.RestServer.Run()

	return fmt.Errorf("failed running RestServer: %w", err)
}

func (app *Application) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.MongoClient.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed disconnecting MongoClient: %w", err)
	}

	return nil
}

func GetMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, fmt.Errorf("failed connecting to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed pinging mongodb client: %w", err)
	}

	return client, nil
}
