package main

import (
	"net/http"

	"github.com/polyglot-k/request_response/http_handler"
	"github.com/polyglot-k/request_response/logger"
)
func main() {
	http.HandleFunc("/", http_handler.Handler)
	logger.Info("Server started on port 8080")
	http.ListenAndServe(":8080", nil)

}