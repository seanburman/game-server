package main

import (
	"log"

	"github.com/seanburman/game-ws-server/routes"
	"github.com/seanburman/game-ws-server/server"
)

func main() {
	s := server.NewServer()

	s.UseHealthCheck("/health-check")
	s.UseRouter("/auth", routes.NewAuthRouter())
	s.UseRouter("/game", routes.NewGameRouter())

	log.Fatal(s.ListenAndServe())
}
