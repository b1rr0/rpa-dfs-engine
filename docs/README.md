# Go Robot Project Documentation

Welcome to the Go Robot project documentation! This directory contains all the documentation and standards for development.

## 📁 Documentation Structure

```
docs/
├── README.md                    # This file - main documentation index
└── go-standards/               # Go project layout standards and guidelines
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
2. **During development**: Keep the quick reference handy
3. **When refactoring**: Use the migration guide

### For Code Reviews
- Ensure new projects follow the standard layout
- Reference these documents when suggesting improvements
- Check that directory structure matches the standards

### For CI/CD
- Include structure validation in your build pipeline
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