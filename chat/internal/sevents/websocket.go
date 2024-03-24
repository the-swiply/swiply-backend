package sevents

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"github.com/the-swiply/swiply-backend/pkg/houston/auf"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WS struct {
	chatService *service.ChatService
}

func NewWS(chatService *service.ChatService) *WS {
	return &WS{
		chatService: chatService,
	}
}

func (w *WS) Connect(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		errMsg := fmt.Sprintf("can't upgrade connection to ws: %v", err)
		fmt.Fprint(rw, errMsg)
		return
	}

	userID := auf.ExtractUserIDFromContext[uuid.UUID](r.Context())
	go w.handleClientConn(conn, userID)
}

func (w *WS) handleClientConn(conn *websocket.Conn, userID uuid.UUID) {
	defer func() {
		w.chatService.RemoveChatClient(userID)
		conn.Close()
	}()

	w.chatService.AddChatClient(userID, &Client{
		conn: conn,
	})

	for {
		mt, _, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
	}
}
