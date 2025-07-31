package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if err := initLogger(); err != nil {
		fmt.Printf("❌ Logger init error: %v\n", err)
	}
	defer closeLogger()

	LogInfo("Facebook Auto Login started")

	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "siteparser://") {
		LogInfo("Started via protocol: %s", os.Args[1])
		handleProtocolLaunch(os.Args[1])
		return
	}

	handleFirstRun()

	if len(os.Args) > 1 {
		LogInfo("Processing command line arguments")
		handleCommandLine()
	} else {
		LogInfo("Showing help")
		showHelp()
	}

	LogInfo("Application finished")
}

func handleFirstRun() {
	LogInfo("Checking protocol registration")
	createHtmlInterface()

	if !isProtocolRegistered() {
		LogInfo("Protocol not registered, starting registration")
		fmt.Println("=== Facebook Auto Login - First Run ===")
		fmt.Println("Registering siteparser:// protocol...")

		if registerProtocol() {
			LogSuccess("Protocol registered successfully")
			fmt.Println("✅ Protocol registered successfully!")
			LogSuccess("HTML interface created")
			fmt.Println("✅ HTML interface created!")
			fmt.Println("\nNow you can:")
			fmt.Println("1. Open 'facebook-login.html' in browser")
			fmt.Println("2. Use commands: facebook-login.exe -l LOGIN -p PASSWORD")
			fmt.Println("3. Use protocol: siteparser://browser/?login=...&password=...")
		} else {
			LogWarning("Failed to register protocol")
			fmt.Println("⚠️ Failed to register protocol.")
			fmt.Println("💡 To register protocol run: run-as-admin.bat")
			fmt.Println("📄 Application will continue without protocol")
			LogInfo("HTML interface already created")
			fmt.Println("✅ HTML interface created!")
		}
	} else {
		LogInfo("Protocol already registered")
	}
}
