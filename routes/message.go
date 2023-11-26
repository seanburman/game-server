package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/seanburman/game-ws-server/server"
	"github.com/seanburman/game-ws-server/services"
)

type MessageRouter struct {
	server.Router
	*services.MessageService
}

func NewMessageRouter() *MessageRouter {
	return &MessageRouter{}
}

func (mr *MessageRouter) Register(prefix string, mux *http.ServeMux, ctx *server.ServerContext) {
	mr.MessageService = services.NewMessageService(ctx)

	mux.HandleFunc("message", mr.HandleMessages)
}

func (mr *MessageRouter) HandleMessages(w http.ResponseWriter, r *http.Request) {
	ws, err := mr.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not upgrade connection"))
	}

	mr.Conns[ws] = true

	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}

		var msg server.Message
		err = json.Unmarshal(b, &msg)
		if err != nil {
			fmt.Println("unmarshal error:", err)
			continue
		}
		fmt.Println(msg)

		// go func() {
		// broadcast
		// }()
	}
}
