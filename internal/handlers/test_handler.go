package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
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
		navigationTimeout: 7777777777 * time.Second,
	}
}

// Execute implements the Handler interface for test operations.
// Почему я не вижу этот браузер?
// Дело в том, что по умолчанию chromedp запускает браузер в headless-режиме (без интерфейса).
// Чтобы увидеть окно браузера, нужно явно указать запуск в режиме с интерфейсом (headless=false).

func (h *TestHandler) Execute() error {
	// Подавляем ошибки cookie парсинга в логах
	originalOutput := log.Writer()
	log.SetOutput(&cookieErrorFilter{originalOutput})
	defer log.SetOutput(originalOutput)

	logger.LogInfo("=== RPA DFS Engine - Test Mode ===")
	logger.LogInfo("Opening test website: %s", h.url)

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "VizDisplayCompositor,NetworkService,CookieStore,NavigationThreadingOptimizations"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("ignore-ssl-errors", true),
		chromedp.Flag("ignore-certificate-errors-spki-list", true),
		chromedp.Flag("ignore-certificate-errors-ssl-version-fallback-min", true),
		// Отключаем cookies полностью
		chromedp.Flag("disable-cookies", true),
		chromedp.Flag("disable-local-storage", true),
		chromedp.Flag("disable-session-storage", true),
		// Дополнительные флаги для стабильности
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-plugins", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-background-networking", true),
		// Флаги для решения проблем с навигацией
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-domain-reliability", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	// Устанавливаем таймаут навигации
	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, h.navigationTimeout)
	defer cancelTimeout()

	// Навигация на тестовый сайт
	if err := h.navigateToWebsite(ctxTimeout); err != nil {
		logger.LogError("Failed to navigate to website: %v", err)
		return fmt.Errorf("website navigation failed: %w", err)
	}

	logger.LogSuccess("Successfully navigated to test website")
	logger.LogInfo("Press ENTER in the console to navigate to the next link...")

	// Ждём нажатия Enter в консоли
	fmt.Print("Press ENTER to continue to the next link...")
	_, _ = fmt.Scanln()

	// Пример перехода на другой сайт после нажатия Enter
	nextURL := "https://example.com/"
	logger.LogInfo("Navigating to next website: %s", nextURL)
	if err := chromedp.Run(ctxTimeout,
		chromedp.Navigate(nextURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
	); err != nil {
		logger.LogError("Failed to navigate to next website: %v", err)
		return fmt.Errorf("next website navigation failed: %w", err)
	}
	logger.LogSuccess("Successfully navigated to next website")

	logger.LogInfo("Waiting for browser window to close...")

	// Ожидаем закрытия браузера пользователем
	<-ctx.Done()

	logger.LogInfo("Browser window closed. Exiting program.")
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

// cookieErrorFilter фильтрует ошибки связанные с cookie парсингом
type cookieErrorFilter struct {
	writer io.Writer
}

func (f *cookieErrorFilter) Write(p []byte) (n int, err error) {
	message := string(p)
	// Фильтруем ошибки связанные с cookie парсингом и навигацией
	if (strings.Contains(message, "could not unmarshal event") &&
		(strings.Contains(message, "cookiePart") || strings.Contains(message, "partitionKey"))) ||
		strings.Contains(message, "unknown ClientNavigationReason value") ||
		strings.Contains(message, "initialFrameNavigation") {
		// Игнорируем эти ошибки
		return len(p), nil
	}
	// Пропускаем остальные сообщения
	return f.writer.Write(p)
}
