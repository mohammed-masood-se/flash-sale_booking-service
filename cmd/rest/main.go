package main

import (
	"log"
	"os"
)

func main() {

	restServerPort := os.Getenv("REST_SERVER_PORT")
	mongoUri := os.Getenv("MONGODB_CONNECTION_STRING")
	redisAddr := os.Getenv("REDIS_ADDRESS")

	app, err := NewApplication(ApplicationConfig{
		RestServerPort: restServerPort,
		RedisAddr:      redisAddr,
		MongoURI:       mongoUri,
	})

	if err != nil {
		log.Printf("failed creating NewApplication: %v", err)
		return
	}
	defer app.Shutdown()

	err = app.Run()

	if err != nil {
		log.Printf("failed running application: %v", err)
		return
	}
}
