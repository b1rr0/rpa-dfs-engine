package protocol

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"rpa-dfs-engine/internal/browser"
	"rpa-dfs-engine/internal/fileutils"
	"rpa-dfs-engine/internal/logger"
)

func RegisterProtocol() bool {
	if runtime.GOOS != "windows" {
		logger.LogError("Protocol registration only supported on Windows")
		return false
	}

	exePath, err := os.Executable()
	if err != nil {
		logger.LogError("Error getting executable path: %v", err)
		return false
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		logger.LogError("Error getting absolute path: %v", err)
		return false
	}

	logger.LogInfo("Executable path: %s", exePath)

	commands := [][]string{
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser`, "/ve", "/d", "URL:Facebook Auto Login Protocol", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser`, "/v", "URL Protocol", "/d", "", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser\shell\open\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%1"`, exePath), "/f"},
	}

	logger.LogInfo("Executing protocol registration commands...")

	for i, cmd := range commands {
		logger.LogInfo("Command %d: %v", i+1, cmd)
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			logger.LogError("Error executing command %v: %v", cmd, err)

			if strings.Contains(err.Error(), "Access is denied") || strings.Contains(err.Error(), "access denied") {
				logger.LogError("Access denied. Run as administrator")
			}
			return false
		}
		logger.LogSuccess("Command %d executed successfully", i+1)
	}

	logger.LogSuccess("Protocol registered successfully")
	return true
}

func HandleProtocolLaunch(protocolURL string) {
	logger.LogInfo("Started via protocol: %s", protocolURL)

	u, err := url.Parse(protocolURL)
	if err != nil {
		logger.LogError("URL parsing error: %v", err)
		return
	}

	path := u.Path
	query := u.Query()

	logger.LogInfo("Protocol parsing - path: %s, params: %v", path, query)

	if path == "/test" {
		logger.LogSuccess("Protocol test successful")
		logger.LogInfo("Application ready to work")
		return
	}

	if path == "" || path == "/" {
		if query.Get("login") != "" || query.Get("password") != "" {
			path = "/browser"
			logger.LogInfo("Determined operation type: %s", path)
		}
	}

	if path == "/browser" {
		username := query.Get("login")
		password := query.Get("password")

		if username == "" || password == "" {
			logger.LogError("Login or password not specified in protocol")
			logger.LogInfo("Use: siteparser://browser/?login=...&password=...")
			return
		}

		logger.LogInfo("Facebook automation via protocol")
		logger.LogInfo("Login: %s", username)

		result := browser.OpenBrowserWithLogin(username, password)

		if err := fileutils.SaveBrowserResultToFile(result); err != nil {
			logger.LogError("Error saving result: %v", err)
		}

		if result.Success {
			logger.LogSuccess("Facebook automation successful")
			logger.LogInfo("📝 Message: %s", result.Message)
		} else {
			logger.LogError("Automation error: %s", result.Error)
		}
	} else {
		logger.LogError("Invalid protocol format: %s", path)
		logger.LogInfo("Use:")
		logger.LogInfo("  siteparser://browser/?login=...&password=...")
		logger.LogInfo("  siteparser://test")
	}
}
