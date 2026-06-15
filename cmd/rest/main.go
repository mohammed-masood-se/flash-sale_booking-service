package main

import (
	"log"
	"os"
)

func main() {

	restServerPort := os.Getenv("REST_SERVER_PORT")

	app, err := NewApplication(ApplicationConfig{
		RestServerPort: restServerPort,
	})

	if err != nil {
		log.Printf("failed creating NewApplication: %v", err)
		return
	}

	err = app.Run()

	if err != nil {
		log.Printf("failed running application: %v", err)
		return
	}
}
