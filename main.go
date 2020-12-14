package main

import "github.com/coinflipgamesllc/api.playtest-coop.com/app"

// @title Playtest Co-op API
// @version 1.0
// @description This is the backend for all Playtest Co-op related data
// @termsOfService https://playtest-coop.com/terms-of-service

// @contact.name Coin Flip Games
// @contact.email hi@coinflipgames.co

// @host api.playtest-coop.com
// @BasePath /v1
func main() {
	server := app.NewServer()

	server.Run()
}
