package mocks

import (
	"fmt"
	"rpa-dfs-engine/internal/types"
)

type MockChromeDetector struct {
	IsInstalledResult bool
	VersionResult     string
}

func (m *MockChromeDetector) IsInstalled() bool {
	return m.IsInstalledResult
}

func (m *MockChromeDetector) GetVersion() string {
	return m.VersionResult
}

type MockBrowserAutomation struct {
	Elements        map[string]bool // Element selector -> exists
	ShouldFailLogin bool
	FailureError    string
	LoginResult     types.BrowserResult
	NavigationDelay bool
}

func NewMockBrowserAutomation() *MockBrowserAutomation {
	return &MockBrowserAutomation{
		Elements: make(map[string]bool),
	}
}

func (m *MockBrowserAutomation) PerformLogin(username, password string) types.BrowserResult {
	if m.ShouldFailLogin {
		return types.BrowserResult{
			Success:  false,
			URL:      "https://facebook.com",
			Username: username,
			Error:    m.FailureError,
		}
	}

	if m.LoginResult.Username == "" {
		m.LoginResult.Username = username
	}

	return m.LoginResult
}

func (m *MockBrowserAutomation) ElementExists(selector string) bool {
	exists, found := m.Elements[selector]
	return found && exists
}

func (m *MockBrowserAutomation) ClickElement(selector string) error {
	if !m.ElementExists(selector) {
		return fmt.Errorf("element not found: %s", selector)
	}
	return nil
}

func (m *MockBrowserAutomation) SendKeys(selector, text string) error {
	if !m.ElementExists(selector) {
		return fmt.Errorf("element not found: %s", selector)
	}
	return nil
}

func (m *MockBrowserAutomation) SetElement(selector string, exists bool) {
	m.Elements[selector] = exists
}

func (m *MockBrowserAutomation) SetFailure(shouldFail bool, errorMessage string) {
	m.ShouldFailLogin = shouldFail
	m.FailureError = errorMessage
}

func (m *MockBrowserAutomation) Reset() {
	m.Elements = make(map[string]bool)
	m.ShouldFailLogin = false
	m.FailureError = ""
	m.LoginResult = types.BrowserResult{}
	m.NavigationDelay = false
}
