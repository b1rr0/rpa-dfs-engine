package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File
var logger *log.Logger

func initLogger() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exeDir := filepath.Dir(exePath)
	logDir := exeDir

	date := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("facebook-login_%s.log", date)
	logFilePath := filepath.Join(logDir, logFileName)

	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger = log.New(logFile, "", log.LstdFlags)

	logger.Println("=== SESSION START ===")
	logger.Printf("Start time: %s", time.Now().Format("2006-01-02 15:04:05"))
	logger.Printf("Arguments: %v", os.Args)

	workDir, _ := os.Getwd()
	logger.Printf("Working directory: %s", workDir)

	return nil
}

func closeLogger() {
	if logger != nil {
		logger.Println("=== SESSION END ===")
	}
	if logFile != nil {
		logFile.Close()
	}
}

func LogInfo(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf("ℹ️ %s\n", message)
	if logger != nil {
		logger.Printf("INFO: %s", message)
	}
}

func LogSuccess(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf("✅ %s\n", message)
	if logger != nil {
		logger.Printf("SUCCESS: %s", message)
	}
}

func LogError(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf("❌ %s\n", message)
	if logger != nil {
		logger.Printf("ERROR: %s", message)
	}
}

func LogWarning(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf("⚠️ %s\n", message)
	if logger != nil {
		logger.Printf("WARNING: %s", message)
	}
}

func LogDebug(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	fmt.Printf("🐛 %s\n", message)
	if logger != nil {
		logger.Printf("DEBUG: %s", message)
	}
}
