package main

import (
	"net/http"

	"github.com/knu-k/logger"
	"github.com/knu-k/server_sent_event/http_handler"
)



func main() {
	http.HandleFunc("/", http_handler.TestHandler)
	http.HandleFunc("/events", http_handler.SseHandler)
	logger.Info("🚀 SSE server started at http://localhost:8080/events")
	http.ListenAndServe(":8080", nil)
}
