package unit

import (
	"encoding/json"
	"testing"
	"time"

	"rpa-dfs-engine/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestBrowserResult_WithSuccessfulResult_MarshalesToCorrectJSON(t *testing.T) {
	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com/feed",
		Username:  "test@example.com",
		Message:   "Login successful",
		Timestamp: 1234567890,
	}

	jsonData, err := json.Marshal(result)

	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":true`)
	assert.Contains(t, string(jsonData), `"url":"https://facebook.com/feed"`)
	assert.Contains(t, string(jsonData), `"username":"test@example.com"`)
	assert.Contains(t, string(jsonData), `"message":"Login successful"`)
	assert.Contains(t, string(jsonData), `"timestamp":1234567890`)
	assert.NotContains(t, string(jsonData), `"error"`)
}

func TestBrowserResult_WithErrorResult_MarshalesToCorrectJSON(t *testing.T) {
	result := types.BrowserResult{
		Success:   false,
		URL:       "https://facebook.com",
		Username:  "test@example.com",
		Error:     "Chrome not found",
		Timestamp: 1234567890,
	}

	jsonData, err := json.Marshal(result)

	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":false`)
	assert.Contains(t, string(jsonData), `"error":"Chrome not found"`)
}

func TestBrowserResult_WithValidJSON_UnmarshalCorrectly(t *testing.T) {
	jsonData := `{
		"success": true,
		"url": "https://facebook.com/feed",
		"username": "test@example.com",
		"message": "Login successful",
		"timestamp": 1234567890
	}`

	var result types.BrowserResult
	err := json.Unmarshal([]byte(jsonData), &result)

	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "https://facebook.com/feed", result.URL)
	assert.Equal(t, "test@example.com", result.Username)
	assert.Equal(t, "Login successful", result.Message)
	assert.Equal(t, int64(1234567890), result.Timestamp)
	assert.Empty(t, result.Error)
}

func TestBrowserResult_WithInvalidJSON_ReturnsUnmarshalError(t *testing.T) {
	invalidJSON := `{"success": "invalid_boolean"}`

	var result types.BrowserResult
	err := json.Unmarshal([]byte(invalidJSON), &result)

	assert.Error(t, err)
}

func TestBrowserResult_WithCurrentTimestamp_SetsCorrectTimestamp(t *testing.T) {
	beforeTime := time.Now().Unix()

	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com",
		Username:  "test@example.com",
		Message:   "Test",
		Timestamp: time.Now().Unix(),
	}

	afterTime := time.Now().Unix()

	assert.True(t, result.Timestamp >= beforeTime)
	assert.True(t, result.Timestamp <= afterTime)
}

func TestBrowserResult_WithEmptyFields_HandlesGracefully(t *testing.T) {
	result := types.BrowserResult{}

	jsonData, err := json.Marshal(result)

	assert.NoError(t, err)
	assert.Contains(t, string(jsonData), `"success":false`)
	assert.Contains(t, string(jsonData), `"timestamp":0`)
}
