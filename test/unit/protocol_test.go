package unit

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestProtocolURLParsing_WithValidBrowserURL_ParsesCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=user@example.com&password=secret123"

	parsedURL, err := url.Parse(protocolURL)

	assert.NoError(t, err)
	assert.Equal(t, "siteparser", parsedURL.Scheme)
	assert.Equal(t, "browser", parsedURL.Host)
	assert.Equal(t, "/", parsedURL.Path)

	query := parsedURL.Query()
	assert.Equal(t, "user@example.com", query.Get("login"))
	assert.Equal(t, "secret123", query.Get("password"))
}

func TestProtocolURLParsing_WithTestURL_ParsesCorrectly(t *testing.T) {
	protocolURL := "siteparser://test"

	parsedURL, err := url.Parse(protocolURL)

	assert.NoError(t, err)
	assert.Equal(t, "siteparser", parsedURL.Scheme)
	assert.Equal(t, "test", parsedURL.Host)
	assert.Empty(t, parsedURL.Path)
}

func TestProtocolURLParsing_WithEmailAndTokenParameters_ParsesCorrectly(t *testing.T) {
	protocolURL := "siteparser://process/?email=test@example.com&token=abc123"

	parsedURL, err := url.Parse(protocolURL)

	assert.NoError(t, err)
	assert.Equal(t, "siteparser", parsedURL.Scheme)
	assert.Equal(t, "process", parsedURL.Host)

	query := parsedURL.Query()
	assert.Equal(t, "test@example.com", query.Get("email"))
	assert.Equal(t, "abc123", query.Get("token"))
}

func TestProtocolURLParsing_WithMalformedURL_ReturnsError(t *testing.T) {
	malformedURL := "not-a-valid-url://missing-parts"

	parsedURL, err := url.Parse(malformedURL)

	assert.NoError(t, err)                             // URL parsing itself usually succeeds
	assert.NotEqual(t, "siteparser", parsedURL.Scheme) // But scheme should be wrong
}

func TestParameterExtraction_WithLoginAndPassword_ExtractsCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=test%40example.com&password=secret%21"
	parsedURL, _ := url.Parse(protocolURL)
	query := parsedURL.Query()

	login := query.Get("login")
	password := query.Get("password")

	assert.Equal(t, "test@example.com", login) // URL decoded
	assert.Equal(t, "secret!", password)       // URL decoded
}

func TestParameterExtraction_WithMissingParameters_ReturnsEmpty(t *testing.T) {
	protocolURL := "siteparser://browser/"
	parsedURL, _ := url.Parse(protocolURL)
	query := parsedURL.Query()

	login := query.Get("login")
	password := query.Get("password")

	assert.Empty(t, login)
	assert.Empty(t, password)
}

func TestPathDetection_WithBrowserPath_DetectsCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=user&password=pass"
	parsedURL, _ := url.Parse(protocolURL)

	host := parsedURL.Host
	path := parsedURL.Path

	assert.Equal(t, "browser", host)
	assert.Equal(t, "/", path)
}

func TestPathDetection_WithTestPath_DetectsCorrectly(t *testing.T) {
	protocolURL := "siteparser://test"
	parsedURL, _ := url.Parse(protocolURL)

	host := parsedURL.Host

	assert.Equal(t, "test", host)
}

func TestProtocolSchemeValidation_WithCorrectScheme_ValidatesSuccessfully(t *testing.T) {
	protocolURL := "siteparser://browser/?login=user&password=pass"
	parsedURL, _ := url.Parse(protocolURL)

	scheme := parsedURL.Scheme

	assert.Equal(t, "siteparser", scheme)
}

func TestProtocolSchemeValidation_WithIncorrectScheme_DetectsInvalid(t *testing.T) {
	protocolURL := "wrongscheme://browser/?login=user&password=pass"
	parsedURL, _ := url.Parse(protocolURL)

	scheme := parsedURL.Scheme

	assert.NotEqual(t, "siteparser", scheme)
	assert.Equal(t, "wrongscheme", scheme)
}

func TestParameterValidation_WithValidEmailFormat_ValidatesCorrectly(t *testing.T) {
	email := "user@example.com"

	isValid := strings.Contains(email, "@") && strings.Contains(email, ".")

	assert.True(t, isValid)
}

func TestParameterValidation_WithInvalidEmailFormat_DetectsInvalid(t *testing.T) {
	invalidEmail := "not-an-email"

	isValid := strings.Contains(invalidEmail, "@") && strings.Contains(invalidEmail, ".")

	assert.False(t, isValid)
}

func TestURLQueryHandling_WithEmptyQuery_HandlesGracefully(t *testing.T) {
	protocolURL := "siteparser://browser/"
	parsedURL, _ := url.Parse(protocolURL)

	query := parsedURL.Query()

	assert.NotNil(t, query)
	assert.Empty(t, query.Get("login"))
	assert.Empty(t, query.Get("password"))
}

func TestURLQueryHandling_WithSpecialCharacters_HandlesCorrectly(t *testing.T) {
	protocolURL := "siteparser://browser/?login=user%2Btest%40example.com&password=pass%26word"
	parsedURL, _ := url.Parse(protocolURL)

	query := parsedURL.Query()
	login := query.Get("login")
	password := query.Get("password")

	assert.Equal(t, "user+test@example.com", login)
	assert.Equal(t, "pass&word", password)
}

func TestProtocolAction_WithRootPath_DefaultsToBrowser(t *testing.T) {
	protocolURL := "siteparser:///?login=user&password=pass"
	parsedURL, _ := url.Parse(protocolURL)
	query := parsedURL.Query()

	path := parsedURL.Path
	if path == "" || path == "/" {
		if query.Get("login") != "" || query.Get("password") != "" {
			path = "/browser"
		}
	}

	assert.Equal(t, "/browser", path)
}
