package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func saveBrowserResultToFile(result BrowserResult) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %v", err)
	}

	exeDir := filepath.Dir(exePath)

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("facebook_result_%s.txt", timestamp)
	filePath := filepath.Join(exeDir, filename)

	content := fmt.Sprintf(`=== Facebook Auto Login Result ===
Time: %s
Success: %t
URL: %s
Login: %s
Message: %s
`,
		time.Unix(result.Timestamp, 0).Format("2006-01-02 15:04:05"),
		result.Success,
		result.URL,
		result.Username,
		result.Message,
	)

	if result.Error != "" {
		content += fmt.Sprintf("Error: %s\n", result.Error)
	}

	content += "\n=== END OF RESULT ===\n"

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Printf("✅ Result saved to file: %s\n", filename)
	return nil
}
