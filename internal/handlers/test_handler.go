package handlers

import (
	"context"
	"fmt"
	"time"

	"rpa-dfs-engine/internal/logger"

	"github.com/chromedp/chromedp"
)

// TestHandler handles opening test websites using Chrome automation.
// It provides functionality to test the RPA engine with known websites.
type TestHandler struct {
	url               string
	navigationTimeout time.Duration
	displayDuration   time.Duration
}

// NewTestHandler creates a new TestHandler instance with default configuration.
// It returns the Handler interface to promote loose coupling.
func NewTestHandler() Handler {
	return &TestHandler{
		url:               "https://skrooge.ai/",
		navigationTimeout: 30 * time.Second,
		displayDuration:   5 * time.Second,
	}
}

// Execute implements the Handler interface for test operations.
// It opens the configured test website using Chrome automation.
func (h *TestHandler) Execute() error {
	logger.LogInfo("=== RPA DFS Engine - Test Mode ===")
	logger.LogInfo("Opening test website: %s", h.url)

	// Create Chrome context with default options
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set navigation timeout
	ctx, cancel = context.WithTimeout(ctx, h.navigationTimeout)
	defer cancel()

	// Navigate to the test website
	if err := h.navigateToWebsite(ctx); err != nil {
		logger.LogError("Failed to navigate to website: %v", err)
		return fmt.Errorf("website navigation failed: %w", err)
	}

	logger.LogSuccess("Successfully navigated to test website")
	logger.LogInfo("Keeping browser open for %v to observe results...", h.displayDuration)

	// Keep browser open briefly to observe the result
	time.Sleep(h.displayDuration)

	logger.LogSuccess("Test mode completed successfully")
	return nil
}

// GetDescription implements the Handler interface and returns a description
// of what this handler does.
func (h *TestHandler) GetDescription() string {
	return "Opens test websites using Chrome automation for RPA engine testing"
}

// SetURL allows changing the test URL after handler creation.
// This is useful for testing different websites.
func (h *TestHandler) SetURL(url string) {
	if url != "" {
		h.url = url
	}
}

// SetNavigationTimeout allows customizing the navigation timeout.
func (h *TestHandler) SetNavigationTimeout(timeout time.Duration) {
	if timeout > 0 {
		h.navigationTimeout = timeout
	}
}

// navigateToWebsite handles the Chrome automation for website navigation.
func (h *TestHandler) navigateToWebsite(ctx context.Context) error {
	return chromedp.Run(ctx,
		chromedp.Navigate(h.url),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	)
}
