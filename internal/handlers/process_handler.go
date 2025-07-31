package handlers

import (
	"rpa-dfs-engine/internal/logger"
)

// ProcessHandler handles processing requests that include email and token credentials.
// Currently this is a placeholder implementation that logs processing information.
type ProcessHandler struct{}

// NewProcessHandler creates a new ProcessHandler instance.
// It returns the Handler interface to promote loose coupling.
func NewProcessHandler() Handler {
	return &ProcessHandler{}
}

// Execute implements the Handler interface for processing with credentials.
// This is currently a placeholder implementation that logs the processing mode.
func (h *ProcessHandler) Execute() error {
	logger.LogInfo("=== RPA DFS Engine - Process Mode ===")
	logger.LogInfo("Processing handler initialized with email and token")
	logger.LogInfo("Note: This is currently a placeholder implementation")
	logger.LogSuccess("Process handler completed successfully")

	return nil
}

// GetDescription implements the Handler interface and returns a description
// of what this handler does.
func (h *ProcessHandler) GetDescription() string {
	return "Processes requests using email and token credentials (placeholder implementation)"
}
