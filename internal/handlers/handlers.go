package handlers

import (
	"os"
	"strings"
)

// GetHandler returns the appropriate handler based on command line arguments.
// It analyzes the execution context and arguments to determine which handler to use:
// - TestHandler: when run from siteparser with isTest=true
// - ProcessHandler: when run from siteparser with email and token
// - SetupHandler: default fallback for initial setup
func GetHandler() Handler {
	args := os.Args

	// Check if running from siteparser context
	isFromSiteparser := strings.Contains(args[0], "/siteparser")

	switch {
	case isFromSiteparser && containsArg(args, "isTest=true"):
		return NewTestHandler()
	case isFromSiteparser && hasEmailAndToken(args):
		return NewProcessHandler()
	default:
		return NewSetupHandler()
	}
}

// containsArg checks if any command line argument contains the specified substring.
func containsArg(args []string, target string) bool {
	for _, arg := range args {
		if strings.Contains(arg, target) {
			return true
		}
	}
	return false
}

// hasEmailAndToken checks if command line arguments contain both email and token parameters.
// It supports multiple formats: flags (-e, --email) and key-value pairs (email=).
func hasEmailAndToken(args []string) bool {
	emailFlags := []string{"-e", "--email", "email="}
	tokenFlags := []string{"-t", "--token", "token="}

	hasEmail := hasAnyFlag(args, emailFlags)
	hasToken := hasAnyFlag(args, tokenFlags)

	return hasEmail && hasToken
}

// hasAnyFlag checks if any of the provided flags exist in the command line arguments.
func hasAnyFlag(args []string, flags []string) bool {
	for _, flag := range flags {
		if containsArg(args, flag) {
			return true
		}
	}
	return false
}
