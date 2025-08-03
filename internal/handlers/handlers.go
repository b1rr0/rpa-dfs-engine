package handlers

import (
	"net/url"
	"os"
	"strings"

	"rpa-dfs-engine/internal/config"
	"rpa-dfs-engine/internal/logger"
)

var protocolPrefix = config.PROTOCOL_NAME + "://"

func GetHandler() Handler {
	args := os.Args
	logger.LogInfo("args: %v", args)
	if !isProtocolCall(args) {
		return NewSetupHandler()
	}

	query, err := parseProtocolQuery(args[1])
	if err != nil {
		logger.LogError("Error parsing protocol URL: %v", err)
		return NewSetupHandler()
	}

	logger.LogInfo("Parsed query parameters: %v", query)

	if isTestMode(query) {
		return NewTestHandler()
	}

	if hasEmailAndToken(query) {
		return NewProcessHandler()
	}

	return NewSetupHandler()
}

func isProtocolCall(args []string) bool {
	return len(args) > 1 && strings.HasPrefix(args[1], protocolPrefix)
}

func parseProtocolQuery(protocolURL string) (url.Values, error) {
	parsedURL, err := url.Parse(protocolURL)
	if err != nil {
		return nil, err
	}
	return parsedURL.Query(), nil
}

func isTestMode(query url.Values) bool {
	return query.Get("isTest") == "true"
}

func hasEmailAndToken(query url.Values) bool {
	email := query.Get("email")
	token := query.Get("token")
	return email != "" && token != ""
}
