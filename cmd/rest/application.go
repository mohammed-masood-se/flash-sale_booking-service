package main

import (
	"booking-service/internal/adapters/mongorepo"
	"booking-service/internal/adapters/rediscache"
	"booking-service/internal/adapters/rest"
	"booking-service/internal/core/services"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type ApplicationConfig struct {
	MongoURI       string
	RedisAddr      string
	RestServerPort string
}

type Application struct {
	config      *ApplicationConfig
	MongoClient *mongo.Client
	RedisClient *redis.Client

	RestServer *rest.RestServer
}

func NewApplication(config ApplicationConfig) (*Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	mongoClient, err := GetMongoClient(ctx, config.MongoURI)

	if err != nil {
		return nil, fmt.Errorf("failed getting mongo client: %w", err)
	}
	log.Println("[mongo-client] connected successfully")

	database := mongoClient.Database("flash-sale")
	txmanager := mongorepo.NewMongoTransactionManager(mongoClient)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed pinging redis client: %w", err)
	}
	log.Println("[redis-client] connected successfully")

	usersCollection := database.Collection("users")
	registrationCollection := database.Collection("registration")

	userRepository, err := mongorepo.NewUserRepository(ctx, usersCollection, registrationCollection)
	if err != nil {
		return nil, fmt.Errorf("failed creating NewUserRepository: %w", err)
	}
	userCache := rediscache.NewUserCache(redisClient)
	userService := services.NewUserService(txmanager, userRepository, userCache)

	restServer, err := rest.NewRestServer(rest.RestServerConfig{
		Port: config.RestServerPort,
	}, &rest.Services{
		UserService: userService,
	})

	if err != nil {
		_ = mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed creating NewRestServer: %w", err)
	}

	return &Application{
		config:      &config,
		MongoClient: mongoClient,
		RedisClient: redisClient,

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

	if err := app.RedisClient.Close(); err != nil {
		return fmt.Errorf("failed closing redis client connection: %w", err)
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
