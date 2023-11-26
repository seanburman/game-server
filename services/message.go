package services

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/seanburman/game-ws-server/server"
)

type MessageService struct {
	*server.ServerContext
	websocket.Upgrader
	Conns map[*websocket.Conn]bool
	Rooms map[string][]*websocket.Conn
}

func NewMessageService(ctx *server.ServerContext) *MessageService {
	return &MessageService{
		ServerContext: ctx,
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				fmt.Println(r.Header)
				return true
			}},
		Rooms: make(map[string][]*websocket.Conn),
	}
}

// func (ms *MessageService) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
// 	conn, err := ms.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	m, _ := json.Marshal(server.Message{Name: "Sean"})

// 	conn.WriteMessage(1, []byte(m))

// 	ms. = append(clients, conn)

// 	for {
// 		msgType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			return
// 		}
// 		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

// 		for _, client := range clients {
// 			if err = client.WriteMessage(msgType, msg); err != nil {
// 				return
// 			}
// 		}
// 	}

// }
