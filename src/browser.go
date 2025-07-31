package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/chromedp/chromedp"
)

func openBrowserWithLogin(username, password string) BrowserResult {
	LogInfo("Starting browser for Facebook")
	LogInfo("Login: %s", username)

	if !isChromeInstalled() {
		LogError("Chrome not found")
		return BrowserResult{
			Success:   false,
			URL:       FACEBOOK_URL,
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
		chromedp.Navigate(FACEBOOK_URL),
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys(FACEBOOK_LOGIN_SELECTOR, username),
		chromedp.Sleep(1*time.Second),

		chromedp.SendKeys(FACEBOOK_PASSWORD_SELECTOR, password),
		chromedp.Sleep(1*time.Second),

		chromedp.Click(FACEBOOK_LOGIN_BUTTON_SELECTOR),
		chromedp.Sleep(3*time.Second),

		chromedp.Location(&result),
	)

	if err != nil {
		LogError("Browser automation error: %v", err)
		return BrowserResult{
			Success:   false,
			URL:       FACEBOOK_URL,
			Username:  username,
			Error:     fmt.Sprintf("Automation error: %v", err),
			Timestamp: time.Now().Unix(),
		}
	}

	LogSuccess("Browser opened and data entered successfully")
	return BrowserResult{
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
