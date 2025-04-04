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
	progressMap     = make(map[int32]int) // 진행 중인 작업 저장
	completedTasks  = make(map[int32]int) // 완료된 작업 저장 (100%)
	mutex           = &sync.Mutex{}
)

func processTask(jobId int32) {
	logger.Info("Processing task started: JobId " + strconv.Itoa(int(jobId)))
	for i := 1; i <= 10; i++ {
		time.Sleep(400 * time.Millisecond)
		mutex.Lock()
		progressMap[jobId] = i * 10
		logger.Info("Job Id " + strconv.Itoa(int(jobId)) + " progress: " + strconv.Itoa(progressMap[jobId]) + "%")
		mutex.Unlock()
	}

	mutex.Lock()
	completedTasks[jobId] = 100 // 100% 완료 상태 저장
	delete(progressMap, jobId)  // 메모리 최적화를 위해 제거
	mutex.Unlock()

	logger.Info("JobId " + strconv.Itoa(int(jobId)) + " completed and stored")
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var requestBody RequestBody
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding JSON request body", http.StatusBadRequest)
		return
	}

	jobId := requestBody.JobId
	logger.Info("Received request from: " + r.RemoteAddr)
	logger.Info("JobId: " + strconv.Itoa(int(jobId)))

	// 중복 실행 방지: 완료된 작업이면 새로 실행 X
	mutex.Lock()
	_, exists := completedTasks[jobId]
	mutex.Unlock()
	if exists {
		logger.Info("JobId " + strconv.Itoa(int(jobId)) + " already completed, skipping")
		response := map[string]interface{}{
			"message":  "Task already completed",
			"progress": 100,
		}
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		return
	}

	go processTask(jobId)

	response := map[string]string{"message": "Task created successfully"}
	jsonResponse, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

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

	mutex.Lock()
	progress, inProgress := progressMap[int32(jobId)]
	completedProgress, completed := completedTasks[int32(jobId)]
	mutex.Unlock()
	
	// 진행 중이면 진행률 반환
	if inProgress {
		response := map[string]int{"job_id": int(jobId), "progress": progress}
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		logger.Info("CheckStatusHandler: JobId " + jobIdStr + " progress: " + strconv.Itoa(progress) + "%")
		return
	}

	// 완료된 작업이면 완료된 진행률 반환 (100%)
	if completed {
		response := map[string]int{"job_id": int(jobId), "progress": completedProgress}
		jsonResponse, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		logger.Info("CheckStatusHandler: JobId " + jobIdStr + " already completed, returning 100%")
		return
	}

	// 존재하지 않는 작업이면 에러 반환
	http.Error(w, "Job not found", http.StatusNotFound)
	logger.Info("CheckStatusHandler: JobId " + jobIdStr + " not found")
}
