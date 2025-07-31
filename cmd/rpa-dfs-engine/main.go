package main

import (
	"fmt"
	"os"
	"strings"

	"rpa-dfs-engine/internal/cli"
	"rpa-dfs-engine/internal/html"
	"rpa-dfs-engine/internal/logger"
	"rpa-dfs-engine/internal/protocol"
)

func main() {
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("❌ Logger init error: %v\n", err)
	}
	defer logger.CloseLogger()

	logger.LogInfo("Facebook Auto Login started")

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "siteparser://") {
		logger.LogInfo("Started via protocol: %s", os.Args[1])
		protocol.HandleProtocolLaunch(os.Args[1])
		return
	}

	handleFirstRun()

	if len(os.Args) > 1 {
		logger.LogInfo("Processing command line arguments")
		cli.HandleCommandLine()
	} else {
		logger.LogInfo("Showing help")
		cli.ShowHelp()
	}

	logger.LogInfo("Application finished")
}

func handleFirstRun() {
	logger.LogInfo("Checking protocol registration")
	html.CreateHtmlInterface()

	if !protocol.IsProtocolRegistered() {
		logger.LogInfo("Protocol not registered, starting registration")
		fmt.Println("=== Facebook Auto Login - First Run ===")
		fmt.Println("Registering siteparser:// protocol...")

		if protocol.RegisterProtocol() {
			logger.LogSuccess("Protocol registered successfully")
			fmt.Println("✅ Protocol registered successfully!")
			logger.LogSuccess("HTML interface created")
			fmt.Println("✅ HTML interface created!")
			fmt.Println("\nNow you can:")
			fmt.Println("1. Open 'facebook-login.html' in browser")
			fmt.Println("2. Use commands: facebook-login.exe -l LOGIN -p PASSWORD")
			fmt.Println("3. Use protocol: siteparser://browser/?login=...&password=...")
		} else {
			logger.LogWarning("Failed to register protocol")
			fmt.Println("⚠️ Failed to register protocol.")
			fmt.Println("💡 To register protocol run: run-as-admin.bat")
			fmt.Println("📄 Application will continue without protocol")
			logger.LogInfo("HTML interface already created")
			fmt.Println("✅ HTML interface created!")
		}
	} else {
		logger.LogInfo("Protocol already registered")
	}
}
