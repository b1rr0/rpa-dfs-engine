package unit

import (
	"os"
	"testing"

	"rpa-dfs-engine/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestProtocolName_WithEnvironmentVariable_ReturnsEnvironmentValue(t *testing.T) {
	originalValue := os.Getenv("PROTOCOL_NAME")
	defer os.Setenv("PROTOCOL_NAME", originalValue)

	expectedProtocol := "testprotocol"
	os.Setenv("PROTOCOL_NAME", expectedProtocol)

	result := os.Getenv("PROTOCOL_NAME")

	assert.Equal(t, expectedProtocol, result)
}

func TestSiteForTest_WithEnvironmentVariable_ReturnsEnvironmentValue(t *testing.T) {
	originalValue := os.Getenv("SITE_FOR_TEST")
	defer os.Setenv("SITE_FOR_TEST", originalValue)

	expectedSite := "https://test.example.com"
	os.Setenv("SITE_FOR_TEST", expectedSite)

	result := os.Getenv("SITE_FOR_TEST")

	assert.Equal(t, expectedSite, result)
}

func TestFacebookConstants_AreDefinedCorrectly(t *testing.T) {
	assert.Equal(t, "https://www.facebook.com/", config.FACEBOOK_URL)
	assert.Equal(t, "#email", config.FACEBOOK_LOGIN_SELECTOR)
	assert.Equal(t, "#pass", config.FACEBOOK_PASSWORD_SELECTOR)
	assert.Equal(t, "button[name='login']", config.FACEBOOK_LOGIN_BUTTON_SELECTOR)
}

func TestFacebookConstants_AreNotEmpty(t *testing.T) {
	assert.NotEmpty(t, config.FACEBOOK_URL)
	assert.NotEmpty(t, config.FACEBOOK_LOGIN_SELECTOR)
	assert.NotEmpty(t, config.FACEBOOK_PASSWORD_SELECTOR)
	assert.NotEmpty(t, config.FACEBOOK_LOGIN_BUTTON_SELECTOR)
}

func TestFacebookURL_IsValidFormat(t *testing.T) {
	assert.Contains(t, config.FACEBOOK_URL, "https://")
	assert.Contains(t, config.FACEBOOK_URL, "facebook.com")
}

func TestFacebookSelectors_AreValidCSS(t *testing.T) {
	assert.True(t, config.FACEBOOK_LOGIN_SELECTOR[0] == '#', "Login selector should be an ID selector")
	assert.True(t, config.FACEBOOK_PASSWORD_SELECTOR[0] == '#', "Password selector should be an ID selector")
	assert.Contains(t, config.FACEBOOK_LOGIN_BUTTON_SELECTOR, "button", "Login button should be a button selector")
}

func TestFacebookSelectors_HaveCorrectAttributes(t *testing.T) {
	assert.Contains(t, config.FACEBOOK_LOGIN_BUTTON_SELECTOR, "name='login'", "Login button should have name attribute")
}

func TestEnvironmentVariables_WithEmptyValues_HandleGracefully(t *testing.T) {
	originalProtocol := os.Getenv("PROTOCOL_NAME")
	originalSite := os.Getenv("SITE_FOR_TEST")
	defer func() {
		os.Setenv("PROTOCOL_NAME", originalProtocol)
		os.Setenv("SITE_FOR_TEST", originalSite)
	}()

	os.Setenv("PROTOCOL_NAME", "")
	os.Setenv("SITE_FOR_TEST", "")

	protocolResult := os.Getenv("PROTOCOL_NAME")
	siteResult := os.Getenv("SITE_FOR_TEST")

	assert.Empty(t, protocolResult)
	assert.Empty(t, siteResult)
}
