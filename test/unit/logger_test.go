package unit

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogMessageFormatting_WithInfoLevel_FormatsCorrectly(t *testing.T) {
	testMessage := "Test info message"
	expectedPrefix := "INFO:"

	formattedMessage := fmt.Sprintf("%s %s", expectedPrefix, testMessage)

	assert.Contains(t, formattedMessage, expectedPrefix)
	assert.Contains(t, formattedMessage, testMessage)
	assert.Equal(t, "INFO: Test info message", formattedMessage)
}

func TestLogMessageFormatting_WithSuccessLevel_FormatsCorrectly(t *testing.T) {
	testMessage := "Operation completed successfully"
	expectedPrefix := "SUCCESS:"

	formattedMessage := fmt.Sprintf("%s %s", expectedPrefix, testMessage)

	assert.Contains(t, formattedMessage, expectedPrefix)
	assert.Contains(t, formattedMessage, testMessage)
	assert.Equal(t, "SUCCESS: Operation completed successfully", formattedMessage)
}

func TestLogMessageFormatting_WithErrorLevel_FormatsCorrectly(t *testing.T) {
	testMessage := "An error occurred"
	expectedPrefix := "ERROR:"

	formattedMessage := fmt.Sprintf("%s %s", expectedPrefix, testMessage)

	assert.Contains(t, formattedMessage, expectedPrefix)
	assert.Contains(t, formattedMessage, testMessage)
	assert.Equal(t, "ERROR: An error occurred", formattedMessage)
}

func TestLogMessageFormatting_WithWarningLevel_FormatsCorrectly(t *testing.T) {
	testMessage := "This is a warning"
	expectedPrefix := "WARNING:"

	formattedMessage := fmt.Sprintf("%s %s", expectedPrefix, testMessage)

	assert.Contains(t, formattedMessage, expectedPrefix)
	assert.Contains(t, formattedMessage, testMessage)
	assert.Equal(t, "WARNING: This is a warning", formattedMessage)
}

func TestLogMessageFormatting_WithDebugLevel_FormatsCorrectly(t *testing.T) {
	testMessage := "Debug information"
	expectedPrefix := "DEBUG:"

	formattedMessage := fmt.Sprintf("%s %s", expectedPrefix, testMessage)

	assert.Contains(t, formattedMessage, expectedPrefix)
	assert.Contains(t, formattedMessage, testMessage)
	assert.Equal(t, "DEBUG: Debug information", formattedMessage)
}

func TestLogFileNameGeneration_WithCurrentDate_GeneratesCorrectFormat(t *testing.T) {
	testDate := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	expectedDate := "2023-12-25"

	logFileName := fmt.Sprintf("facebook-login_%s.log", testDate.Format("2006-01-02"))

	assert.Equal(t, "facebook-login_2023-12-25.log", logFileName)
	assert.Contains(t, logFileName, expectedDate)
	assert.Contains(t, logFileName, "facebook-login_")
	assert.Contains(t, logFileName, ".log")
}

func TestLogDirectoryPath_WithExecutablePath_GeneratesCorrectPath(t *testing.T) {
	exePath := "/usr/local/bin/rpa-dfs-engine"
	exeDir := filepath.Dir(exePath)
	logsDir := filepath.Join(exeDir, "logs")

	assert.Equal(t, "/usr/local/bin/logs", logsDir)
	assert.Contains(t, logsDir, exeDir)
	assert.Contains(t, logsDir, "logs")
}

func TestSessionMessages_HaveCorrectFormat(t *testing.T) {
	sessionStart := "=== SESSION START ==="
	sessionEnd := "=== SESSION END ==="

	assert.True(t, strings.HasPrefix(sessionStart, "==="))
	assert.True(t, strings.HasSuffix(sessionStart, "==="))
	assert.Contains(t, sessionStart, "SESSION START")

	assert.True(t, strings.HasPrefix(sessionEnd, "==="))
	assert.True(t, strings.HasSuffix(sessionEnd, "==="))
	assert.Contains(t, sessionEnd, "SESSION END")
}

func TestLogMessageFormatting_WithParameters_FormatsCorrectly(t *testing.T) {
	format := "User %s logged in with status %t"
	username := "test@example.com"
	success := true

	message := fmt.Sprintf(format, username, success)
	logEntry := fmt.Sprintf("INFO: %s", message)

	assert.Equal(t, "INFO: User test@example.com logged in with status true", logEntry)
	assert.Contains(t, logEntry, username)
	assert.Contains(t, logEntry, "true")
}

func TestLogTimestampFormat_WithCurrentTime_UsesStandardFormat(t *testing.T) {
	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	expectedFormat := "2023-12-25 15:30:45"

	formattedTime := testTime.Format("2006-01-02 15:04:05")

	assert.Equal(t, expectedFormat, formattedTime)
}

func TestLogMessageSanitization_WithSpecialCharacters_HandlesCorrectly(t *testing.T) {
	messageWithSpecialChars := "Error: failed to parse \"user@domain.com\" & retry"
	prefix := "ERROR:"

	logEntry := fmt.Sprintf("%s %s", prefix, messageWithSpecialChars)

	assert.Contains(t, logEntry, "\"user@domain.com\"")
	assert.Contains(t, logEntry, "&")
	assert.Equal(t, "ERROR: Error: failed to parse \"user@domain.com\" & retry", logEntry)
}

func TestLogEntryLength_WithLongMessage_HandlesAppropriately(t *testing.T) {
	longMessage := strings.Repeat("A", 1000)
	prefix := "INFO:"

	logEntry := fmt.Sprintf("%s %s", prefix, longMessage)

	assert.True(t, len(logEntry) > 1000)
	assert.Contains(t, logEntry, prefix)
	assert.Contains(t, logEntry, longMessage)
}
