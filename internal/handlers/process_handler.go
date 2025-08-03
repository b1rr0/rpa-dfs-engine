package handlers

import (
	"fmt"
	"net/url"
	"os"

	"rpa-dfs-engine/internal/browser"
	"rpa-dfs-engine/internal/logger"
	"rpa-dfs-engine/internal/traverser"
)

// ProcessHandler handles processing requests that include email and token credentials.
// Executes a workflow using the traverser system with workflow.json from templates.
type ProcessHandler struct {
	email   string
	token   string
	website string
}

// NewProcessHandler creates a new ProcessHandler instance.
// It returns the Handler interface to promote loose coupling.
func NewProcessHandler() Handler {
	return &ProcessHandler{}
}

// Execute implements the Handler interface for processing with credentials.
// Single entry point that loads and executes workflow.json from templates directory.
func (h *ProcessHandler) Execute() error {
	logger.LogInfo("=== RPA DFS Engine - Process Mode ===")
	logger.LogInfo("Starting traverser workflow execution")

	// Extract parameters from command line arguments (assuming protocol URL format)
	if err := h.extractParameters(); err != nil {
		logger.LogError("Failed to extract parameters: %v", err)
		return fmt.Errorf("parameter extraction failed: %w", err)
	}

	// Create context data
	contextData := h.createContextData()

	// Execute workflow using traverser
	if err := h.executeWorkflow(contextData); err != nil {
		logger.LogError("Workflow execution failed: %v", err)
		return fmt.Errorf("workflow execution failed: %w", err)
	}

	logger.LogSuccess("Process handler completed successfully")
	return nil
}

// GetDescription implements the Handler interface and returns a description
// of what this handler does.
func (h *ProcessHandler) GetDescription() string {
	return "Executes traverser workflow using templates/workflow.json with email and token credentials"
}

// extractParameters extracts email, token, and website from command line arguments
func (h *ProcessHandler) extractParameters() error {
	args := os.Args
	if len(args) < 2 {
		return fmt.Errorf("no protocol URL provided")
	}

	// Parse protocol URL to extract parameters
	protocolURL := args[1]
	parsedURL, err := url.Parse(protocolURL)
	if err != nil {
		return fmt.Errorf("failed to parse protocol URL: %w", err)
	}

	query := parsedURL.Query()
	h.email = query.Get("email")
	h.token = query.Get("token")

	logger.LogInfo("Extracted parameters - Email: %s, Website: %s", h.email, h.website)
	return nil
}

// createContextData creates the context data for template resolution
func (h *ProcessHandler) createContextData() map[string]interface{} {
	return map[string]interface{}{
		"user": map[string]interface{}{
			"email": h.email,
		},
	}
}

// executeWorkflow executes the workflow using the traverser system
func (h *ProcessHandler) executeWorkflow(contextData map[string]interface{}) error {
	// Create traverser instance
	t := traverser.NewWithBrowser(browser.NewAutomation())
	defer t.Close()
	// Get workflow file path
	workflowPath, err := traverser.GetBaseWorkflowPath()
	if err != nil {
		logger.LogError("Failed to locate workflow file: %v", err)
		return fmt.Errorf("workflow file not found: %w", err)
	}
	context := traverser.NewContext(contextData)
	err = t.ExecuteWorkflow(workflowPath, context)
	if err != nil {
		return fmt.Errorf("failed to execute workflow: %w", err)
	}

	logger.LogSuccess("Workflow executed successfully")
	return nil
}
