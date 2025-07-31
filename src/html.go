package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createHtmlInterface() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	exeDir := filepath.Dir(exePath)
	templatePath := filepath.Join(filepath.Dir(exeDir), "src", "template.html")
	htmlPath := filepath.Join(exeDir, "facebook-login.html")

	htmlContent, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Println("Error reading template file:", err)
		return
	}

	err = os.WriteFile(htmlPath, htmlContent, 0644)
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}

	fmt.Printf("✅ HTML interface created: %s\n", htmlPath)
}
