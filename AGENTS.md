# AGENTS.md - Go App Development Guidelines

## Build, Lint & Test Commands

- **Format code**: `make format` or `go fmt ./...`
- **Lint code**: `make lint` or `go vet ./...`
- **Run server**: `make run` or `go run ./cmd/server/main.go`
- **Install dependencies**: `make install` or `go mod tidy && go mod download`
- **Run tests**: `go test ./...` (single test: `go test -run TestName ./path/to/package`)

## Code Style & Conventions

### Imports
- Group imports: standard library, third-party, internal (separated by blank lines)
- Use absolute module paths: `github.com/haidang666/go-app/internal/...`

### Naming
- **Files**: use snake_case (e.g., `user_entity.go`, `sign_up.go`)
- **Packages**: lowercase, no underscores, named after directory
- **Functions**: CamelCase, exported functions capitalized (e.g., `NewSignUpUseCase()`, `Execute()`)
- **Variables**: camelCase for local, PascalCase for exported struct fields
- **Constants**: ALL_CAPS with underscores

### Type System
- Use interfaces for dependency injection (repository pattern)
- Avoid `interface{}`, use concrete types or generics
- Use pointer receivers for methods that modify state

### Error Handling
- Return `error` as last return value
- Check errors immediately: `if err != nil { return nil, err }`
- Use `fmt.Errorf()` for wrapped errors with context
- No silent failures; log or propagate all errors

### Formatting & Structure
- Max line length: ~80-100 characters (Go standard)
- Use `go fmt` automatically (integrated with linters)
- Single responsibility: one struct per file when possible
- Dependency injection via constructors (e.g., `NewSignUpUseCase(userRepo)`)
