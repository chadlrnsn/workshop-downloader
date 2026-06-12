package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

func InitLogger() error {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	appDir := filepath.Join(home, ".workshop-downloader")
	_ = os.MkdirAll(appDir, 0755)

	logPath := filepath.Join(appDir, "debug_logs.txt")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	logFile = file

	WriteLog("=== Application started / Logger initialized ===")
	return nil
}

func WriteLog(message string) {
	if logFile == nil {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] %s\n", timestamp, message)
	_, _ = logFile.WriteString(line)
}

func WriteError(err error, context string) {
	if err == nil {
		return
	}
	WriteLog(fmt.Sprintf("ERROR in %s: %s", context, err.Error()))
}

func Close() {
	if logFile != nil {
		WriteLog("=== Application stopping ===")
		_ = logFile.Close()
	}
}
