package main

import (
	"log"

	"github.com/seanburman/game-ws-server/routes"
	"github.com/seanburman/game-ws-server/server"
)

func main() {
	s := server.Server
	s.UseRouter("")
	s.ServeHealthCheck("/health-check")
	s.ServeStaticFiles("/", "./static")
	// s.Router.UseSequence(middleware.MiddlewareVerifyJWT)

	s.Handle("/auth", routes.NewAuthRoute().Route)
	s.Handle("/game", routes.NewGameRoute().Route)

	log.Fatal(s.ListenAndServe())
}
