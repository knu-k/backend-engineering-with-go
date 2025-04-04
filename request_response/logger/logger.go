package logger

import (
	"log"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

func Info(message string) {
	log.Println(Green + "[INFO] " + message + Reset)
}

func Warn(message string) {
	log.Println(Yellow + "[WARN] " + message + Reset)
}

func Error(message string) {
	log.Println(Red + "[ERROR] " + message + Reset)
}

func Debug(message string) {
	log.Println(Blue + "[DEBUG] " + message + Reset)
}
