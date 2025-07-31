# Go Project Layout Rules and Standards

> Based on [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

This document outlines the standard Go project layout patterns and best practices for organizing Go projects.

## Core Go Directories

### `/cmd`
**Main applications for this project**

- The directory name for each application should match the name of the executable you want to have (e.g., `/cmd/myapp`)
- Don't put a lot of code in the application directory
- Import and invoke the code from the `/internal` and `/pkg` directories and nothing else
- It's common to have a small `main` function that imports and invokes the code from the `/internal` and `/pkg` directories

**Example structure:**
```
/cmd/
  /myapp/
    main.go
  /anotherapp/
    main.go
```

### `/internal`
**Private application and library code**

- Code you don't want others importing in their applications or libraries
- This layout pattern is enforced by the Go compiler itself
- Not limited to the top level `internal` directory - you can have more than one `internal` directory at any level of your project tree
- You can optionally add extra structure to separate shared and non-shared internal code

**Recommended structure:**
```
/internal/
  /app/          # Your actual application code
    /myapp/
  /pkg/          # Code shared by those apps
    /myprivlib/
```

### `/pkg`
**Library code that's ok to use by external applications**

- Other projects will import these libraries expecting them to work, so think twice before you put something here
- The `/internal` directory is a better way to ensure your private packages are not importable
- Good way to explicitly communicate that the code in that directory is safe for use by others
- Common layout pattern, but not universally accepted in the Go community

**When to use:**
- When your root directory contains lots of non-Go components and directories
- To group Go code in one place
- When you want to make it easier to run various Go tools

### `/vendor`
**Application dependencies**

- Managed manually or by your favorite dependency management tool
- The `go mod vendor` command will create the `/vendor` directory for you
- You might need to add the `-mod=vendor` flag to your `go build` command if not using Go 1.14+
- **Don't commit your application dependencies if you are building a library**

## Service Application Directories

### `/api`
**OpenAPI/Swagger specs, JSON schema files, protocol definition files**

- Store your API definitions
- Protocol definition files
- JSON schema files
- OpenAPI/Swagger specifications

## Web Application Directories

### `/web`
**Web application specific components**

- Static web assets
- Server side templates
- SPAs (Single Page Applications)

## Common Application Directories

### `/configs`
**Configuration file templates or default configs**

- Put your `confd` or `consul-template` template files here
- Default configuration files
- Configuration templates

### `/init`
**System init and process manager configs**

- System init (systemd, upstart, sysv) configs
- Process manager/supervisor (runit, supervisord) configs

### `/scripts`
**Scripts to perform various build, install, analysis, etc operations**

- Keep the root level Makefile small and simple
- Build scripts
- Installation scripts
- Analysis scripts
- Various automation scripts

### `/build`
**Packaging and Continuous Integration**

**Structure:**
```
/build/
  /package/     # Cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations
  /ci/          # CI (travis, circle, drone) configurations and scripts
```

**Note:** Some CI tools are picky about config file locations. Try putting config files in `/build/ci` and linking them to where CI tools expect them.

### `/deployments`
**IaaS, PaaS, system and container orchestration deployment configurations**

- Docker-compose files
- Kubernetes/Helm configurations
- Terraform configurations
- Deployment templates

**Note:** In some repos (especially apps deployed with Kubernetes) this directory is called `/deploy`.

### `/test`
**Additional external test apps and test data**

- Structure the `/test` directory however you want
- For bigger projects, consider having a data subdirectory: `/test/data` or `/test/testdata`
- Go will ignore directories or files that begin with "." or "_"

## Other Directories

### `/docs`
**Design and user documents**

- In addition to your godoc generated documentation
- Design documents
- User documentation
- Architecture documentation

### `/tools`
**Supporting tools for this project**

- These tools can import code from the `/pkg` and `/internal` directories
- Build tools
- Development tools
- Utility tools

### `/examples`
**Examples for your applications and/or public libraries**

- Usage examples
- Sample code
- Tutorials

### `/third_party`
**External helper tools, forked code and other 3rd party utilities**

- Swagger UI
- Forked code
- External helper tools
- 3rd party utilities

### `/githooks`
**Git hooks**

- Pre-commit hooks
- Post-commit hooks
- Other git hooks

### `/assets`
**Other assets to go along with your repository**

- Images
- Logos
- Icons
- Other media assets

### `/website`
**Project's website data**

- If you are not using GitHub pages
- Website source code
- Documentation sites

## Directories You Shouldn't Have

### ❌ `/src`
**Avoid this Java pattern in Go projects**

- Some Go projects do have a `src` folder, but it usually happens when developers come from the Java world
- You really don't want your Go code or Go projects to look like Java
- Don't confuse the project level `/src` directory with the `/src` directory Go uses for its workspaces

**Why avoid it:**
- It's a Java pattern, not a Go pattern
- Go has its own workspace structure with `$GOPATH`
- With Go modules (Go 1.11+), you can have your project outside of `GOPATH`

## Best Practices

### 1. Start Simple
- It's ok not to use all directories if your app project is really small
- Don't add extra levels of nesting that don't add value
- Think about structure when your project gets big enough

### 2. Think About External Usage
- Use `/internal` for code you don't want others to import
- Use `/pkg` for code that's safe for others to use
- Be intentional about what you expose

### 3. Follow Go Conventions
- Use `gofmt` for formatting
- Follow Go naming conventions
- Use `go vet` and `golint` for code quality

### 4. Documentation
- Use godoc comments
- Keep documentation up to date
- Provide examples

### 5. Testing
- Write tests for your code
- Use the standard `testing` package
- Consider table-driven tests

## Badges and Quality Indicators

Consider adding these badges to your README:

- **Go Report Card** - Scans code with `gofmt`, `go vet`, `gocyclo`, `golint`, `ineffassign`, `license` and `misspell`
- **Pkg.go.dev** - Go discovery & docs
- **Release** - Shows latest release number

## References

- [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- [How to Write Go Code](https://golang.org/doc/code.html)
- [Effective Go](https://golang.org/doc/effective_go.html)

---

> **Note:** This is a more opinionated project template. The Go community has varying opinions on some of these patterns, especially the `/pkg` directory. Use your judgment and consider your project's specific needs. 