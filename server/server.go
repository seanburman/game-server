package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/seanburman/game-ws-server/config"
	db "github.com/seanburman/game-ws-server/db/clients"
)

var Server = NewServer("")

type server struct {
	mux  *http.ServeMux
	Port string
	*Router
	*ServerContext
	*db.Clients
	Conns map[*websocket.Conn]bool
	Rooms map[string][]*websocket.Conn
}

func NewServer(basepath string) *server {
	s := &server{
		mux:           http.NewServeMux(),
		Port:          config.Env().PORT,
		ServerContext: NewServerContext(),
		Clients:       db.NewClients(),
		Conns:         make(map[*websocket.Conn]bool),
	}
	s.Router = newRouter(basepath, s)
	return s
}

func (s *server) ListenAndServe() error {
	if s.Router == nil {
		return errors.New("missing router initializing in root")
	}
	s.ApplyMiddleware()

	for _, r := range s.Routes {
		for _, sr := range r.subroutes {
			fmt.Println(sr)
			s.mux.Handle(sr.prefix, *sr.handler)
		}
	}
	s.printInfo()
	CORS := cors.Default().Handler(s.mux)
	return http.ListenAndServe(s.Port, CORS)
}

func (s *server) UseRouter(prefix string) {
	s.Router = newRouter(prefix, s)
}

func (s *server) ServeHealthCheck(prefix string) {
	s.mux.HandleFunc(prefix, func(w http.ResponseWriter, r *http.Request) {
		m, err := json.Marshal(Message{Message: "server is healthy"})
		if err != nil {
			log.Panic("error marshalling response message")
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(m)
	})
}

func (s *server) ServeStaticFiles(prefix string, dir string) {
	s.mux.Handle(prefix, http.FileServer(http.Dir(dir)))
}

func (s *server) printInfo() {
	var Reset = "\033[0m"
	var Blue = "\033[34m"
	// var Orange = "\033[48:5:208m%s\033[m\n"
	var Red = "\033[31m"
	var Yellow = "\033[33m"
	var White = "\033[97m"
	logo1 := `   ______                          _____`
	logo2 := `  / ____/___  _________   ___     / ___/__  ______   _____  _____`
	logo3 := ` / / __/ __ \|  __  __ \/  _ \    \__\/ _ \/ ___/ | / / _ \/ ___/`
	logo4 := `/ /_/ / /_/  / / / / / /  __/   ___/ /  __/ /   | |/ /  __/ /`
	logo5 := `\____/\__/,_/_/ /_/ /_/\___/   /____/\___/_/    |___/\___/_/`
	fmt.Println(White + logo1 + Reset)
	fmt.Println(White + logo2 + Reset)
	fmt.Println(Yellow + logo3 + Reset)
	fmt.Println(Red + logo4 + Reset)
	fmt.Println(Blue + logo5 + Reset)

	fmt.Printf(White+"\nListening on port %s..."+Reset, s.Port)
}
