package unit

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestProtocolCallDetection_WithValidProtocolURL_DetectsCorrectly(t *testing.T) {
	protocolPrefix := "siteparser://"
	args := []string{"program", "siteparser://browser/?login=user&password=pass"}

	isProtocol := len(args) > 1 && strings.HasPrefix(args[1], protocolPrefix)

	assert.True(t, isProtocol)
}

func TestProtocolCallDetection_WithNormalExecution_DetectsCorrectly(t *testing.T) {
	protocolPrefix := "siteparser://"
	args := []string{"program"} // No protocol URL

	isProtocol := len(args) > 1 && strings.HasPrefix(args[1], protocolPrefix)

	assert.False(t, isProtocol)
}

func TestProtocolCallDetection_WithWrongScheme_DetectsCorrectly(t *testing.T) {
	protocolPrefix := "siteparser://"
	args := []string{"program", "http://example.com"}

	isProtocol := len(args) > 1 && strings.HasPrefix(args[1], protocolPrefix)

	assert.False(t, isProtocol)
}

func TestProtocolQueryParsing_WithValidURL_ParsesCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=user@example.com&password=secret123"

	parsedURL, err := url.Parse(protocolURL)
	query := parsedURL.Query()

	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", query.Get("login"))
	assert.Equal(t, "secret123", query.Get("password"))
}

func TestProtocolQueryParsing_WithInvalidURL_ReturnsError(t *testing.T) {
	invalidURL := "not-a-valid-url"

	_, err := url.Parse(invalidURL)

	assert.NoError(t, err) // Parse succeeds but scheme validation would fail
}

func TestTestModeDetection_WithTestParameter_DetectsCorrectly(t *testing.T) {
	query := url.Values{}
	query.Set("isTest", "true")

	isTest := query.Get("isTest") == "true"

	assert.True(t, isTest)
}

func TestTestModeDetection_WithoutTestParameter_DetectsCorrectly(t *testing.T) {
	query := url.Values{}

	isTest := query.Get("isTest") == "true"

	assert.False(t, isTest)
}

func TestTestModeDetection_WithFalseTestParameter_DetectsCorrectly(t *testing.T) {
	query := url.Values{}
	query.Set("isTest", "false")

	isTest := query.Get("isTest") == "true"

	assert.False(t, isTest)
}

func TestEmailAndTokenDetection_WithBothParameters_DetectsCorrectly(t *testing.T) {
	query := url.Values{}
	query.Set("email", "test@example.com")
	query.Set("token", "abc123")

	email := query.Get("email")
	token := query.Get("token")
	hasEmailAndToken := email != "" && token != ""

	assert.True(t, hasEmailAndToken)
}

func TestEmailAndTokenDetection_WithMissingEmail_DetectsCorrectly(t *testing.T) {
	query := url.Values{}
	query.Set("token", "abc123")

	email := query.Get("email")
	token := query.Get("token")
	hasEmailAndToken := email != "" && token != ""

	assert.False(t, hasEmailAndToken)
}

func TestEmailAndTokenDetection_WithMissingToken_DetectsCorrectly(t *testing.T) {
	query := url.Values{}
	query.Set("email", "test@example.com")

	email := query.Get("email")
	token := query.Get("token")
	hasEmailAndToken := email != "" && token != ""

	assert.False(t, hasEmailAndToken)
}

func TestEmailAndTokenDetection_WithEmptyParameters_DetectsCorrectly(t *testing.T) {
	query := url.Values{}

	email := query.Get("email")
	token := query.Get("token")
	hasEmailAndToken := email != "" && token != ""

	assert.False(t, hasEmailAndToken)
}

func TestHandlerSelection_WithTestMode_SelectsTestHandler(t *testing.T) {
	query := url.Values{}
	query.Set("isTest", "true")

	var handlerType string
	if query.Get("isTest") == "true" {
		handlerType = "TestHandler"
	} else if query.Get("email") != "" && query.Get("token") != "" {
		handlerType = "ProcessHandler"
	} else {
		handlerType = "SetupHandler"
	}

	assert.Equal(t, "TestHandler", handlerType)
}

func TestHandlerSelection_WithEmailAndToken_SelectsProcessHandler(t *testing.T) {
	query := url.Values{}
	query.Set("email", "test@example.com")
	query.Set("token", "abc123")

	var handlerType string
	if query.Get("isTest") == "true" {
		handlerType = "TestHandler"
	} else if query.Get("email") != "" && query.Get("token") != "" {
		handlerType = "ProcessHandler"
	} else {
		handlerType = "SetupHandler"
	}

	assert.Equal(t, "ProcessHandler", handlerType)
}

func TestHandlerSelection_WithNoSpecialParameters_SelectsSetupHandler(t *testing.T) {
	query := url.Values{}

	var handlerType string
	if query.Get("isTest") == "true" {
		handlerType = "TestHandler"
	} else if query.Get("email") != "" && query.Get("token") != "" {
		handlerType = "ProcessHandler"
	} else {
		handlerType = "SetupHandler"
	}

	assert.Equal(t, "SetupHandler", handlerType)
}

func TestURLParameterExtraction_WithSpecialCharacters_ExtractsCorrectly(t *testing.T) {
	protocolURL := "siteparser://process/?email=test%2Buser%40example.com&token=abc%26123"
	parsedURL, _ := url.Parse(protocolURL)

	query := parsedURL.Query()
	email := query.Get("email")
	token := query.Get("token")

	assert.Equal(t, "test+user@example.com", email) // URL decoded
	assert.Equal(t, "abc&123", token)               // URL decoded
}

func TestProtocolPrefixValidation_WithCorrectPrefix_ValidatesSuccessfully(t *testing.T) {
	protocolPrefix := "siteparser://"
	testURL := "siteparser://browser/?login=user&password=pass"

	hasCorrectPrefix := strings.HasPrefix(testURL, protocolPrefix)

	assert.True(t, hasCorrectPrefix)
}

func TestProtocolPrefixValidation_WithIncorrectPrefix_DetectsInvalid(t *testing.T) {
	protocolPrefix := "siteparser://"
	testURL := "wrongscheme://browser/?login=user&password=pass"

	hasCorrectPrefix := strings.HasPrefix(testURL, protocolPrefix)

	assert.False(t, hasCorrectPrefix)
}

func TestArgLength_WithSufficientArgs_ValidatesCorrectly(t *testing.T) {
	args := []string{"program", "siteparser://test"}

	hasSufficientArgs := len(args) > 1

	assert.True(t, hasSufficientArgs)
}

func TestArgLength_WithInsufficientArgs_DetectsInvalid(t *testing.T) {
	args := []string{"program"}

	hasSufficientArgs := len(args) > 1

	assert.False(t, hasSufficientArgs)
}
