package protocol

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"rpa-dfs-engine/internal/config"
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
		{"reg", "add", `HKEY_CLASSES_ROOT\` + config.PROTOCOL_NAME, "/ve", "/d", "URL:Scrooge PRotocol", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\` + config.PROTOCOL_NAME, "/v", "URL Protocol", "/d", "", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\` + config.PROTOCOL_NAME + `\shell\open\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%1"`, exePath), "/f"},
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
