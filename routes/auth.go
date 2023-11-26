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
	AuthRouter struct {
		*server.ServerContext
		authService *services.AuthService
	}
)

func NewAuthRouter() *AuthRouter {
	return &AuthRouter{}
}

func (ar *AuthRouter) Register(prefix string, mux *http.ServeMux, ctx *server.ServerContext) {
	ar.ServerContext = ctx
	ar.authService = services.NewAuthService(ctx)
	mux.HandleFunc(prefix+"/user", ar.HandleAuthenticateUser)
	mux.HandleFunc(prefix+"/register", ar.HandleRegisterUser)
}

func (ar *AuthRouter) HandleAuthenticateUser(w http.ResponseWriter, r *http.Request) {
	// 	// u := r.FormValue("username")
	// 	// p := r.FormValue("password")
	// TODO: Use will
	log.Println("Authentication request")
	t := Access{
		Token: ar.authService.NewToken("userId"),
	}

	res, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}

	w.Write([]byte(res))
}

func (ar *AuthRouter) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// e := r.FormValue("email")
	// u := r.FormValue("username")
	// p := r.FormValue("password")game
}

func (ar *AuthRouter) HandleCSRFAuthenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	w.Write([]byte("success"))
}
