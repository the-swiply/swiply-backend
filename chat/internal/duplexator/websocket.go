package duplexator

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"net/http"
	"sync"
)

const (
	authorizationHeader = "Authorization"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type messageHandler interface {
}

type WS struct {
	clients   map[uuid.UUID]*websocket.Conn
	clientsMu sync.Mutex
}

func NewWS() *WS {
	return &WS{
		clients:   make(map[uuid.UUID]*websocket.Conn),
		clientsMu: sync.Mutex{},
	}
}

func (w *WS) Connect(rw http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get(authorizationHeader)
	if auth == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(rw, "authorization header is not set")
		return
	}

	// TODO: parse id from auth jwt.
	t := ""
	id, err := uuid.Parse(t)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "can't parse id from token")
		return
	}

	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		errMsg := fmt.Sprintf("can't upgrade connection to ws: %v", err)
		loggy.Warnln(errMsg)
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, errMsg)
		return
	}

	w.clientsMu.Lock()
	w.clients[id] = conn
	w.clientsMu.Unlock()

	go w.readMessages(id, conn)
}

func (w *WS) readMessages(id uuid.UUID, conn *websocket.Conn) {
	defer func() {
		w.clientsMu.Lock()
		delete(w.clients, id)
		w.clientsMu.Unlock()
	}()

	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
	}
}

func (w *WS) SendMessageToClients(ctx context.Context, ids ...uuid.UUID) error {

}
