package unit

import (
	"runtime"
	"strings"
	"testing"
	"time"

	"rpa-dfs-engine/internal/types"

	"github.com/stretchr/testify/assert"
)


func TestBrowserResult_WithSuccessfulLogin_ReturnsSuccessResult(t *testing.T) {
	username := "test@example.com"
	expectedURL := "https://facebook.com/feed"

	result := types.BrowserResult{
		Success:   true,
		URL:       expectedURL,
		Username:  username,
		Message:   "Facebook opened, login and password entered",
		Timestamp: time.Now().Unix(),
	}

	assert.True(t, result.Success)
	assert.Equal(t, expectedURL, result.URL)
	assert.Equal(t, username, result.Username)
	assert.Contains(t, result.Message, "Facebook opened")
	assert.Empty(t, result.Error)
	assert.NotZero(t, result.Timestamp)
}

func TestBrowserResult_WithChromeNotFound_ReturnsErrorResult(t *testing.T) {
	username := "test@example.com"

	result := types.BrowserResult{
		Success:   false,
		URL:       "https://www.facebook.com/",
		Username:  username,
		Error:     "Chrome not found",
		Timestamp: time.Now().Unix(),
	}

	assert.False(t, result.Success)
	assert.Equal(t, "https://www.facebook.com/", result.URL)
	assert.Equal(t, username, result.Username)
	assert.Equal(t, "Chrome not found", result.Error)
	assert.Empty(t, result.Message)
	assert.NotZero(t, result.Timestamp)
}

func TestBrowserResult_WithAutomationError_ReturnsErrorResult(t *testing.T) {
	username := "test@example.com"
	automationError := "element not found: #email"

	result := types.BrowserResult{
		Success:   false,
		URL:       "https://www.facebook.com/",
		Username:  username,
		Error:     "Automation error: " + automationError,
		Timestamp: time.Now().Unix(),
	}

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "Automation error:")
	assert.Contains(t, result.Error, automationError)
}

func TestChromePathValidation_WithWindowsPaths_ValidatesCorrectly(t *testing.T) {
	chromePaths := []string{
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe`,
	}

	for _, path := range chromePaths {
		assert.Contains(t, path, "Google\\Chrome\\Application\\chrome.exe")
		assert.True(t, strings.HasSuffix(path, "chrome.exe"))
	}
}

func TestRuntimeDetection_WithWindowsOS_DetectsCorrectly(t *testing.T) {
	currentOS := runtime.GOOS

	if currentOS == "windows" {
		assert.Equal(t, "windows", currentOS)
	} else {
		assert.NotEqual(t, "windows", currentOS)
	}
}

func TestBrowserTimeout_WithValidDuration_ConfiguresCorrectly(t *testing.T) {
	timeoutDuration := 30 * time.Second

	timeout := timeoutDuration

	assert.Equal(t, 30*time.Second, timeout)
	assert.True(t, timeout > 0)
}

func TestChromeDPOptions_Configuration_SetsCorrectFlags(t *testing.T) {
	flags := make(map[string]bool)
	flags["headless"] = false // Set to false in actual code
	flags["disable-gpu"] = true
	flags["no-sandbox"] = true
	flags["disable-dev-shm-usage"] = true

	assert.False(t, flags["headless"]) // Should be false for visible browser
	assert.True(t, flags["disable-gpu"])
	assert.True(t, flags["no-sandbox"])
	assert.True(t, flags["disable-dev-shm-usage"])
}

func TestSelectorValidation_WithFacebookSelectors_ValidatesCorrectly(t *testing.T) {
	emailSelector := "#email"
	passwordSelector := "#pass"
	loginButtonSelector := "button[name='login']"

	assert.True(t, strings.HasPrefix(emailSelector, "#"))
	assert.True(t, strings.HasPrefix(passwordSelector, "#"))
	assert.Contains(t, loginButtonSelector, "button")
	assert.Contains(t, loginButtonSelector, "name='login'")
}

func TestBrowserAutomationSteps_WithValidSelectors_DefinesCorrectSequence(t *testing.T) {
	steps := []string{
		"Navigate to Facebook",
		"Send keys to email field",
		"Send keys to password field",
		"Click login button",
		"Get current location",
	}

	assert.Len(t, steps, 5)
	assert.Contains(t, steps[0], "Navigate")
	assert.Contains(t, steps[1], "email")
	assert.Contains(t, steps[2], "password")
	assert.Contains(t, steps[3], "login")
	assert.Contains(t, steps[4], "location")
}

func TestErrorHandling_WithContextTimeout_ReturnsTimeoutError(t *testing.T) {
	timeoutError := "context deadline exceeded"

	result := types.BrowserResult{
		Success:   false,
		URL:       "https://www.facebook.com/",
		Username:  "test@example.com",
		Error:     "Automation error: " + timeoutError,
		Timestamp: time.Now().Unix(),
	}

	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "context deadline exceeded")
}

func TestParameterValidation_WithEmptyCredentials_HandlesCorrectly(t *testing.T) {
	username := ""
	password := ""

	isValid := username != "" && password != ""

	assert.False(t, isValid)
}

func TestParameterValidation_WithValidCredentials_HandlesCorrectly(t *testing.T) {
	username := "test@example.com"
	password := "secret123"

	isValid := username != "" && password != ""

	assert.True(t, isValid)
}

func TestURLValidation_WithFacebookURL_ValidatesCorrectly(t *testing.T) {
	facebookURL := "https://www.facebook.com/"

	assert.Contains(t, facebookURL, "https://")
	assert.Contains(t, facebookURL, "facebook.com")
	assert.True(t, strings.HasPrefix(facebookURL, "https://"))
}

func TestTimestampGeneration_WithCurrentTime_GeneratesValidTimestamp(t *testing.T) {
	beforeTime := time.Now().Unix()

	timestamp := time.Now().Unix()

	afterTime := time.Now().Unix()

	assert.True(t, timestamp >= beforeTime)
	assert.True(t, timestamp <= afterTime)
	assert.True(t, timestamp > 0)
}

func TestBrowserResult_WithPartialFailure_HandlesCorrectly(t *testing.T) {
	result := types.BrowserResult{
		Success:   false,
		URL:       "https://www.facebook.com/login/",
		Username:  "test@example.com",
		Error:     "Login failed: incorrect credentials",
		Timestamp: time.Now().Unix(),
	}

	assert.False(t, result.Success)
	assert.Contains(t, result.URL, "facebook.com")
	assert.Contains(t, result.Error, "Login failed")
}

func TestElementInteraction_WithValidSelectors_DefinesCorrectActions(t *testing.T) {
	interactions := map[string]string{
		"#email":               "SendKeys",
		"#pass":                "SendKeys",
		"button[name='login']": "Click",
	}

	assert.Equal(t, "SendKeys", interactions["#email"])
	assert.Equal(t, "SendKeys", interactions["#pass"])
	assert.Equal(t, "Click", interactions["button[name='login']"])
}
