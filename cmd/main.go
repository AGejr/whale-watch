package main

import (
	"log"

	"github.com/agejr/go-beans/internal/docker"
	"github.com/agejr/go-beans/internal/ui"
)

func main() {
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	app := ui.NewApp(dockerClient)

	if err := app.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
