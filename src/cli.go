package main

import (
	"flag"
	"fmt"
)

func handleCommandLine() {
	var (
		username = flag.String("l", "", "Facebook login (email)")
		password = flag.String("p", "", "Facebook password")
		help     = flag.Bool("h", false, "Show help")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *username == "" || *password == "" {
		LogError("Login and password required")
		fmt.Println("❌ Login and password required")
		fmt.Println("💡 Use: facebook-login.exe -l LOGIN -p PASSWORD")
		showHelp()
		return
	}

	LogInfo("Starting Facebook automation")
	LogInfo("Login: %s", *username)

	result := openBrowserWithLogin(*username, *password)

	if err := saveBrowserResultToFile(result); err != nil {
		LogError("Error saving to file: %v", err)
		fmt.Printf("❌ Error saving to file: %v\n", err)
	}

	if result.Success {
		LogSuccess("Facebook automation successful")
		fmt.Printf("✅ Facebook automation successful!\n")
		fmt.Printf("🌐 URL: %s\n", result.URL)
		fmt.Printf("📝 Message: %s\n", result.Message)
	} else {
		LogError("Automation error: %s", result.Error)
		fmt.Printf("❌ Error: %s\n", result.Error)
	}
}

func showHelp() {
	fmt.Println("🌐 Facebook Auto Login - Automatic Facebook Login")
	fmt.Println("==================================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  facebook-login.exe -l LOGIN -p PASSWORD")
	fmt.Println()
	fmt.Println("Parameters:")
	fmt.Println("  -l LOGIN     Facebook login (email)")
	fmt.Println("  -p PASSWORD  Facebook password")
	fmt.Println("  -h           Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  facebook-login.exe -l user@example.com -p mypassword")
	fmt.Println("  facebook-login.exe -l john.doe@gmail.com -p secret123")
	fmt.Println()
	fmt.Println("Result:")
	fmt.Println("  - Chrome browser will open")
	fmt.Println("  - Automatically navigate to Facebook")
	fmt.Println("  - Enter login and password")
	fmt.Println("  - Click login button")
	fmt.Println("  - Result saved to facebook_result_*.txt file")
	fmt.Println()
	fmt.Println("Requirements:")
	fmt.Println("  - Google Chrome installed")
	fmt.Println("  - Internet connection")
	fmt.Println("  - Valid Facebook credentials")
}
