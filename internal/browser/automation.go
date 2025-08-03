package browser

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"rpa-dfs-engine/internal/logger"

	"github.com/chromedp/chromedp"
)

// Automation provides general-purpose browser automation capabilities
type Automation struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewAutomation creates a new browser automation instance
func NewAutomation() *Automation {
	logger.LogInfo("Initializing browser automation")

	if !isChromeInstalled() {
		logger.LogError("Chrome not found")
		return nil
	}

	// Create Chrome context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Run in visible mode
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	automation := &Automation{
		ctx:    ctx,
		cancel: cancel,
	}

	logger.LogSuccess("Browser automation initialized successfully")
	return automation
}

// NavigateTo navigates to the specified URL
func (a *Automation) NavigateTo(url string) error {
	logger.LogInfo("Navigating to: %s", url)

	err := chromedp.Run(a.ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Wait for page load
	)

	if err != nil {
		logger.LogError("Failed to navigate to %s: %v", url, err)
		return fmt.Errorf("failed to navigate to %s: %w", url, err)
	}

	logger.LogSuccess("Successfully navigated to: %s", url)
	return nil
}

// FillField fills a form field with the specified value
func (a *Automation) FillField(selector, value string) error {
	logger.LogInfo("Filling field %s with value: %s", selector, value)

	err := chromedp.Run(a.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Clear(selector),
		chromedp.SendKeys(selector, value),
		chromedp.Sleep(500*time.Millisecond), // Brief pause after input
	)

	if err != nil {
		logger.LogError("Failed to fill field %s: %v", selector, err)
		return fmt.Errorf("failed to fill field %s: %w", selector, err)
	}

	logger.LogSuccess("Successfully filled field: %s", selector)
	return nil
}

// ClickElement clicks on the specified element
func (a *Automation) ClickElement(selector string) error {
	logger.LogInfo("Clicking element: %s", selector)

	err := chromedp.Run(a.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Click(selector),
		chromedp.Sleep(1*time.Second), // Wait after click
	)

	if err != nil {
		logger.LogError("Failed to click element %s: %v", selector, err)
		return fmt.Errorf("failed to click element %s: %w", selector, err)
	}

	logger.LogSuccess("Successfully clicked element: %s", selector)
	return nil
}

// SendFile uploads a file to the specified file input
func (a *Automation) SendFile(selector, filePath string) error {
	logger.LogInfo("Uploading file %s to selector: %s", filePath, selector)

	err := chromedp.Run(a.ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.SendKeys(selector, filePath),
		chromedp.Sleep(1*time.Second), // Wait for file processing
	)

	if err != nil {
		logger.LogError("Failed to upload file %s to %s: %v", filePath, selector, err)
		return fmt.Errorf("failed to upload file %s to %s: %w", filePath, selector, err)
	}

	logger.LogSuccess("Successfully uploaded file %s to: %s", filePath, selector)
	return nil
}

// Close closes the browser instance
func (a *Automation) Close() error {
	logger.LogInfo("Closing browser")

	if a.cancel != nil {
		a.cancel()
	}

	logger.LogSuccess("Browser closed successfully")
	return nil
}

// WaitForElement waits for an element to be visible (utility method)
func (a *Automation) WaitForElement(selector string, timeout time.Duration) error {
	logger.LogDebug("Waiting for element: %s (timeout: %v)", selector, timeout)

	ctx, cancel := context.WithTimeout(a.ctx, timeout)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(selector, chromedp.ByQuery),
	)

	if err != nil {
		logger.LogError("Element not found within timeout %s: %v", selector, err)
		return fmt.Errorf("element not found within timeout %s: %w", selector, err)
	}

	logger.LogSuccess("Element found: %s", selector)
	return nil
}

// GetCurrentURL returns the current page URL (utility method)
func (a *Automation) GetCurrentURL() (string, error) {
	var url string
	err := chromedp.Run(a.ctx,
		chromedp.Location(&url),
	)

	if err != nil {
		logger.LogError("Failed to get current URL: %v", err)
		return "", fmt.Errorf("failed to get current URL: %w", err)
	}

	logger.LogDebug("Current URL: %s", url)
	return url, nil
}

func isChromeInstalled() bool {
	chromePaths := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe`,
	}

	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("where", "chrome")
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}
