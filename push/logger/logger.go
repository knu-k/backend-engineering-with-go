package logger

import (
	"log"
)

// ANSI 색상 코드
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

// Info 로그 (초록색)
func Info(message string) {
	log.Println(Green + "[INFO] " + message + Reset)
}

// Warn 로그 (노란색)
func Warn(message string) {
	log.Println(Yellow + "[WARN] " + message + Reset)
}

// Error 로그 (빨간색)
func Error(message string) {
	log.Println(Red + "[ERROR] " + message + Reset)
}

// Debug 로그 (파란색)
func Debug(message string) {
	log.Println(Blue + "[DEBUG] " + message + Reset)
}
