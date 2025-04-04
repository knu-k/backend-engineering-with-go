package main

import (
	"net/http"

	"github.com/polyglot-k/push/logger"
	ws "github.com/polyglot-k/push/websocket"
	"golang.org/x/net/websocket"
)

func main() {
	http.Handle("/ws", websocket.Handler(ws.HandleWebSocket))

	logger.Info("WebSocket Server Running on port :8080") 
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Info("서버 오류: " + err.Error()) 
	}
}
