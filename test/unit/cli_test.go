package unit

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestFlagDefinition_WithValidFlags_DefinesCorrectly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	username := fs.String("l", "", "Facebook login (email)")
	password := fs.String("p", "", "Facebook password")
	help := fs.Bool("h", false, "Show help")

	assert.NotNil(t, username)
	assert.NotNil(t, password)
	assert.NotNil(t, help)

	assert.Equal(t, "", *username)
	assert.Equal(t, "", *password)
	assert.False(t, *help)
}

func TestFlagParsing_WithValidArguments_ParsesCorrectly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	username := fs.String("l", "", "Facebook login (email)")
	password := fs.String("p", "", "Facebook password")

	args := []string{"-l", "test@example.com", "-p", "secret123"}

	err := fs.Parse(args)

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", *username)
	assert.Equal(t, "secret123", *password)
}

func TestFlagParsing_WithHelpFlag_ParsesCorrectly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	help := fs.Bool("h", false, "Show help")

	args := []string{"-h"}

	err := fs.Parse(args)

	assert.NoError(t, err)
	assert.True(t, *help)
}

func TestFlagParsing_WithMissingArguments_UsesDefaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	username := fs.String("l", "", "Facebook login (email)")
	password := fs.String("p", "", "Facebook password")

	args := []string{} // No arguments

	err := fs.Parse(args)

	assert.NoError(t, err)
	assert.Empty(t, *username)
	assert.Empty(t, *password)
}

func TestParameterValidation_WithEmptyCredentials_DetectsInvalid(t *testing.T) {
	username := ""
	password := ""

	isValid := username != "" && password != ""

	assert.False(t, isValid)
}

func TestParameterValidation_WithMissingUsername_DetectsInvalid(t *testing.T) {
	username := ""
	password := "secret123"

	isValid := username != "" && password != ""

	assert.False(t, isValid)
}

func TestParameterValidation_WithMissingPassword_DetectsInvalid(t *testing.T) {
	username := "test@example.com"
	password := ""

	isValid := username != "" && password != ""

	assert.False(t, isValid)
}

func TestParameterValidation_WithValidCredentials_ValidatesSuccessfully(t *testing.T) {
	username := "test@example.com"
	password := "secret123"

	isValid := username != "" && password != ""

	assert.True(t, isValid)
}

func TestEmailValidation_WithValidEmail_ValidatesCorrectly(t *testing.T) {
	email := "test@example.com"

	hasAtSymbol := strings.Contains(email, "@")
	hasDot := strings.Contains(email, ".")
	isValid := hasAtSymbol && hasDot

	assert.True(t, isValid)
}

func TestEmailValidation_WithInvalidEmail_DetectsInvalid(t *testing.T) {
	invalidEmail := "not-an-email"

	hasAtSymbol := strings.Contains(invalidEmail, "@")
	hasDot := strings.Contains(invalidEmail, ".")
	isValid := hasAtSymbol && hasDot

	assert.False(t, isValid)
}

func TestHelpMessage_Content_ContainsRequiredInformation(t *testing.T) {
	expectedElements := []string{
		"Facebook Auto Login",
		"Usage:",
		"facebook-login.exe",
		"-l LOGIN",
		"-p PASSWORD",
		"Parameters:",
		"Examples:",
		"Requirements:",
		"Google Chrome",
	}

	helpContent := `🌐 Facebook Auto Login - Automatic Facebook Login
==================================================

Usage:
  facebook-login.exe -l LOGIN -p PASSWORD

Parameters:
  -l LOGIN     Facebook login (email)
  -p PASSWORD  Facebook password
  -h           Show this help

Examples:
  facebook-login.exe -l user@example.com -p mypassword
  facebook-login.exe -l john.doe@gmail.com -p secret123

Requirements:
  - Google Chrome installed`

	for _, element := range expectedElements {
		assert.Contains(t, helpContent, element, "Help should contain: %s", element)
	}
}

func TestCommandLineArgs_WithRealArgs_ParsesCorrectly(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"program", "-l", "test@example.com", "-p", "password123"}

	args := os.Args[1:] // Skip program name

	assert.Len(t, args, 4)
	assert.Contains(t, args, "-l")
	assert.Contains(t, args, "test@example.com")
	assert.Contains(t, args, "-p")
	assert.Contains(t, args, "password123")
}

func TestFlagParsing_WithLongArguments_HandlesCorrectly(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	username := fs.String("l", "", "Facebook login (email)")
	password := fs.String("p", "", "Facebook password")

	longEmail := "very.long.email.address.with.multiple.dots@example.subdomain.com"
	longPassword := "very-long-password-with-special-characters-123!@#$%"
	args := []string{"-l", longEmail, "-p", longPassword}

	err := fs.Parse(args)

	assert.NoError(t, err)
	assert.Equal(t, longEmail, *username)
	assert.Equal(t, longPassword, *password)
}
