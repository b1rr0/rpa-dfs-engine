package html

import (
	"os"
	"path/filepath"
	"rpa-dfs-engine/internal/config"
	"rpa-dfs-engine/internal/logger"
	"strings"
)

func CreateHtmlInterface() {
	exePath, err := os.Executable()
	logger.LogInfo("Checking protocol registration")

	logger.LogInfo("Facebook Auto Login started")
	if err != nil {
		logger.LogError("Error getting executable path: %v", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	if filepath.Base(exeDir) == "logs" {
		exeDir = filepath.Dir(exeDir)
	}
	templatePath := filepath.Join(filepath.Dir(exeDir), "web", "template.html")
	htmlPath := filepath.Join(exeDir, "facebook-login.html")

	htmlContent, err := os.ReadFile(templatePath)
	logger.LogInfo("htmlContent: %v", htmlContent)
	htmlContent = []byte(strings.ReplaceAll(string(htmlContent), "#PROTOCOL_NAME#", config.PROTOCOL_NAME))
	logger.LogInfo("htmlContent: %v", htmlContent)

	if err != nil {
		logger.LogError("Error reading template file: %v", err)
		return
	}

	err = os.WriteFile(htmlPath, htmlContent, 0644)
	if err != nil {
		logger.LogError("Error creating HTML file: %v", err)
		return
	}

	logger.LogSuccess("HTML interface created: %s", htmlPath)
}
