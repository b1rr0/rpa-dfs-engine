package handlers

import (
	"rpa-dfs-engine/internal/html"
	"rpa-dfs-engine/internal/logger"
	"rpa-dfs-engine/internal/protocol"
)

// SetupHandler handles the setup/first-run functionality
type SetupHandler struct{}

// NewSetupHandler creates a new setup handler instance
func NewSetupHandler() Handler {
	return &SetupHandler{}
}

// Execute implements the Handler interface
func (h *SetupHandler) Execute() error {
	logger.LogInfo("=== RPA DFS Engine - Setup Mode ===")
	logger.LogInfo("Checking protocol registration")

	html.CreateHtmlInterface()

	logger.LogInfo("Registering siteparser:// protocol...")

	if protocol.RegisterProtocol() {
		logger.LogSuccess("Protocol registered successfully")
		logger.LogSuccess("HTML interface created")
	} else {
		logger.LogWarning("Failed to register protocol")
		logger.LogInfo("💡 To register protocol run: run-as-admin.bat")
		logger.LogInfo("📄 Application will continue without protocol")
		logger.LogSuccess("HTML interface created")
	}

	return nil
}

// GetDescription implements the Handler interface
func (h *SetupHandler) GetDescription() string {
	return "Sets up protocol registration and HTML interface"
}
