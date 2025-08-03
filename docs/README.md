# Go Robot Project Documentation

Welcome to the Go Robot project documentation! This directory contains all the documentation and standards for development.

## 📁 Documentation Structure

```
docs/
├── README.md                           # This file - main documentation index
├── BROWSER_INTERACTION_GUIDELINES.md  # Browser automation standards and guidelines
├── PROJECT_STRUCTURE_ANALYSIS.md      # Project structure analysis
├── test/                               # Testing documentation and best practices
│   ├── TESTING.md                       # Main testing guide
│   ├── UNIT_TESTING_BEST_PRACTICES.md # Detailed unit testing principles
│   ├── GO_TESTING_GUIDE.md            # Go-specific testing techniques
└── go-standards/                       # Go project layout standards and guidelines
    ├── GO_PROJECT_LAYOUT_RULES.md      # Complete Go layout rules
    ├── GO_LAYOUT_QUICK_REFERENCE.md    # Quick reference cheat sheet
    └── GO_LAYOUT_MIGRATION_GUIDE.md    # Migration guide for existing projects
```

## 🚀 Quick Start

### For New Go Projects
1. Read the [Go Project Layout Rules](go-standards/GO_PROJECT_LAYOUT_RULES.md) for comprehensive guidelines
2. Use the [Quick Reference](go-standards/GO_LAYOUT_QUICK_REFERENCE.md) as a cheat sheet during development

### For Existing Projects
1. Review the [Migration Guide](go-standards/GO_LAYOUT_MIGRATION_GUIDE.md)
2. Follow the step-by-step instructions to restructure your project

## 🤖 Browser Automation Guidelines

### [🌐 Browser Interaction Guidelines](BROWSER_INTERACTION_GUIDELINES.md)
**Essential guidelines** for all browser automation within the RPA DFS Engine.

**What you'll find:**
- Mandatory use of `internal/browser` package
- Standard browser automation patterns
- ChromeDP configuration best practices
- Error handling and result patterns
- Do's and don'ts for browser interactions

## 🧪 Testing Documentation

### [📋 Testing Guide](test/README.md)
**Comprehensive testing strategy** for the RPA DFS Engine desktop application.

**What you'll find:**
- Testing philosophy: TDD as a design process
- RPA-specific testing approaches
- Unit vs Integration vs E2E testing strategies
- Mock patterns for browser automation
- Cross-platform testing considerations

### [⚡ Unit Testing Best Practices](test/UNIT_TESTING_BEST_PRACTICES.md)
**Detailed principles** based on [Steve Sanderson's testing guide](https://gist.github.com/vadymhimself/763e96dd8495bb77325efd082e63c9f5).

**What you'll find:**
- The six golden rules of unit testing
- Good vs bad test examples
- RPA-specific testing patterns
- Common anti-patterns to avoid
- Testing checklist and guidelines

### [🔧 Go Testing Guide](test/GO_TESTING_GUIDE.md)
**Practical Go testing techniques** for desktop RPA automation.

**What you'll find:**
- Go testing framework usage
- Mocking patterns and interfaces
- Test file organization
- Coverage reporting and benchmarking
- Performance testing strategies

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
go test ./test/unit/types_test.go
```

### **Individual Tests**
```bash
# Run specific test function
go test -run TestBrowserResult_WithSuccessfulResult ./test/unit/

# Run tests matching pattern
go test -run TestProtocol ./test/unit/
```


## 📚 Go Standards Documentation

### [📖 Go Project Layout Rules](go-standards/GO_PROJECT_LAYOUT_RULES.md)
**Complete comprehensive guide** covering all standard Go project layout patterns based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout).

**What you'll find:**
- Detailed explanations of each directory (`/cmd`, `/internal`, `/pkg`, etc.)
- Best practices and conventions
- What to avoid and why
- References and additional resources

### [⚡ Quick Reference](go-standards/GO_LAYOUT_QUICK_REFERENCE.md)
**Fast cheat sheet** for daily development use.

**What you'll find:**
- Directory purpose table
- Common project structure examples
- Decision tree for directory usage
- Quick start rules

### [🔄 Migration Guide](go-standards/GO_LAYOUT_MIGRATION_GUIDE.md)
**Step-by-step guide** for migrating existing Go projects to standard layout.

**What you'll find:**
- Common migration scenarios
- Phase-by-phase migration checklist
- Migration tools and commands
- Common pitfalls and solutions
- Example migration scripts

## 🎯 How to Use These Standards

### For Development Team
1. **Before starting a new Go project**: Read the layout rules
2. **Before implementing browser automation**: Read the browser interaction guidelines
3. **Before writing code**: Read the testing best practices and use TDD
4. **During development**: Keep the quick reference and testing guide handy
5. **When refactoring**: Use the migration guide and maintain test coverage

### For Code Reviews
- Ensure new projects follow the standard layout
- Verify browser interactions use `internal/browser` package
- Check that new code includes appropriate unit tests
- Validate test naming follows Subject/Scenario/Result pattern
- Ensure external dependencies are properly mocked
- Reference these documents when suggesting improvements

### For CI/CD
- Include structure validation in your build pipeline
- Run unit tests on every commit with coverage reporting
- Execute integration tests in the pipeline
- Use cross-platform testing matrix (Windows, macOS, Linux)
- Enforce minimum test coverage thresholds
- Use the migration scripts as templates for automation

## 📋 Development Workflow

### New Project Setup
```bash
# 1. Create project structure
mkdir -p cmd/myapp internal pkg

# 2. Initialize Go module
go mod init myproject

# 3. Create main.go in cmd/myapp/
# 4. Follow the layout rules for additional directories
```

### Code Organization
- **Executables** → `/cmd`
- **Private code** → `/internal`
- **Public libraries** → `/pkg`
- **API specs** → `/api`
- **Web assets** → `/web`
- **Configs** → `/configs`
- **Scripts** → `/scripts`

## 🔗 External References

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) - Original source
- [Official Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Modules Reference](https://golang.org/ref/mod)

## 📝 Contributing to Documentation

To improve or update these standards:

1. Follow the same standards documented here
2. Keep documentation clear and practical
3. Include examples and real-world scenarios
4. Test any scripts or commands before documenting them

---

> 💡 **Remember**: These standards are meant to help, not hinder. Adapt them to your project's specific needs while maintaining the core principles. 