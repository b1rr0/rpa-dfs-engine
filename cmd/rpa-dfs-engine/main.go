package main

import (
	"fmt"
	"os"

	"rpa-dfs-engine/internal/handlers"
	"rpa-dfs-engine/internal/logger"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("❌ Logger init error: %v\n", err)
	}
	defer logger.CloseLogger()

	appStart()
}

func appStart() {
	logger.LogInfo("RPA DFS Engine started")
	handler := handlers.GetHandler()

	logger.LogInfo("Executing handler: %s", handler.GetDescription())

	if err := handler.Execute(); err != nil {
		logger.LogError("Handler execution failed: %v", err)
		os.Exit(1)
	}

	logger.LogInfo("Application finished successfully")
}
