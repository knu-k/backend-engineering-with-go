package http_handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/knu-k/logger"
)

type RequestBody struct {
	JobId int32 `json:"job_id"`
}

var (
	progressMap    = make(map[int32]int)
	completedTasks = make(map[int32]int)
	mutex          = &sync.Mutex{}
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func processTask(jobId int32) {
	logger.Info("Started processing JobId: " + strconv.Itoa(int(jobId)))
	for i := 1; i <= 10; i++ {
		time.Sleep(400 * time.Millisecond)

		mutex.Lock()
		progressMap[jobId] = i * 10
		mutex.Unlock()

		logger.Info("JobId " + strconv.Itoa(int(jobId)) + " progress: " + strconv.Itoa(i*10) + "%")
	}

	mutex.Lock()
	completedTasks[jobId] = 100
	delete(progressMap, jobId)
	mutex.Unlock()

	logger.Info("JobId " + strconv.Itoa(int(jobId)) + " completed")
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	jobId := requestBody.JobId
	logger.Info("Received request from " + r.RemoteAddr + " for JobId: " + strconv.Itoa(int(jobId)))

	mutex.Lock()
	_, exists := completedTasks[jobId]
	mutex.Unlock()

	if exists {
		logger.Info("JobId " + strconv.Itoa(int(jobId)) + " already completed")
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"message":  "Task already completed",
			"progress": 100,
		})
		return
	}

	go processTask(jobId)

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Task created successfully",
	})
}

func isJobCompleted(jobId int32) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, completed := completedTasks[jobId]
	return completed
}

func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	jobIdStr := r.URL.Query().Get("job_id")
	if jobIdStr == "" {
		http.Error(w, "Missing job_id parameter", http.StatusBadRequest)
		return
	}

	jobId, err := strconv.Atoi(jobIdStr)
	if err != nil {
		http.Error(w, "Invalid job_id format", http.StatusBadRequest)
		return
	}

	if !isJobCompleted(int32(jobId)) {
		// 작업이 존재하지 않는 경우도 확인
		mutex.Lock()
		_, exists := progressMap[int32(jobId)]
		mutex.Unlock()

		if !exists {
			http.Error(w, "Job not found", http.StatusNotFound)
			logger.Info("CheckStatusHandler: JobId " + jobIdStr + " not found")
			return
		}

		for !isJobCompleted(int32(jobId)) {
			time.Sleep(100 * time.Millisecond)
		}
	}

	logger.Info("CheckStatusHandler: JobId " + jobIdStr + " completed")
	respondJSON(w, http.StatusOK, map[string]int{
		"job_id":   jobId,
		"progress": 100,
	})
}
