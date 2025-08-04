package mocks

import (
	"fmt"
	"strings"
)

// MockTraverserBrowser implements the traverser.Browser interface for testing
type MockTraverserBrowser struct {
	NavigationHistory []string
	FormData          map[string]string // selector -> value
	ClickedElements   []string
	UploadedFiles     map[string]string // selector -> filePath
	ShouldFailOn      map[string]bool   // operation -> should fail
	ErrorMessages     map[string]string // operation -> error message
	IsClosed          bool
}

// NewMockTraverserBrowser creates a new mock browser for testing
func NewMockTraverserBrowser() *MockTraverserBrowser {
	return &MockTraverserBrowser{
		NavigationHistory: make([]string, 0),
		FormData:          make(map[string]string),
		ClickedElements:   make([]string, 0),
		UploadedFiles:     make(map[string]string),
		ShouldFailOn:      make(map[string]bool),
		ErrorMessages:     make(map[string]string),
		IsClosed:          false,
	}
}

// NavigateTo implements Browser interface
func (m *MockTraverserBrowser) NavigateTo(url string) error {
	if m.ShouldFailOn["navigate"] {
		return fmt.Errorf(m.ErrorMessages["navigate"])
	}
	m.NavigationHistory = append(m.NavigationHistory, url)
	return nil
}

// FillField implements Browser interface
func (m *MockTraverserBrowser) FillField(selector, value string) error {
	if m.ShouldFailOn["fillField"] {
		return fmt.Errorf(m.ErrorMessages["fillField"])
	}
	m.FormData[selector] = value
	return nil
}

// ClickElement implements Browser interface
func (m *MockTraverserBrowser) ClickElement(selector string) error {
	if m.ShouldFailOn["clickElement"] {
		return fmt.Errorf(m.ErrorMessages["clickElement"])
	}
	m.ClickedElements = append(m.ClickedElements, selector)
	return nil
}

// SendFile implements Browser interface
func (m *MockTraverserBrowser) SendFile(selector, filePath string) error {
	if m.ShouldFailOn["sendFile"] {
		return fmt.Errorf(m.ErrorMessages["sendFile"])
	}
	m.UploadedFiles[selector] = filePath
	return nil
}

// Close implements Browser interface
func (m *MockTraverserBrowser) Close() error {
	if m.ShouldFailOn["close"] {
		return fmt.Errorf(m.ErrorMessages["close"])
	}
	m.IsClosed = true
	return nil
}

// Test helper methods

// SetFailureFor configures the mock to fail for specific operations
func (m *MockTraverserBrowser) SetFailureFor(operation, errorMessage string) {
	m.ShouldFailOn[operation] = true
	m.ErrorMessages[operation] = errorMessage
}

// ClearFailures removes all configured failures
func (m *MockTraverserBrowser) ClearFailures() {
	m.ShouldFailOn = make(map[string]bool)
	m.ErrorMessages = make(map[string]string)
}

// Reset clears all recorded actions and failures
func (m *MockTraverserBrowser) Reset() {
	m.NavigationHistory = make([]string, 0)
	m.FormData = make(map[string]string)
	m.ClickedElements = make([]string, 0)
	m.UploadedFiles = make(map[string]string)
	m.ClearFailures()
	m.IsClosed = false
}

// GetLastNavigation returns the last URL navigated to
func (m *MockTraverserBrowser) GetLastNavigation() string {
	if len(m.NavigationHistory) == 0 {
		return ""
	}
	return m.NavigationHistory[len(m.NavigationHistory)-1]
}

// GetFieldValue returns the value filled for a specific selector
func (m *MockTraverserBrowser) GetFieldValue(selector string) string {
	return m.FormData[selector]
}

// WasElementClicked checks if an element was clicked
func (m *MockTraverserBrowser) WasElementClicked(selector string) bool {
	for _, clicked := range m.ClickedElements {
		if clicked == selector {
			return true
		}
	}
	return false
}

// GetUploadedFile returns the file path uploaded to a selector
func (m *MockTraverserBrowser) GetUploadedFile(selector string) string {
	return m.UploadedFiles[selector]
}

// GetNavigationCount returns the number of navigations performed
func (m *MockTraverserBrowser) GetNavigationCount() int {
	return len(m.NavigationHistory)
}

// GetFormFieldCount returns the number of form fields filled
func (m *MockTraverserBrowser) GetFormFieldCount() int {
	return len(m.FormData)
}

// GetClickCount returns the number of elements clicked
func (m *MockTraverserBrowser) GetClickCount() int {
	return len(m.ClickedElements)
}

// HasNavigatedTo checks if the browser navigated to a specific URL
func (m *MockTraverserBrowser) HasNavigatedTo(url string) bool {
	for _, nav := range m.NavigationHistory {
		if nav == url {
			return true
		}
	}
	return false
}

// HasNavigatedToPattern checks if the browser navigated to a URL matching a pattern
func (m *MockTraverserBrowser) HasNavigatedToPattern(pattern string) bool {
	for _, nav := range m.NavigationHistory {
		if strings.Contains(nav, pattern) {
			return true
		}
	}
	return false
}
