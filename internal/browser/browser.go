package browser

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"rpa-dfs-engine/internal/config"
	"rpa-dfs-engine/internal/logger"
	"rpa-dfs-engine/internal/types"

	"github.com/chromedp/chromedp"
)

func OpenBrowserWithLogin(username, password string) types.BrowserResult {
	logger.LogInfo("Starting browser for Facebook")
	logger.LogInfo("Login: %s", username)

	if !isChromeInstalled() {
		logger.LogError("Chrome not found")
		return types.BrowserResult{
			Success:   false,
			URL:       config.FACEBOOK_URL,
			Username:  username,
			Error:     "Chrome not found",
			Timestamp: time.Now().Unix(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var result string
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(config.FACEBOOK_URL),
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys(config.FACEBOOK_LOGIN_SELECTOR, username),
		chromedp.Sleep(1*time.Second),

		chromedp.SendKeys(config.FACEBOOK_PASSWORD_SELECTOR, password),
		chromedp.Sleep(1*time.Second),

		chromedp.Click(config.FACEBOOK_LOGIN_BUTTON_SELECTOR),
		chromedp.Sleep(3*time.Second),

		chromedp.Location(&result),
	)

	if err != nil {
		logger.LogError("Browser automation error: %v", err)
		return types.BrowserResult{
			Success:   false,
			URL:       config.FACEBOOK_URL,
			Username:  username,
			Error:     fmt.Sprintf("Automation error: %v", err),
			Timestamp: time.Now().Unix(),
		}
	}

	logger.LogSuccess("Browser opened and data entered successfully")
	return types.BrowserResult{
		Success:   true,
		URL:       result,
		Username:  username,
		Message:   "Facebook opened, login and password entered",
		Timestamp: time.Now().Unix(),
	}
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
