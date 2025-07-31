# Go Project Layout Migration Guide

> How to migrate existing Go projects to follow standard layout patterns

## 🔄 Common Migration Scenarios

### Scenario 1: Everything in Root Directory

**Before:**
```
myproject/
├── main.go
├── handler.go
├── service.go
├── model.go
├── utils.go
├── go.mod
└── go.sum
```

**After:**
```
myproject/
├── cmd/
│   └── myproject/
│       └── main.go
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── service/
│   │   └── service.go
│   ├── model/
│   │   └── model.go
│   └── utils/
│       └── utils.go
├── go.mod
└── go.sum
```

**Migration Steps:**
1. Create `cmd/myproject/` directory
2. Move `main.go` to `cmd/myproject/main.go`
3. Create `internal/` directory
4. Group related files into subdirectories under `internal/`
5. Update import paths in all files

### Scenario 2: Using `/src` Directory (Java Pattern)

**Before:**
```
myproject/
├── src/
│   ├── main.go
│   ├── api/
│   ├── services/
│   └── models/
├── go.mod
└── go.sum
```

**After:**
```
myproject/
├── cmd/
│   └── myproject/
│       └── main.go
├── internal/
│   ├── api/
│   ├── service/
│   └── model/
├── go.mod
└── go.sum
```

**Migration Steps:**
1. Create `cmd/` and `internal/` directories
2. Move `main.go` from `src/` to `cmd/myproject/`
3. Move other packages from `src/` to `internal/`
4. Remove empty `src/` directory
5. Update all import paths

### Scenario 3: Mixed Public/Private Code

**Before:**
```
myproject/
├── main.go
├── publicapi/      # Should be public
├── utils/          # Should be private
├── config/         # Should be private
└── examples/       # Should stay as examples
```

**After:**
```
myproject/
├── cmd/
│   └── myproject/
│       └── main.go
├── pkg/
│   └── publicapi/
├── internal/
│   ├── utils/
│   └── config/
└── examples/
```

**Migration Steps:**
1. Create `cmd/`, `pkg/`, and `internal/` directories
2. Move `main.go` to `cmd/myproject/`
3. Move public packages to `pkg/`
4. Move private packages to `internal/`
5. Keep `examples/` at root level

## 📋 Migration Checklist

### Phase 1: Plan the Migration
- [ ] Identify all executable entry points
- [ ] Classify packages as public or private
- [ ] Plan new directory structure
- [ ] Document current import dependencies

### Phase 2: Create New Structure
- [ ] Create `cmd/` directory for applications
- [ ] Create `internal/` for private code
- [ ] Create `pkg/` for public libraries (if needed)
- [ ] Create other directories as needed (`api/`, `web/`, etc.)

### Phase 3: Move Files
- [ ] Move `main.go` files to appropriate `cmd/` subdirectories
- [ ] Move private packages to `internal/`
- [ ] Move public packages to `pkg/` (if applicable)
- [ ] Move configuration files to `configs/`
- [ ] Move scripts to `scripts/`

### Phase 4: Update Code
- [ ] Update all import statements
- [ ] Update build scripts and Makefiles
- [ ] Update documentation
- [ ] Update CI/CD configurations

### Phase 5: Test and Validate
- [ ] Run `go mod tidy`
- [ ] Test all builds: `go build ./...`
- [ ] Run all tests: `go test ./...`
- [ ] Verify external imports still work
- [ ] Update version control (git) if needed

## 🛠️ Migration Tools and Commands

### Check Current Module Structure
```bash
go list -m all
go list ./...
```

### Find Import Dependencies
```bash
go list -deps ./...
go mod graph
```

### Update Import Paths (Example)
```bash
# If using sed on Unix systems
find . -name "*.go" -exec sed -i 's|"myproject/|"myproject/internal/|g' {} \;

# Manual approach recommended for safety
grep -r "myproject/" . --include="*.go"
```

### Validate After Migration
```bash
go mod tidy
go build ./...
go test ./...
go vet ./...
```

## ⚠️ Common Pitfalls

### 1. Import Path Issues
**Problem:** Broken imports after moving files
**Solution:** 
- Use IDE refactoring tools when possible
- Update imports systematically
- Test frequently during migration

### 2. Circular Dependencies
**Problem:** Moving files reveals circular imports
**Solution:**
- Extract shared code to a separate package
- Redesign package boundaries
- Use dependency injection

### 3. External Dependencies
**Problem:** External packages importing your moved code
**Solution:**
- Keep public API stable during migration
- Use go.mod replace directives for testing
- Consider deprecation warnings

### 4. Build Script Issues
**Problem:** Build scripts and Makefiles break
**Solution:**
- Update all build commands
- Update Docker files
- Update CI/CD configurations

## 📝 Example Migration Script

```bash
#!/bin/bash
# migrate-to-standard-layout.sh

set -e

echo "🚀 Starting Go project migration to standard layout..."

# Backup current state
echo "📦 Creating backup..."
git add -A
git commit -m "Pre-migration backup" || echo "No changes to commit"

# Create new directories
echo "📁 Creating standard directories..."
mkdir -p cmd/myproject
mkdir -p internal
mkdir -p pkg

# Move main.go
echo "🔄 Moving main.go..."
mv main.go cmd/myproject/

# Move other files (customize as needed)
echo "🔄 Moving internal packages..."
for dir in handler service model utils; do
    if [ -d "$dir" ]; then
        mv "$dir" internal/
    fi
done

# Update go.mod if needed
echo "🔧 Updating dependencies..."
go mod tidy

# Test the migration
echo "🧪 Testing migration..."
go build ./...
go test ./...

echo "✅ Migration completed successfully!"
echo "🔍 Please review and test your application thoroughly"
echo "📝 Don't forget to update documentation and CI/CD configs"
```

## 🔗 Resources for Migration

- [Official Go Documentation](https://golang.org/doc/)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

---

> ⚡ **Pro Tip**: Migrate gradually and test frequently. Consider doing the migration in a feature branch to easily rollback if needed. 