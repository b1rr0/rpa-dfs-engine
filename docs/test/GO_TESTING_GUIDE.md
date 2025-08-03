# Go Testing Guide for RPA DFS Engine

> Practical Go testing techniques and patterns for desktop RPA automation

## 🏃 **Quick Start with Go Testing**

### **Basic Test Structure**

```go
package browser_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "rpa-dfs-engine/internal/browser"
)

func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {
    // Arrange
    input := "test data"
    
    // Act
    result := browser.ProcessInput(input)
    
    // Assert
    assert.Equal(t, "expected", result)
}
```

### **Running Tests**

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/browser

# Run specific test function
go test -run TestOpenBrowserWithLogin ./internal/browser

# Run tests and generate coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🛠️ **Essential Testing Packages**

### **1. Standard `testing` Package**
```go
import "testing"

func TestBasicFunction(t *testing.T) {
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### **2. Testify - Enhanced Assertions**
```bash
go get github.com/stretchr/testify
```

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)

func TestWithTestify(t *testing.T) {
    assert.Equal(t, expected, actual)
    assert.True(t, condition)
    assert.NoError(t, err)
    assert.Contains(t, "hello world", "world")
    
    // require stops test execution on failure
    require.NoError(t, err)
    // This line won't run if err != nil
}
```

## 🎭 **Mocking Patterns**

### **1. Interface-Based Mocking**

#### **Define Interfaces for Dependencies**
```go
// internal/browser/interfaces.go
package browser

type ChromeDetector interface {
    IsInstalled() bool
    GetVersion() string
}

type FileWriter interface {
    WriteFile(filename string, data []byte) error
}
```

#### **Create Mock Implementations**
```go
// internal/browser/browser_test.go
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
```

#### **Use Mocks in Tests**
```go
func TestBrowser_WhenChromeNotInstalled_ReturnsError(t *testing.T) {
    // Arrange
    mockDetector := &MockChromeDetector{
        IsInstalledResult: false,
    }
    browser := NewBrowser(mockDetector)
    
    // Act
    result := browser.OpenWithLogin("user", "pass")
    
    // Assert
    assert.False(t, result.Success)
    assert.Contains(t, result.Error, "Chrome not found")
}
```

### **2. Testify Mock Framework**

```go
type MockFileWriter struct {
    mock.Mock
}

func (m *MockFileWriter) WriteFile(filename string, data []byte) error {
    args := m.Called(filename, data)
    return args.Error(0)
}

func TestSaveResult_CallsFileWriter(t *testing.T) {
    // Arrange
    mockWriter := new(MockFileWriter)
    mockWriter.On("WriteFile", "test.json", mock.Anything).Return(nil)
    
    saver := NewResultSaver(mockWriter)
    
    // Act
    err := saver.SaveResult("test.json", []byte("data"))
    
    // Assert
    assert.NoError(t, err)
    mockWriter.AssertExpectations(t) // Verify mock was called as expected
}
```

## 📁 **Test File Organization**

### **1. Test Files Alongside Source**
```
internal/browser/
├── browser.go          # Source code
├── browser_test.go     # Tests for browser.go
├── chrome.go           # Source code  
├── chrome_test.go      # Tests for chrome.go
└── mocks.go           # Mock implementations
```

### **2. Package Naming Convention**
```go
// Option 1: Same package (white-box testing)
package browser

import "testing"

func TestInternalFunction(t *testing.T) {
    // Can access private functions and variables
}

// Option 2: External package (black-box testing)  
package browser_test

import (
    "testing"
    "rpa-dfs-engine/internal/browser"
)

func TestPublicAPI(t *testing.T) {
    // Can only access exported functions
    result := browser.OpenBrowserWithLogin("user", "pass")
}
```

## 🔧 **Testing Patterns for RPA Components**

### **1. Testing Browser Automation**

```go
func TestBrowserAutomation_LoginFlow_Success(t *testing.T) {
    // Arrange
    mockBrowser := &MockBrowser{
        Elements: map[string]bool{
            "#email":    true,
            "#pass":     true,
            "#loginbutton": true,
        },
    }
    
    automation := NewBrowserAutomation(mockBrowser)
    
    // Act
    result := automation.PerformLogin("user@test.com", "password")
    
    // Assert
    assert.True(t, result.Success)
    assert.Equal(t, "login successful", result.Message)
}

func TestBrowserAutomation_ElementNotFound_ReturnsError(t *testing.T) {
    // Arrange
    mockBrowser := &MockBrowser{
        Elements: map[string]bool{
            "#email": false, // Element not found
        },
    }
    
    automation := NewBrowserAutomation(mockBrowser)
    
    // Act
    result := automation.PerformLogin("user@test.com", "password")
    
    // Assert
    assert.False(t, result.Success)
    assert.Contains(t, result.Error, "element not found: #email")
}
```

### **2. Testing Configuration Loading**

```go
func TestConfigLoader_ValidYAML_LoadsCorrectly(t *testing.T) {
    // Arrange
    yamlContent := `
browser_timeout: 30s
facebook_url: "https://facebook.com"
`
    mockFileReader := &MockFileReader{
        Content: []byte(yamlContent),
    }
    
    loader := NewConfigLoader(mockFileReader)
    
    // Act
    config, err := loader.LoadConfig("config.yaml")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 30*time.Second, config.BrowserTimeout)
    assert.Equal(t, "https://facebook.com", config.FacebookURL)
}

func TestConfigLoader_InvalidYAML_ReturnsError(t *testing.T) {
    // Arrange
    invalidYAML := `
browser_timeout: [invalid yaml
`
    mockFileReader := &MockFileReader{
        Content: []byte(invalidYAML),
    }
    
    loader := NewConfigLoader(mockFileReader)
    
    // Act
    _, err := loader.LoadConfig("config.yaml")
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "yaml")
}
```

### **3. Testing Protocol Parsing**

```go
func TestProtocolParser_ValidURL_ParsesParameters(t *testing.T) {
    tests := []struct {
        name     string
        url      string
        expected ProtocolParams
    }{
        {
            name: "browser with login and password",
            url:  "siteparser://browser/?login=user&password=pass",
            expected: ProtocolParams{
                Action:   "browser",
                Login:    "user",
                Password: "pass",
            },
        },
        {
            name: "test mode",
            url:  "siteparser://test",
            expected: ProtocolParams{
                Action: "test",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            result, err := ParseProtocolURL(tt.url)
            
            // Assert
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### **4. Testing Error Handling**

```go
func TestFileUtils_SaveResult_HandlesErrors(t *testing.T) {
    tests := []struct {
        name          string
        mockError     error
        expectedError string
    }{
        {
            name:          "permission denied",
            mockError:     os.ErrPermission,
            expectedError: "permission denied",
        },
        {
            name:          "file not found",
            mockError:     os.ErrNotExist,
            expectedError: "file not found",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockWriter := &MockFileWriter{
                Error: tt.mockError,
            }
            
            fileUtils := NewFileUtils(mockWriter)
            
            // Act
            err := fileUtils.SaveResult("test.json", []byte("data"))
            
            // Assert
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedError)
        })
    }
}
```

## 🧪 **Test Helpers and Utilities**

### **1. Test Data Setup**

```go
// testhelpers/fixtures.go
package testhelpers

func CreateTestBrowserResult() types.BrowserResult {
    return types.BrowserResult{
        Success:   true,
        URL:       "https://facebook.com/feed",
        Username:  "test@example.com",
        Message:   "Login successful",
        Timestamp: time.Now().Unix(),
    }
}

func CreateTestConfig() *config.Config {
    return &config.Config{
        BrowserTimeout: 30 * time.Second,
        FacebookURL:    "https://facebook.com",
        ProtocolName:   "siteparser",
    }
}
```

### **2. Common Test Assertions**

```go
// testhelpers/assertions.go
package testhelpers

func AssertBrowserResultSuccess(t *testing.T, result types.BrowserResult) {
    t.Helper() // This function is a test helper
    
    assert.True(t, result.Success, "Expected successful browser result")
    assert.NotEmpty(t, result.URL, "Expected non-empty URL")
    assert.NotEmpty(t, result.Username, "Expected non-empty username")
    assert.Empty(t, result.Error, "Expected no error message")
}

func AssertBrowserResultError(t *testing.T, result types.BrowserResult, expectedError string) {
    t.Helper()
    
    assert.False(t, result.Success, "Expected failed browser result")
    assert.NotEmpty(t, result.Error, "Expected error message")
    assert.Contains(t, result.Error, expectedError, "Error message should contain expected text")
}
```

## 📊 **Test Coverage and Reporting**

### **1. Coverage Commands**

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage by package
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage for specific package
go test -cover ./internal/browser
```

### **2. Coverage Targets**

```bash
# Aim for these coverage levels:
# - Critical components (browser, protocol): 90%+
# - Business logic (handlers, cli): 80%+  
# - Utilities (logger, fileutils): 70%+
# - Configuration: 60%+
```

## 🚀 **Performance Testing**

### **1. Benchmark Tests**

```go
func BenchmarkProtocolParsing(b *testing.B) {
    url := "siteparser://browser/?login=user&password=pass"
    
    for i := 0; i < b.N; i++ {
        _, _ = ParseProtocolURL(url)
    }
}

func BenchmarkConfigLoading(b *testing.B) {
    configContent := []byte(`browser_timeout: 30s`)
    
    b.ResetTimer() // Don't count setup time
    
    for i := 0; i < b.N; i++ {
        _, _ = ParseConfig(configContent)
    }
}
```

### **2. Running Benchmarks**

```bash
# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkProtocolParsing ./internal/protocol

# Run benchmarks with memory stats
go test -bench=. -benchmem ./...
```

## 📋 **Testing Checklist**

- [ ] **Test files named with `_test.go` suffix**
- [ ] **Test functions start with `Test`**
- [ ] **One logical assertion per test**
- [ ] **Descriptive test names (Subject_Scenario_Result)**
- [ ] **Mock all external dependencies**
- [ ] **Tests can run independently**
- [ ] **No hardcoded file paths or URLs**
- [ ] **Error cases are tested**
- [ ] **Edge cases are covered**
- [ ] **Tests are fast (< 1 second each)**

---

> 💡 **Go Testing Best Practice**: Keep tests simple, focused, and fast. Use table-driven tests for multiple scenarios and always test error conditions alongside happy paths. 