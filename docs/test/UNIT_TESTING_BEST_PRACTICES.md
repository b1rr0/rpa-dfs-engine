# Unit Testing Best Practices for RPA DFS Engine

> Detailed principles for writing excellent unit tests, based on [Steve Sanderson's guide](https://gist.github.com/vadymhimself/763e96dd8495bb77325efd082e63c9f5)

## 🎯 **Core Philosophy: TDD as a Design Process**

### **Unit Testing vs. Other Testing Types**

**Unit testing is NOT about finding bugs.** It's about designing robust software components through Test-Driven Development (TDD).

| Testing Type | Purpose | When to Use | Tools for RPA DFS Engine |
|--------------|---------|-------------|--------------------------|
| **Unit Tests** | Design software components | During development (TDD) | Go's `testing` package |
| **Integration Tests** | Detect regressions in complete workflows | CI/CD pipeline | Browser automation with real Chrome |
| **Manual Testing** | Find bugs and usability issues | Before releases | Actual Facebook login testing |

### **The Good vs. Bad Unit Test Scale**

```
[Sweet Spot A] ←────────────────────────→ [Sweet Spot B]
Pure Unit Tests                      Pure Integration Tests
(Test single components)             (Test complete workflows)
        ↓                                        ↓
✅ Fast, Focused                          ✅ Proves real functionality
✅ Test design decisions                  ✅ Catches integration issues  
✅ Easy to maintain                       ✅ Stable across refactoring
✅ Cheap to run                          ✅ End-user focused

                    [Danger Zone]
               Hybrid/Fragile Tests
                        ↓
               ❌ Slow and brittle
               ❌ Hard to maintain
               ❌ Unclear purpose
               ❌ Break on refactoring
```

## 📏 **The Six Golden Rules**

### **1. Make Each Test Orthogonal (Independent)**

**Rule:** Any given behavior should be specified in one and only one test.

#### ✅ **Good Example:**
```go
// Test ONLY the username validation logic
func TestValidateCredentials_WithEmptyUsername_ReturnsUsernameError(t *testing.T) {
    result := ValidateCredentials("", "password123")
    
    assert.False(t, result.IsValid)
    assert.Contains(t, result.Errors, "username is required")
    // DON'T test password validation here
}

// Test ONLY the password validation logic  
func TestValidateCredentials_WithEmptyPassword_ReturnsPasswordError(t *testing.T) {
    result := ValidateCredentials("user@email.com", "")
    
    assert.False(t, result.IsValid)
    assert.Contains(t, result.Errors, "password is required")
    // DON'T test username validation here
}
```

#### ❌ **Bad Example:**
```go
// Testing multiple behaviors in one test
func TestValidateCredentials_WithBadInputs_ReturnsErrors(t *testing.T) {
    result := ValidateCredentials("", "")
    
    assert.False(t, result.IsValid)
    assert.Contains(t, result.Errors, "username is required") // Testing username logic
    assert.Contains(t, result.Errors, "password is required") // Testing password logic
    // This test will break if EITHER username OR password validation changes
}
```

### **2. Test Only One Code Unit at a Time**

**Rule:** Test individual packages/classes in isolation using mocks for dependencies.

#### ✅ **Good Example:**
```go
// Testing browser logic in isolation
func TestOpenBrowserWithLogin_WhenChromeInstalled_ReturnsSuccess(t *testing.T) {
    // Mock the Chrome detection dependency
    mockDetector := &MockChromeDetector{
        IsInstalled: true,
    }
    
    browser := NewBrowser(mockDetector)
    result := browser.OpenWithLogin("user", "pass")
    
    assert.True(t, result.Success)
}
```

#### ❌ **Bad Example:**
```go
// Testing browser + file utils + config all together
func TestBrowserWorkflow_Complete_CreatesFileAndLogsSuccess(t *testing.T) {
    // This test depends on:
    // - Chrome being installed
    // - File system being writable  
    // - Config file existing
    // - Logger working correctly
    // If ANY of these change, this test breaks
    result := CompleteWorkflow("user", "pass")
    assert.True(t, result.Success)
}
```

### **3. Mock Out All External Services and State**

**Rule:** Don't depend on Chrome, filesystem, network, or any external state.

#### ✅ **Good Example:**
```go
type MockFileWriter struct {
    WrittenFiles map[string][]byte
    ShouldFail   bool
}

func (m *MockFileWriter) WriteFile(path string, data []byte) error {
    if m.ShouldFail {
        return errors.New("mock write failure")
    }
    m.WrittenFiles[path] = data
    return nil
}

func TestSaveResult_WhenFileWriteFails_ReturnsError(t *testing.T) {
    mockWriter := &MockFileWriter{ShouldFail: true}
    saver := NewResultSaver(mockWriter)
    
    err := saver.SaveResult("test.json", testData)
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "mock write failure")
}
```

#### ❌ **Bad Example:**
```go
func TestSaveResult_CreatesFile_FileExists(t *testing.T) {
    // Depends on actual filesystem
    err := SaveResult("test.json", testData)
    
    // This test can fail due to:
    // - Permission issues
    // - Disk space problems  
    // - File already existing
    // - Different OS file systems
    assert.NoError(t, err)
    assert.FileExists(t, "test.json") // Brittle assertion
}
```

### **4. Avoid Unnecessary Preconditions**

**Rule:** Don't create complex setup that multiple unrelated tests depend on.

#### ✅ **Good Example:**
```go
func TestProtocolParser_WithValidURL_ParsesCorrectly(t *testing.T) {
    // Simple, focused setup just for this test
    url := "siteparser://browser/?login=user&password=pass"
    
    result, err := ParseProtocolURL(url)
    
    assert.NoError(t, err)
    assert.Equal(t, "user", result.Login)
    assert.Equal(t, "pass", result.Password)
}
```

#### ❌ **Bad Example:**
```go
// Complex setup used by many tests
func setupComplexTestEnvironment() *TestEnvironment {
    // 50 lines of setup code
    // Creates files, initializes databases, sets up configs
    // Used by 20 different tests
    // If this changes, ALL 20 tests might break
}

func TestProtocolParser_WithValidURL_ParsesCorrectly(t *testing.T) {
    env := setupComplexTestEnvironment() // Unnecessary for URL parsing!
    
    result, err := env.Parser.ParseProtocolURL("siteparser://browser/?login=user")
    
    assert.NoError(t, err)
}
```

### **5. Don't Unit Test Configuration Settings**

**Rule:** Configuration is external to your code units. Don't test that you can copy and paste.

#### ✅ **Good Example:**
```go
// Test the LOGIC that uses config, not the config values themselves
func TestBrowserTimeout_WithCustomTimeout_UsesConfiguredValue(t *testing.T) {
    config := &Config{BrowserTimeout: 5 * time.Second}
    browser := NewBrowser(config)
    
    // Test that the browser USES the timeout, not that config contains it
    start := time.Now()
    browser.WaitForElement("selector", config.BrowserTimeout)
    duration := time.Since(start)
    
    assert.True(t, duration >= 5*time.Second)
}
```

#### ❌ **Bad Example:**
```go
// Testing configuration values (pointless)
func TestConfig_HasCorrectBrowserTimeout_Returns30Seconds(t *testing.T) {
    config := LoadConfig("config.yaml")
    
    // This just proves you can copy and paste the same value
    assert.Equal(t, 30*time.Second, config.BrowserTimeout)
    // If you change the config, you have to change this test too
    // It doesn't prove anything about your code's behavior
}
```

### **6. Name Your Unit Tests Clearly and Consistently**

**Rule:** Use the **Subject/Scenario/Result** (S/S/R) pattern.

#### ✅ **Good Examples:**
```go
// Subject: OpenBrowserWithLogin function
// Scenario: When Chrome is not installed  
// Result: Returns error with Chrome not found message
func TestOpenBrowserWithLogin_WhenChromeNotInstalled_ReturnsErrorWithChromeNotFoundMessage(t *testing.T) {
    // Test implementation
}

// Subject: ParseProtocolURL function
// Scenario: With malformed URL
// Result: Returns parsing error
func TestParseProtocolURL_WithMalformedURL_ReturnsParsingError(t *testing.T) {
    // Test implementation
}

// Subject: ValidateCredentials function  
// Scenario: With valid email and password
// Result: Returns success
func TestValidateCredentials_WithValidEmailAndPassword_ReturnsSuccess(t *testing.T) {
    // Test implementation
}
```

#### ❌ **Bad Examples:**
```go
func TestBrowser(t *testing.T) {
    // What about the browser? What scenario? What result?
}

func TestLogin(t *testing.T) {
    // Login what? Under what conditions? What should happen?
}

func TestError(t *testing.T) {
    // Which error? When? Why?
}
```

## 🏗️ **RPA-Specific Testing Patterns**

### **Testing Browser Automation Logic**

```go
func TestBrowserAutomation_WhenElementNotFound_ReturnsTimeoutError(t *testing.T) {
    mockBrowser := &MockBrowser{
        Elements: map[string]bool{
            "#login": false, // Element not found
        },
    }
    
    automation := NewBrowserAutomation(mockBrowser)
    
    err := automation.ClickElement("#login")
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "element not found")
    assert.Contains(t, err.Error(), "#login")
}
```

### **Testing Protocol Parsing**

```go
func TestProtocolHandler_WithTestParameter_ReturnsTestHandler(t *testing.T) {
    url := "siteparser://test"
    
    handler := GetHandlerFromProtocol(url)
    
    assert.IsType(t, &TestHandler{}, handler)
}

func TestProtocolHandler_WithEmailAndTokenParameters_ReturnsProcessHandler(t *testing.T) {
    url := "siteparser://browser/?email=test@email.com&token=abc123"
    
    handler := GetHandlerFromProtocol(url)
    
    assert.IsType(t, &ProcessHandler{}, handler)
}
```

### **Testing Error Handling**

```go
func TestFileUtils_WhenDirectoryNotWritable_ReturnsPermissionError(t *testing.T) {
    mockFS := &MockFileSystem{
        Permissions: map[string]bool{
            "/readonly": false,
        },
    }
    
    fileUtils := NewFileUtils(mockFS)
    
    err := fileUtils.SaveToFile("/readonly/test.json", []byte("data"))
    
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "permission denied")
}
```

## 📋 **Quick Reference Checklist**

Before writing a test, ask yourself:

- [ ] **Is this testing ONE specific behavior?**
- [ ] **Am I mocking all external dependencies?**
- [ ] **Will this test break if I refactor internal implementation?**
- [ ] **Does the test name clearly describe Subject/Scenario/Result?**
- [ ] **Can this test run independently of all other tests?**
- [ ] **Am I testing behavior, not configuration?**
- [ ] **Is this helping me design better code?**

## 🚫 **Common Anti-Patterns to Avoid**

### **1. The "Everything" Test**
```go
// DON'T: Test everything in one massive test
func TestCompleteWorkflow_DoesEverything_Successfully(t *testing.T) {
    // 100 lines of testing browser + files + config + logging + protocol...
}
```

### **2. The "Brittle" Test**
```go
// DON'T: Test implementation details
func TestBrowser_InternalState_HasCorrectFields(t *testing.T) {
    browser := NewBrowser()
    
    // Testing internal state instead of behavior
    assert.Equal(t, "chrome", browser.browserType)
    assert.Equal(t, 30, browser.timeoutSeconds)
}
```

### **3. The "Copy-Paste" Test**
```go
// DON'T: Test configuration values
func TestConfig_FacebookURL_IsCorrectURL(t *testing.T) {
    config := LoadConfig()
    
    // Just copying the same value from config to test
    assert.Equal(t, "https://facebook.com", config.FacebookURL)
}
```

---

> 💡 **Remember**: "Unit tests are a design specification of how a certain behavior should work, not a list of observations of everything the code happens to do." - [Steve Sanderson](https://gist.github.com/vadymhimself/763e96dd8495bb77325efd082e63c9f5) 