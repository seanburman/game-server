package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/seanburman/game-ws-server/server"
)

type SessionID uuid.UUID

type Session struct {
	SessionID SessionID
	User      *websocket.Conn
	Client    *websocket.Conn
}

type SessionService struct {
	*server.ServerContext
	Conns    map[*websocket.Conn]bool
	Sessions map[SessionID]*Session
}

func NewSessionService() *SessionService {
	return &SessionService{
		Conns:    make(map[*websocket.Conn]bool),
		Sessions: make(map[SessionID]*Session),
	}
}

func (gs *SessionService) CreateSession(sid SessionID) {
	s := &Session{SessionID: sid}
	gs.Sessions[sid] = s
	log.Println("new session: ", s.SessionID)
}

func (gs *SessionService) JoinUserSession(conn *websocket.Conn, sid SessionID) error {
	if gs.Sessions[sid] != nil {
		gs.Sessions[sid].User = conn
		log.Println("user joined session: ", gs.Sessions[sid])
		return nil
	}
	return errors.New("bad request")
}

func (gs *SessionService) JoinClientSession(conn *websocket.Conn, sid SessionID) error {
	if gs.Sessions[sid] != nil {
		gs.Sessions[sid].Client = conn
		log.Println("client joined session: ", gs.Sessions[sid])
		return nil
	}
	return errors.New("bad session")
}

func (gs SessionService) NewSessionID() (SessionID, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return SessionID(uuid), nil
	}
	log.Println(err)
	return SessionID(uuid), err
}

func (gs *SessionService) SessionIDFromParam(r *http.Request) (SessionID, error) {
	param := r.URL.Query().Get("session")
	var err error
	id, err := uuid.Parse(param)
	if err != nil {
		return SessionID{}, err
	}
	id = uuid.Must(id, err)
	if err != nil {
		return SessionID{}, err
	}
	return SessionID(id), nil
}
