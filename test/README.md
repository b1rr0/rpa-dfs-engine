# RPA DFS Engine Test Suite

This directory contains all tests for the RPA DFS Engine, organized by type and purpose following Go testing best practices.

## 📁 **Directory Structure**

```
test/
├── README.md                          # This file - testing documentation
├── unit/                              # Unit tests
│   ├── traverser_test.go               # Comprehensive traverser component tests
│   ├── browser_test.go                 # Browser automation tests
│   ├── handlers_test.go                # Handler logic tests
│   ├── cli_test.go                     # CLI parameter parsing tests
│   ├── protocol_test.go                # Protocol handling tests
│   ├── logger_test.go                  # Logging functionality tests
│   ├── fileutils_test.go               # File utilities tests
│   ├── config_test.go                  # Configuration tests
│   └── types_test.go                   # Type definition tests
├── integration/                       # Integration tests
│   ├── traverser_integration_test.go   # End-to-end traverser workflow tests
│   └── workflow_test.go                # Complete workflow integration tests
└── mocks/                             # Mock implementations
    ├── mock_traverser_browser.go       # Mock browser for traverser testing
    ├── mock_browser.go                 # Legacy browser mock
    └── mock_filewriter.go              # File writing mock
```

## 🚀 **Running Tests**

### **All Tests**
```bash
# Run all tests
go test ./test/...

# With verbose output
go test -v ./test/...

# With coverage report
go test -cover ./test/...
```

### **Specific Test Suites**
```bash
# Unit tests only
go test ./test/unit/

# Integration tests only
go test ./test/integration/

# Specific test file
go test ./test/unit/traverser_test.go
```

### **Individual Tests**
```bash
# Run specific test function
go test -run TestParser_LoadWorkflow ./test/unit/

# Run tests matching pattern
go test -run TestTraverser ./test/...
```

## 🧪 **Test Categories**

### **Unit Tests** (`/test/unit/`)

#### **Traverser Tests** (`traverser_test.go`)
- **Parser Tests**: JSON workflow and context loading/validation
- **Context Tests**: Template variable resolution and data access
- **Template Resolver Tests**: `{{user.property}}` substitution and condition evaluation
- **Executor Tests**: Node execution, control flow, and error handling
- **High-level API Tests**: Workflow validation and convenience functions

#### **Browser Tests** (`browser_test.go`)
- ChromeDP configuration and browser automation
- Error handling and timeout management
- Cross-platform Chrome detection

#### **Handler Tests** (`handlers_test.go`)
- Protocol URL parsing and parameter extraction
- Handler selection logic (Setup, Test, Process)
- Request routing and validation

#### **CLI Tests** (`cli_test.go`)
- Command-line flag parsing
- Parameter validation and defaults
- Help message generation

### **Integration Tests** (`/test/integration/`)

#### **Traverser Integration** (`traverser_integration_test.go`)
- **End-to-End Workflows**: Complete workflow execution with mock browser
- **Error Handling**: Browser failures, invalid workflows, missing files
- **Template Resolution**: Complex nested template substitution
- **Workflow Validation**: Comprehensive schema validation tests

#### **Legacy Workflow Tests** (`workflow_test.go`)
- Protocol handling integration
- File output generation
- Handler workflow coordination

### **Mocks** (`/test/mocks/`)

#### **MockTraverserBrowser** (`mock_traverser_browser.go`)
- Implements `traverser.Browser` interface
- Records all browser actions for verification
- Configurable failure simulation
- Helper methods for test assertions

#### **Legacy Mocks**
- `MockBrowserAutomation`: Facebook-specific browser mock
- `MockFileWriter`: File operation mock

## ✅ **Testing Standards**

### **Test Naming Convention**
- **Functions**: `TestComponentName_Scenario_ExpectedResult`
- **Subtests**: `t.Run("Descriptive scenario name", func(t *testing.T) {...})`

### **Test Organization**
```go
func TestComponentName_Functionality(t *testing.T) {
    t.Run("Positive case scenario", func(t *testing.T) {
        // Setup
        // Execute  
        // Assert
    })
    
    t.Run("Error case scenario", func(t *testing.T) {
        // Setup with error conditions
        // Execute
        // Assert error handling
    })
}
```

### **Assertion Guidelines**
- Use `require.NoError(t, err)` for critical failures
- Use `assert.Equal(t, expected, actual)` for value comparisons
- Use `assert.True/False(t, condition)` for boolean checks
- Use `assert.Contains(t, haystack, needle)` for substring/element checks

### **Mock Usage**
```go
// Create mock
mockBrowser := mocks.NewMockTraverserBrowser()

// Configure failure
mockBrowser.SetFailureFor("navigate", "connection failed")

// Execute test
err := executor.Execute()

// Verify behavior
assert.True(t, mockBrowser.HasNavigatedTo("https://example.com"))
assert.Equal(t, "test@example.com", mockBrowser.GetFieldValue("#email"))
```

## 🎯 **Test Coverage**

### **Current Coverage Areas**

#### **✅ Comprehensive Coverage**
- **Traverser Components**: Parser, Context, Templates, Executor
- **Workflow Execution**: All node types and control flow
- **Error Handling**: Validation, browser failures, malformed input
- **Template Resolution**: Variable substitution and condition evaluation

#### **✅ Good Coverage** 
- **Browser Automation**: ChromeDP integration and configuration
- **Protocol Handling**: URL parsing and parameter extraction
- **CLI Processing**: Flag parsing and validation
- **Logging**: Message formatting and file handling

#### **✅ Basic Coverage**
- **File Operations**: Reading, writing, path handling
- **Configuration**: Environment variables and constants
- **Type Definitions**: JSON marshaling/unmarshaling

### **Testing Best Practices**

1. **Test Isolation**: Each test should be independent and not rely on external state
2. **Mock External Dependencies**: Use mocks for browser, file system, network calls
3. **Test Both Success and Failure Paths**: Verify error handling and edge cases
4. **Use Descriptive Test Names**: Tests should clearly describe what they verify
5. **Arrange-Act-Assert Pattern**: Structure tests with clear setup, execution, and verification
6. **Clean Up Resources**: Use `defer` statements for cleanup or `t.TempDir()` for temporary files

### **Adding New Tests**

#### **For New Components**
1. Create test file: `test/unit/component_test.go`
2. Follow naming convention: `TestComponentName_Functionality`
3. Add both positive and negative test cases
4. Use appropriate mocks for dependencies

#### **For Integration Scenarios**
1. Add to: `test/integration/`
2. Focus on end-to-end workflows
3. Test real component interactions
4. Verify complete user scenarios

#### **Mock Requirements**
1. Implement required interfaces completely
2. Provide helper methods for test assertions
3. Allow configurable failure simulation
4. Record all interactions for verification

## 📊 **Test Results**

Current test status: **✅ All tests passing**

- **Unit Tests**: 100+ test cases covering all components
- **Integration Tests**: End-to-end workflow scenarios
- **Coverage**: Comprehensive coverage of critical paths
- **Performance**: Fast execution with efficient mocking

## 🔧 **Troubleshooting**

### **Common Test Issues**

#### **Test Failures**
```bash
# Run specific failing test with verbose output
go test -v -run TestSpecificName ./test/unit/

# Check for race conditions
go test -race ./test/...
```

#### **Mock Issues**
- Ensure mock interfaces match actual implementations
- Verify mock expectations are properly set
- Check that mocks are reset between tests

#### **File Path Issues**
- Use `t.TempDir()` for temporary test files
- Use relative paths in test data
- Clean up test files in `defer` statements

### **Performance Testing**
```bash
# Run benchmarks
go test -bench=. ./test/...

# Generate CPU profile
go test -cpuprofile=cpu.prof ./test/...

# Generate memory profile  
go test -memprofile=mem.prof ./test/...
```

---

> 💡 **Remember**: Tests are documentation! They show how components should be used and what behaviors are expected. Keep them clear, comprehensive, and maintainable. 