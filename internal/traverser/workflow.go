package traverser

import (
	"fmt"
	"os"
	"path/filepath"
	"rpa-dfs-engine/internal/logger"
)

// getWorkflowPath returns the path to the workflow.json file in templates directory
func GetBaseWorkflowPath() (string, error) {
	// Get executable directory
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	workflowPath := filepath.Join(exeDir, "internal", "templates", "workflow.json")

	// Try relative path if absolute doesn't exist
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		workflowPath = "internal/templates/workflow.json"
	}

	// Verify file exists
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		return "", fmt.Errorf("workflow.json not found at: %s", workflowPath)
	}

	logger.LogInfo("Using workflow file: %s", workflowPath)
	return workflowPath, nil
}
