package unit

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"rpa-dfs-engine/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestSaveBrowserResultToFile_WithSuccessfulResult_CreatesCorrectContent(t *testing.T) {
	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com/feed",
		Username:  "test@example.com",
		Message:   "Login successful",
		Timestamp: 1234567890,
	}

	expectedContent := fmt.Sprintf(`=== Facebook Auto Login Result ===
Time: %s
Success: %t
URL: %s
Login: %s
Message: %s
`,
		time.Unix(result.Timestamp, 0).Format("2006-01-02 15:04:05"),
		result.Success,
		result.URL,
		result.Username,
		result.Message,
	)
	expectedContent += "\n=== END OF RESULT ===\n"

	assert.Contains(t, expectedContent, "=== Facebook Auto Login Result ===")
	assert.Contains(t, expectedContent, "Success: true")
	assert.Contains(t, expectedContent, "URL: https://facebook.com/feed")
	assert.Contains(t, expectedContent, "Login: test@example.com")
	assert.Contains(t, expectedContent, "Message: Login successful")
	assert.Contains(t, expectedContent, "=== END OF RESULT ===")
	assert.NotContains(t, expectedContent, "Error:")
}

func TestSaveBrowserResultToFile_WithErrorResult_IncludesErrorSection(t *testing.T) {
	result := types.BrowserResult{
		Success:   false,
		URL:       "https://facebook.com",
		Username:  "test@example.com",
		Error:     "Chrome not found",
		Timestamp: 1234567890,
	}

	expectedContent := fmt.Sprintf(`=== Facebook Auto Login Result ===
Time: %s
Success: %t
URL: %s
Login: %s
Message: %s
`,
		time.Unix(result.Timestamp, 0).Format("2006-01-02 15:04:05"),
		result.Success,
		result.URL,
		result.Username,
		result.Message,
	)

	if result.Error != "" {
		expectedContent += fmt.Sprintf("Error: %s\n", result.Error)
	}

	expectedContent += "\n=== END OF RESULT ===\n"

	assert.Contains(t, expectedContent, "Success: false")
	assert.Contains(t, expectedContent, "Error: Chrome not found")
}

func TestFileNameGeneration_WithTimestamp_CreatesCorrectFormat(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)

	timestamp := testTime.Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("facebook_result_%s.txt", timestamp)

	assert.Equal(t, "facebook_result_2023-12-25_15-30-45.txt", filename)
	assert.Contains(t, filename, "facebook_result_")
	assert.Contains(t, filename, ".txt")
}

func TestFileContent_WithEmptyMessage_HandlesGracefully(t *testing.T) {
	result := types.BrowserResult{
		Success:   true,
		URL:       "https://facebook.com",
		Username:  "test@example.com",
		Message:   "",
		Timestamp: 1234567890,
	}

	content := fmt.Sprintf(`=== Facebook Auto Login Result ===
Time: %s
Success: %t
URL: %s
Login: %s
Message: %s
`,
		time.Unix(result.Timestamp, 0).Format("2006-01-02 15:04:05"),
		result.Success,
		result.URL,
		result.Username,
		result.Message,
	)

	assert.Contains(t, content, "Message: ")
	assert.NotContains(t, content, "Message: \n\n")
}

func TestTimestampFormatting_WithValidTimestamp_FormatsCorrectly(t *testing.T) {
	timestamp := int64(1609459200)

	formattedTime := time.Unix(timestamp, 0).UTC().Format("2006-01-02 15:04:05")

	assert.Equal(t, "2021-01-01 00:00:00", formattedTime)
}

func TestFilePathGeneration_WithExecutablePath_GeneratesValidPath(t *testing.T) {
	exePath := "/usr/local/bin/rpa-dfs-engine"
	exeDir := filepath.Dir(exePath)
	timestamp := "2023-12-25_15-30-45"
	filename := fmt.Sprintf("facebook_result_%s.txt", timestamp)

	filePath := filepath.Join(exeDir, filename)

	assert.Equal(t, "/usr/local/bin/facebook_result_2023-12-25_15-30-45.txt", filePath)
	assert.Contains(t, filePath, exeDir)
	assert.Contains(t, filePath, filename)
}

func TestContentGeneration_WithSpecialCharacters_HandlesCorrectly(t *testing.T) {
	result := types.BrowserResult{
		Success:   false,
		URL:       "https://facebook.com/login?error=invalid",
		Username:  "test+user@example.com",
		Error:     "Login failed: invalid credentials",
		Message:   "Error with special chars: & < > \"",
		Timestamp: 1234567890,
	}

	content := fmt.Sprintf("Login: %s\nError: %s\nMessage: %s",
		result.Username, result.Error, result.Message)

	assert.Contains(t, content, "test+user@example.com")
	assert.Contains(t, content, "Login failed: invalid credentials")
	assert.Contains(t, content, "& < > \"")
}
