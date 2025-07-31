package cli

import (
	"flag"

	"rpa-dfs-engine/internal/browser"
	"rpa-dfs-engine/internal/fileutils"
	"rpa-dfs-engine/internal/logger"
)

func HandleCommandLine() {
	var (
		username = flag.String("l", "", "Facebook login (email)")
		password = flag.String("p", "", "Facebook password")
		help     = flag.Bool("h", false, "Show help")
	)

	flag.Parse()

	if *help {
		ShowHelp()
		return
	}

	if *username == "" || *password == "" {
		logger.LogError("Login and password required")
		logger.LogInfo("💡 Use: facebook-login.exe -l LOGIN -p PASSWORD")
		ShowHelp()
		return
	}

	logger.LogInfo("Starting Facebook automation")
	logger.LogInfo("Login: %s", *username)

	result := browser.OpenBrowserWithLogin(*username, *password)

	if err := fileutils.SaveBrowserResultToFile(result); err != nil {
		logger.LogError("Error saving to file: %v", err)
	}

	if result.Success {
		logger.LogSuccess("Facebook automation successful")
		logger.LogInfo("🌐 URL: %s", result.URL)
		logger.LogInfo("📝 Message: %s", result.Message)
	} else {
		logger.LogError("Automation error: %s", result.Error)
	}
}

func ShowHelp() {
	logger.LogInfo("🌐 Facebook Auto Login - Automatic Facebook Login")
	logger.LogInfo("==================================================")
	logger.LogInfo("")
	logger.LogInfo("Usage:")
	logger.LogInfo("  facebook-login.exe -l LOGIN -p PASSWORD")
	logger.LogInfo("")
	logger.LogInfo("Parameters:")
	logger.LogInfo("  -l LOGIN     Facebook login (email)")
	logger.LogInfo("  -p PASSWORD  Facebook password")
	logger.LogInfo("  -h           Show this help")
	logger.LogInfo("")
	logger.LogInfo("Examples:")
	logger.LogInfo("  facebook-login.exe -l user@example.com -p mypassword")
	logger.LogInfo("  facebook-login.exe -l john.doe@gmail.com -p secret123")
	logger.LogInfo("")
	logger.LogInfo("Result:")
	logger.LogInfo("  - Chrome browser will open")
	logger.LogInfo("  - Automatically navigate to Facebook")
	logger.LogInfo("  - Enter login and password")
	logger.LogInfo("  - Click login button")
	logger.LogInfo("  - Result saved to facebook_result_*.txt file")
	logger.LogInfo("")
	logger.LogInfo("Requirements:")
	logger.LogInfo("  - Google Chrome installed")
	logger.LogInfo("  - Internet connection")
	logger.LogInfo("  - Valid Facebook credentials")
}
