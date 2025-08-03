package handlers

import (
	"fmt"

	"rpa-dfs-engine/internal/browser"
	"rpa-dfs-engine/internal/config"
	"rpa-dfs-engine/internal/logger"
)

type TestHandler struct {
	url     string
	browser *browser.Automation
}

func NewTestHandler() Handler {
	return &TestHandler{
		url: config.SITE_FOR_TEST,
	}
}

func (h *TestHandler) Execute() error {
	logger.LogInfo("=== RPA DFS Engine - Test Mode ===")
	logger.LogInfo("Opening test website: %s", h.url)

	// Initialize browser automation
	h.browser = browser.NewAutomation()
	if h.browser == nil {
		logger.LogError("Failed to initialize browser automation")
		return fmt.Errorf("browser initialization failed")
	}
	defer h.browser.Close()

	// Navigate to the test website
	if err := h.browser.NavigateTo(h.url); err != nil {
		logger.LogError("Failed to navigate to website: %v", err)
		return fmt.Errorf("website navigation failed: %w", err)
	}

	logger.LogSuccess("Successfully navigated to test website")
	logger.LogInfo("Press ENTER in the console to navigate to the next link...")

	// Wait for user input
	fmt.Print("Press ENTER to continue to the next link...")
	_, _ = fmt.Scanln()

	// Navigate to the next website
	nextURL := "https://example.com/"
	logger.LogInfo("Navigating to next website: %s", nextURL)

	if err := h.browser.NavigateTo(nextURL); err != nil {
		logger.LogError("Failed to navigate to next website: %v", err)
		return fmt.Errorf("next website navigation failed: %w", err)
	}

	logger.LogSuccess("Successfully navigated to next website")
	logger.LogInfo("Browser will remain open. Close the browser window to exit.")

	// Wait for user to close browser manually
	fmt.Print("Close the browser window to exit the program...")
	_, _ = fmt.Scanln()

	logger.LogInfo("Test handler execution completed")
	return nil
}

func (h *TestHandler) GetDescription() string {
	return "Opens test websites using browser automation for RPA engine testing"
}

// SetURL allows customizing the test URL
func (h *TestHandler) SetURL(url string) {
	if url != "" {
		h.url = url
	}
}

// NavigateToMultipleSites can be used to test navigation to multiple sites
func (h *TestHandler) NavigateToMultipleSites(urls []string) error {
	if h.browser == nil {
		return fmt.Errorf("browser not initialized")
	}

	for i, url := range urls {
		logger.LogInfo("Navigating to site %d/%d: %s", i+1, len(urls), url)

		if err := h.browser.NavigateTo(url); err != nil {
			logger.LogError("Failed to navigate to %s: %v", url, err)
			return err
		}

		logger.LogSuccess("Successfully loaded: %s", url)

		if i < len(urls)-1 {
			fmt.Printf("Press ENTER to continue to next site (%d/%d)...", i+2, len(urls))
			_, _ = fmt.Scanln()
		}
	}

	return nil
}

// TestFormInteraction demonstrates form filling capabilities
func (h *TestHandler) TestFormInteraction() error {
	if h.browser == nil {
		return fmt.Errorf("browser not initialized")
	}

	logger.LogInfo("Testing form interaction capabilities")

	// Example form interactions (would need actual form elements)
	testData := map[string]string{
		"#email":   "test@example.com",
		"#name":    "Test User",
		"#message": "This is a test message",
	}

	for selector, value := range testData {
		logger.LogInfo("Attempting to fill field: %s", selector)
		if err := h.browser.FillField(selector, value); err != nil {
			logger.LogWarning("Could not fill field %s (element may not exist): %v", selector, err)
		} else {
			logger.LogSuccess("Successfully filled field: %s", selector)
		}
	}

	return nil
}
