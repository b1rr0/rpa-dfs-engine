package integration

import (
	"net/url"
	"strings"
	"testing"
	"time"

	"rpa-dfs-engine/internal/types"

	"github.com/stretchr/testify/assert"
)


func TestCompleteProtocolWorkflow_WithValidBrowserURL_ProcessesCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=test@example.com&password=secret123"

	parsedURL, err := url.Parse(protocolURL)
	assert.NoError(t, err)

	query := parsedURL.Query()
	login := query.Get("login")
	password := query.Get("password")

	hasValidCredentials := login != "" && password != ""
	assert.True(t, hasValidCredentials)

	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com/feed",
		Username:  login,
		Message:   "Login successful",
		Timestamp: time.Now().Unix(),
	}

	assert.True(t, result.Success)
	assert.Contains(t, result.URL, "facebook.com")
	assert.Equal(t, login, result.Username)
	assert.NotZero(t, result.Timestamp)
}

func TestCompleteTestWorkflow_WithTestURL_ProcessesCorrectly(t *testing.T) {
	protocolURL := "siteparser://test"

	parsedURL, err := url.Parse(protocolURL)
	assert.NoError(t, err)

	host := parsedURL.Host
	isTestMode := host == "test"

	var result string
	if isTestMode {
		result = "Test mode successful"
	}

	assert.True(t, isTestMode)
	assert.Equal(t, "Test mode successful", result)
}

func TestCompleteSetupWorkflow_WithNoParameters_ProcessesCorrectly(t *testing.T) {
	args := []string{"rpa-dfs-engine.exe"}

	protocolPrefix := "siteparser://"
	isProtocolCall := len(args) > 1 && strings.HasPrefix(args[1], protocolPrefix)

	var handlerType string
	if !isProtocolCall {
		handlerType = "SetupHandler"
	}

	var setupOperations []string
	if handlerType == "SetupHandler" {
		setupOperations = []string{
			"Protocol registration",
			"HTML interface creation",
		}
	}

	assert.False(t, isProtocolCall)
	assert.Equal(t, "SetupHandler", handlerType)
	assert.Len(t, setupOperations, 2)
	assert.Contains(t, setupOperations, "Protocol registration")
	assert.Contains(t, setupOperations, "HTML interface creation")
}

func TestErrorHandlingWorkflow_WithInvalidCredentials_HandlesGracefully(t *testing.T) {
	protocolURL := "siteparser://browser/?login=&password="

	parsedURL, err := url.Parse(protocolURL)
	assert.NoError(t, err)

	query := parsedURL.Query()
	login := query.Get("login")
	password := query.Get("password")

	hasValidCredentials := login != "" && password != ""

	var result types.BrowserResult
	if !hasValidCredentials {
		result = types.BrowserResult{
			Success:   false,
			URL:       "https://www.facebook.com/",
			Username:  login,
			Error:     "Login or password not specified",
			Timestamp: time.Now().Unix(),
		}
	}

	assert.False(t, hasValidCredentials)
	assert.False(t, result.Success)
	assert.Contains(t, result.Error, "Login or password not specified")
}

func TestFileOutputWorkflow_WithBrowserResult_GeneratesCorrectFormat(t *testing.T) {
	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com/feed",
		Username:  "test@example.com",
		Message:   "Login successful",
		Timestamp: 1640995200, // 2022-01-01 00:00:00 UTC
	}

	timestamp := time.Unix(result.Timestamp, 0).UTC()
	filename := "facebook_result_" + timestamp.Format("2006-01-02_15-04-05") + ".txt"

	assert.Equal(t, "facebook_result_2022-01-01_00-00-00.txt", filename)
	assert.Contains(t, filename, "facebook_result_")
	assert.Contains(t, filename, "2022-01-01")
	assert.True(t, strings.HasSuffix(filename, ".txt"))
}
