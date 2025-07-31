# Go Project Layout Quick Reference

> Quick cheat sheet for Go project structure

## ✅ Recommended Directories

| Directory | Purpose | Notes |
|-----------|---------|-------|
| `/cmd` | Main applications | Small `main` functions, import from `/internal` and `/pkg` |
| `/internal` | Private code | Enforced by Go compiler, can't be imported externally |
| `/pkg` | Public libraries | Code safe for external use |
| `/api` | API definitions | OpenAPI/Swagger specs, JSON schemas |
| `/web` | Web components | Static assets, templates, SPAs |
| `/configs` | Configuration | Templates and default configs |
| `/scripts` | Build scripts | Keep root Makefile simple |
| `/build` | CI/CD & packaging | `/build/ci` and `/build/package` |
| `/deployments` | Deployment configs | Docker, K8s, Terraform |
| `/test` | Test data/apps | External test applications |
| `/docs` | Documentation | Beyond godoc |
| `/tools` | Supporting tools | Can import from `/pkg` and `/internal` |
| `/examples` | Usage examples | Sample code and tutorials |
| `/assets` | Media assets | Images, logos, icons |

## ❌ Avoid These Directories

| Directory | Why Avoid | Alternative |
|-----------|-----------|-------------|
| `/src` | Java pattern, not Go | Use root or proper Go structure |

## 🏗️ Typical Project Structure

```
myproject/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── myapp/
│   └── pkg/
│       └── shared/
├── pkg/
│   └── publiclib/
├── api/
│   └── openapi.yaml
├── web/
│   ├── static/
│   └── templates/
├── configs/
├── scripts/
├── build/
│   ├── ci/
│   └── package/
├── deployments/
├── docs/
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Quick Start Rules

1. **Start with `/cmd`** - Put your main applications here
2. **Use `/internal`** - For private code that shouldn't be imported
3. **Consider `/pkg`** - Only for code truly meant to be public
4. **Avoid `/src`** - It's a Java pattern, not Go
5. **Keep `main.go` small** - Import and call, don't implement

## 📝 Common Patterns

### Single Application
```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── model/
├── go.mod
└── go.sum
```

### Multiple Applications
```
myproject/
├── cmd/
│   ├── server/
│   │   └── main.go
│   ├── client/
│   │   └── main.go
│   └── tool/
│       └── main.go
├── internal/
│   ├── shared/
│   ├── server/
│   └── client/
├── go.mod
└── go.sum
```

### Library Project
```
mylib/
├── pkg/
│   ├── feature1/
│   └── feature2/
├── internal/
│   └── helpers/
├── examples/
├── docs/
├── go.mod
└── go.sum
```

## 🔍 Decision Tree

**Should I use `/pkg`?**
- ✅ Yes: If others will import your code
- ❌ No: If it's application-specific code

**Should I use `/internal`?**
- ✅ Yes: For code you don't want others importing
- ✅ Yes: Almost always a good choice for applications

**Should I use `/cmd`?**
- ✅ Yes: If you have executable applications
- ✅ Yes: Even for single applications

**Should I create `/src`?**
- ❌ No: Never use `/src` in Go projects

---

> 💡 **Remember**: Start simple and add structure as your project grows! 