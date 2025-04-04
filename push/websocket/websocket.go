package websocket

import (
	"sync"

	"github.com/polyglot-k/push/logger"
	"golang.org/x/net/websocket"
)

type Connection struct {
	WebSocket *websocket.Conn
	Sequence  int
}

var (
	connections []Connection
	mu          sync.Mutex
)

func HandleWebSocket(ws *websocket.Conn) {
	defer ws.Close()

	clientAddr := ws.Request().RemoteAddr
	logger.Info("클라이언트 연결됨: " + clientAddr)

	appendConnection(ws)

	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			logger.Error("연결 종료: " + err.Error())
			removeConnection(ws)
			break
		}

		logger.Info("수신된 메시지: " + msg)
		sendMessageBroadcast(msg)
	}
}

func sendMessageBroadcast(msg string) {
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range connections {
		err := websocket.Message.Send(conn.WebSocket, "서버 응답: "+msg)
		if err != nil {
			logger.Error("메시지 전송 오류: " + err.Error())
		}
	}
}

func appendConnection(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	connections = append(connections, Connection{WebSocket: ws, Sequence: len(connections) + 1})
}

func removeConnection(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for i, conn := range connections {
		if conn.WebSocket == ws {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}
