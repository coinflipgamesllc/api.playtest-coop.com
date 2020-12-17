package main

import (
	"os"

	"github.com/coinflipgamesllc/api.playtest-coop.com/infrastructure"
	"github.com/coinflipgamesllc/api.playtest-coop.com/ui"
)

// @title Playtest Co-op API
// @version 1.0
// @description This is the backend for all Playtest Co-op related data
// @termsOfService https://playtest-coop.com/terms-of-service

// @contact.name Coin Flip Games
// @contact.email hi@coinflipgames.co

// @host api.playtest-coop.com
// @BasePath /v1
func main() {
	container := &infrastructure.Container{}
	ui.RegisterRoutes(container)

	// Start events handlers
	events := container.EventHandler()
	go func() {
		events.ListenForEvents()
	}()

	router := container.Router()
	router.Run(":" + os.Getenv("PORT"))
}
