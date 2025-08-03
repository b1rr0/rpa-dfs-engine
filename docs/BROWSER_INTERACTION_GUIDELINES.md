# Browser Interaction Guidelines

## Overview

This document outlines the standard practices for all browser interactions within the RPA DFS Engine project.

## ⚠️ Important Rule

**All browser interactions MUST use the `internal/browser` package.**

## Why Use `internal/browser`?

The `internal/browser` package provides:

- ✅ Standardized ChromeDP configuration
- ✅ Consistent error handling and logging
- ✅ Unified browser automation patterns
- ✅ Proper Chrome installation detection
- ✅ Consistent result handling via `types.BrowserResult`

## Usage Examples

### Basic Browser Operation

```go
import "rpa-dfs-engine/internal/browser"

// Correct way to perform browser automation
result := browser.OpenBrowserWithLogin(username, password)

if result.Success {
    logger.LogSuccess("Browser automation successful")
    logger.LogInfo("URL: %s", result.URL)
} else {
    logger.LogError("Browser error: %s", result.Error)
}
```

### Integration Points

The `internal/browser` package is already integrated with:

- **CLI Handler** (`internal/cli/cli.go`)
- **Protocol Handler** (`internal/protocol/protocol.go`)
- **Test Handler** (`internal/handlers/test_handler.go`)

## What NOT to Do

❌ **Don't** create direct ChromeDP instances outside of `internal/browser`
❌ **Don't** implement custom browser automation without using the standard package
❌ **Don't** bypass the centralized browser configuration

## Browser Configuration

The `internal/browser` package includes:

- Headless/non-headless mode control
- Security flags and sandbox settings
- GPU and memory optimization
- Chrome installation validation
- Timeout handling

## Result Handling

All browser operations return a `types.BrowserResult` structure:

```go
type BrowserResult struct {
    Success   bool
    URL       string
    Username  string
    Message   string
    Error     string
    Timestamp int64
}
```

## Adding New Browser Features

When adding new browser automation features:

1. **Extend** the `internal/browser` package
2. **Follow** existing patterns and error handling
3. **Return** `types.BrowserResult` for consistency
4. **Use** the established logging patterns
5. **Test** Chrome installation detection

## Related Documentation

- [Go Project Layout Rules](go-standards/GO_PROJECT_LAYOUT_RULES.md)
- [Project Structure Analysis](PROJECT_STRUCTURE_ANALYSIS.md)

---

> 💡 **Remember**: Consistency in browser interactions ensures reliability, maintainability, and easier debugging across the entire RPA engine. 