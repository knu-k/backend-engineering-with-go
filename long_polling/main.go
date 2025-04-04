package main

import (
	"net/http"

	"github.com/knu-k/logger"
	"github.com/knu-k/long_polling/http_handler"
)



func main() {
	http.HandleFunc("/check_status", http_handler.CheckStatusHandler)
	http.HandleFunc("/task", http_handler.CreateTaskHandler)
	logger.Info("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

/**
curl -X POST http://localhost:8080/task -H "Content-Type: application/json" -d "{\"job_id\": 1234}"

curl -X GET http://localhost:8080/check_status?job_id=1234 
**/