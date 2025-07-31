package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func isProtocolRegistered() bool {
	if runtime.GOOS != "windows" {
		return false
	}

	cmd := exec.Command("reg", "query", `HKEY_CLASSES_ROOT\siteparser`, "/ve")
	return cmd.Run() == nil
}

func registerProtocol() bool {
	if runtime.GOOS != "windows" {
		LogError("Protocol registration only supported on Windows")
		fmt.Println("Protocol registration only supported on Windows")
		return false
	}

	exePath, err := os.Executable()
	if err != nil {
		LogError("Error getting executable path: %v", err)
		fmt.Println("Error getting executable path:", err)
		return false
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		LogError("Error getting absolute path: %v", err)
		fmt.Println("Error getting absolute path:", err)
		return false
	}

	LogInfo("Executable path: %s", exePath)

	commands := [][]string{
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser`, "/ve", "/d", "URL:Facebook Auto Login Protocol", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser`, "/v", "URL Protocol", "/d", "", "/f"},
		{"reg", "add", `HKEY_CLASSES_ROOT\siteparser\shell\open\command`, "/ve", "/d", fmt.Sprintf(`"%s" "%%1"`, exePath), "/f"},
	}

	LogInfo("Executing protocol registration commands...")

	for i, cmd := range commands {
		LogInfo("Command %d: %v", i+1, cmd)
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			LogError("Error executing command %v: %v", cmd, err)
			fmt.Printf("Error executing command %v: %v\n", cmd, err)

			if strings.Contains(err.Error(), "Access is denied") || strings.Contains(err.Error(), "access denied") {
				LogError("Access denied. Run as administrator")
				fmt.Println("❌ Access denied. Run as administrator")
			}
			return false
		}
		LogSuccess("Command %d executed successfully", i+1)
	}

	LogSuccess("Protocol registered successfully")
	return true
}

func handleProtocolLaunch(protocolURL string) {
	fmt.Println("Started via protocol:", protocolURL)

	u, err := url.Parse(protocolURL)
	if err != nil {
		LogError("URL parsing error: %v", err)
		fmt.Println("❌ URL parsing error:", err)
		return
	}

	path := u.Path
	query := u.Query()

	LogInfo("Protocol parsing - path: %s, params: %v", path, query)

	if path == "/test" {
		fmt.Println("✅ Protocol test successful!")
		fmt.Println("Application ready to work.")
		return
	}

	if path == "" || path == "/" {
		if query.Get("login") != "" || query.Get("password") != "" {
			path = "/browser"
			LogInfo("Determined operation type: %s", path)
		}
	}

	if path == "/browser" {
		username := query.Get("login")
		password := query.Get("password")

		if username == "" || password == "" {
			LogError("Login or password not specified in protocol")
			fmt.Println("❌ Login or password not specified. Use: siteparser://browser/?login=...&password=...")
			return
		}

		LogInfo("Facebook automation via protocol")
		LogInfo("Login: %s", username)
		fmt.Println("🌐 Facebook automation")
		fmt.Printf("👤 Login: %s\n", username)

		result := openBrowserWithLogin(username, password)

		if err := saveBrowserResultToFile(result); err != nil {
			LogError("Error saving result: %v", err)
			fmt.Printf("❌ Error saving to file: %v\n", err)
		}

		if result.Success {
			LogSuccess("Facebook automation successful")
			fmt.Printf("✅ Facebook automation successful!\n")
			fmt.Printf("📝 Message: %s\n", result.Message)
		} else {
			LogError("Automation error: %s", result.Error)
			fmt.Printf("❌ Error: %s\n", result.Error)
		}
	} else {
		LogError("Invalid protocol format: %s", path)
		fmt.Println("❌ Invalid protocol format. Use:")
		fmt.Println("  siteparser://browser/?login=...&password=...")
		fmt.Println("  siteparser://test")
	}
}
