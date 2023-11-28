package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/csrf"

	"github.com/seanburman/game-ws-server/server"
	"github.com/seanburman/game-ws-server/services"
)

type (
	Access struct {
		Token     string `json:"token,omitempty"`
		SessionID string `json:"session_id,omitempty"`
	}
)

type authRoute struct {
	server.Route
	*services.AuthService
}

func NewAuthRoute() *authRoute {
	ar := &authRoute{
		Route: *server.NewRoute("/auth"),
	}
	ar.AuthService = services.NewAuthService()
	ar.SubRoutes()
	// ar.Use(middleware.MiddlewareVerifyJWT)
	return ar
}

func (ar *authRoute) SubRoutes() {
	ar.Handle("/user", ar.HandleAuthenticateUser)
}

// mux.HandleFunc(prefix+"/user", ar.HandleAuthenticateUser)
// mux.HandleFunc(prefix+"/register", ar.HandleRegisterUser)

func (ar *authRoute) HandleAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// 	// u := r.FormValue("username")
	// 	// p := r.FormValue("password")
	t := Access{
		Token: ar.AuthService.NewToken("userId"),
	}

	res, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}

	w.Write([]byte(res))
}

func (ar *authRoute) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// e := r.FormValue("email")
	// u := r.FormValue("username")
	// p := r.FormValue("password")game
}

func (ar *authRoute) HandleCSRFAuthenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Write([]byte("success"))
}
