package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/seanburman/game-ws-server/server"
	"github.com/seanburman/game-ws-server/services"
)

type GameRouter struct {
	server.Router
	*services.SessionService
}

func NewGameRouter() *GameRouter {
	return &GameRouter{}
}

func (gr *GameRouter) Register(prefix string, mux *http.ServeMux, ctx *server.ServerContext) {
	gr.SessionService = services.NewSessionService(ctx)
	mux.HandleFunc(prefix+"/session/create", gr.handleCreateSession)
	mux.HandleFunc(prefix+"/new", gr.handleServeGame)
	mux.HandleFunc(prefix+"/ws", gr.handleHandShake)
}

func (gr *GameRouter) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	// TODO: JWT MUST BE VALIDATE BEFORE THIS HANDLER
	sid, err := gr.NewSessionID()
	if err != nil {
		http.Error(w, "error creating session", http.StatusInternalServerError)
	}

	// var t Access
	// err = json.NewDecoder(r.Body).Decode(&t)
	// if err != nil || t.Token == "" {
	// 	http.Error(w, "missing token", http.StatusBadRequest)
	// }
	// fmt.Println(t.Token)
	gr.CreateSession(sid)
	u := uuid.UUID(sid)
	s, err := json.Marshal(Access{SessionID: fmt.Sprint(u)})
	if err != nil {
		http.Error(w, "error marshalling session id", http.StatusInternalServerError)
	}

	w.Write([]byte(s))
}

func (gr *GameRouter) handleServeGame(w http.ResponseWriter, r *http.Request) {
	sid, err := gr.SessionIDFromParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid session id"))
		return
	}
	log.Println("user being served session: ", services.SessionID(sid))
	// Check sessions
	err = gr.JoinUserSession(nil, services.SessionID(sid))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}
	entry := []byte(fmt.Sprintf(
		`
		<!DOCTYPE html>
		<script src="../wasm_exec.js"></script>
		<script>
			function func(){
				return "%s"
			}
		</script>
		<script>
		// Polyfill
		if (!WebAssembly.instantiateStreaming) {
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer();
				return await WebAssembly.instantiate(source, importObject);
			};
		}

		const go = new Go();
		WebAssembly.instantiateStreaming(fetch("../game.wasm"), go.importObject).then(result => {
			go.run(result.instance);
		});
		</script>
		`, sid))

	w.Write([]byte(entry))
}

func (gr GameRouter) handleHandShake(w http.ResponseWriter, r *http.Request) {

}
