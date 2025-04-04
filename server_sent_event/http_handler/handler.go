package http_handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/knu-k/logger"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the SSE server!")
}

func SseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("🔗 Connected:"+ r.RemoteAddr)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			logger.Error("❌ Disconnected:"+ r.RemoteAddr)
			return
		default:
			// 실제 메시지 전송
			fmt.Fprintf(w, "data: Message %d at %s\n\n", i, time.Now().Format(time.RFC3339))
			// ping 전송 (SSE 주석용)
			fmt.Fprintf(w, ": heartbeat\n\n")

			flusher.Flush()
			time.Sleep(2 * time.Second)
		}
	}
}
