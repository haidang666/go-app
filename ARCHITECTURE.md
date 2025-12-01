# Go Application Architecture Analysis

**Project**: go-app  
**Repository**: github.com/haidang666/go-app  
**Go Version**: 1.25.4  
**Analysis Date**: December 2, 2025

---

## Architecture Pattern: Clean Architecture with Hexagonal (Ports & Adapters)

The application follows Clean Architecture principles combined with Hexagonal Architecture. The design ensures:

- Clear separation of concerns across multiple layers
- Dependency inversion through interface-based contracts
- Framework independence - domain layer is pure business logic
- Testability - each layer can be tested independently
- Ports & Adapters - interfaces define ports, implementations are adapters

### Layers

```
External (cmd/)
    |
Bootstrap (DI setup)
    |
Presentation (HTTP handlers)
    |
Application (Use cases)
    |
Domain (Business logic - no dependencies)
    |
Infrastructure (Implementations)
```

---

## Directory Structure

### Root Organization

- cmd/ - Entry points
- internal/ - Private application code
  - api/ - API DTOs (SignUpRequest)
  - bootstrap/ - DI setup (Wire)
  - config/ - Configuration (envconfig)
  - domain/ - Business logic
    - contract/ - Interfaces (UserRepository)
    - dto/ - Domain DTOs (SignUpInput)
    - entity/ - Entities (User)
    - use_case/ - Use cases (SignUpUseCase)
  - infrastructure/ - Implementations
    - http/ - HTTP handlers and router
    - repository/ - Data access (UserRepository)
- pkg/ - Public reusable utilities
  - http/ - Request/response parsing
  - jwt/ - JWT token utilities
  - logger/ - Structured logging

### Layer Responsibilities

| Layer | Location | Purpose |
|-------|----------|---------|
| External | cmd/server/main.go | Entry point, bootstrap, graceful shutdown |
| Presentation | internal/api/, infrastructure/http/handlers/ | HTTP request handling, validation |
| Application | internal/domain/use_case/ | Business logic orchestration |
| Domain | internal/domain/ | Core business rules, entities, contracts |
| Infrastructure | internal/infrastructure/ | Repository implementations, config |

---

## Key Components

### 1. Entry Point: cmd/server/main.go

- Loads configuration from environment
- Creates DI container with Wire
- Starts HTTP server with graceful shutdown
- Handles OS signals (SIGINT, SIGTERM)

### 2. Configuration: internal/config/config.go

- Loads from environment variables (.env and ENV)
- APP_PORT, DB_HOST, DB_PORT, DB_NAME, DB_USERNAME, DB_PASSWORD

### 3. Bootstrap/DI: internal/bootstrap/

**wire.go**: Defines provider functions for Wire
- ProvideUserRepository() -> UserRepository
- ProvideSignUpUseCase() -> SignUpUseCase
- ProvideAuthHandler() -> AuthHandler
- ProvideRouter() -> Router
- ProvideContainer() -> Container

**wire_gen.go**: Auto-generated initialization code
**container.go**: Service container definition
**server.go**: HTTP server startup with graceful shutdown

### 4. Domain Layer: internal/domain/

**entity/user_entity.go**: User domain entity with validation
- ID (UUID), Email, HashedPassword, CreatedAt, UpdatedAt
- Validate() method for business rules

**contract/user_repository.go**: Repository interface (port)
- Create(ctx, user) -> (user, error)

**dto/sign_up_input.go**: Use case input DTO
- Email, Password

**use_case/auth/sign_up.go**: Sign-up use case
1. Hash password (bcrypt, cost 10)
2. Create User entity
3. Validate entity
4. Persist via repository
5. Return created user

### 5. Infrastructure: internal/infrastructure/

**repository/user_repository.go**: Repository implementation (adapter)
- Generates UUID
- Normalizes email to lowercase
- Currently in-memory (ready for database)

**http/router/router.go**: Chi router with middleware
- RequestID: Add request ID to context
- RealIP: Extract real client IP
- Logger: Log requests
- Recoverer: Panic recovery
- Routes: /health, /api/v1/auth/sign-up

**http/handlers/auth/handler.go**: HTTP handler
- Parse JSON request
- Validate request
- Convert to domain DTO
- Execute use case
- Return JSON response

**http/handlers/auth/routes.go**: Route registration
- POST /auth/sign-up -> AuthHandler.SignUp()

### 6. API Layer: internal/api/auth/sign_up_request.go

HTTP request DTO:
- Email, Password
- Validate() method for API-level validation

### 7. Utilities: pkg/

**http/request/parser.go**:
- FromJSON() - Parse and validate JSON (1MB limit, reject unknown fields)
- ToJSON() - Serialize and write JSON response

**jwt/jwt.go**:
- Client for JWT generation and verification
- HS256 algorithm

**logger/logger.go**:
- Singleton Zap logger
- Development vs production modes
- Structured logging

---

## Data Flow: Sign-Up

1. HTTP POST /api/v1/auth/sign-up with JSON body
2. AuthHandler.SignUp()
   - FromJSON() -> SignUpRequest
   - SignUpRequest.Validate()
   - Convert to dto.SignUpInput
3. SignUpUseCase.Execute()
   - Hash password (bcrypt)
   - Create User entity
   - User.Validate()
   - userRepo.Create()
4. UserRepository.Create()
   - Generate UUID
   - Normalize email
   - Return User
5. AuthHandler serializes response
6. HTTP 201 Created with User JSON

---

## Dependency Injection Pattern

### Google Wire Framework

Wire generates compile-time, type-safe dependency injection code.

**Benefits**:
- Zero-cost (generated code, no reflection)
- Type-safe (compile-time validation)
- Clear initialization order
- No runtime overhead

**Dependency Graph**:
```
Container
  <- Router
    <- AuthHandler
      <- SignUpUseCase
        <- UserRepository interface
          <- UserRepository implementation
```

### Providers (wire.go)

Each provider creates one dependency:
1. ProvideUserRepository() - Creates implementation
2. ProvideSignUpUseCase() - Receives repository interface
3. ProvideAuthHandler() - Receives use case
4. ProvideRouter() - Receives handler
5. ProvideContainer() - Receives router

### Generated Code (wire_gen.go)

Wire generates InitializeContainer() that executes providers in dependency order.

### Constructor Injection

All dependencies passed through constructors:
```go
NewSignUpUseCase(userRepo contract.UserRepository)
NewAuthHandler(args NewAuthHandlerArgs)
```

### Interface-Based Dependencies

```go
// Use case depends on interface
type SignUpUseCase struct {
    userRepo contract.UserRepository  // Interface!
}

// Allows swapping implementations
```

---

## Technology Stack

| Technology | Version | Purpose |
|-----------|---------|---------|
| chi/v5 | 5.2.3 | HTTP router |
| Google Wire | 0.7.0 | DI code generation |
| envconfig | 1.4.0 | Env variable loading |
| godotenv | 1.5.1 | .env file support |
| Zap | 1.27.1 | Structured logging |
| golang-jwt/jwt | 5.3.0 | JWT tokens |
| go-playground/validator | 10.28.0 | Validation |
| google/uuid | 1.3.0 | UUID generation |
| golang.org/x/crypto | 0.45.0 | Bcrypt hashing |

---

## Design Patterns Used

1. **Clean Architecture** - Layers with dependency inversion
2. **Hexagonal/Ports & Adapters** - Interface-based contracts
3. **Repository** - Data access abstraction
4. **Use Case/Interactor** - Business logic encapsulation
5. **Handler/Controller** - HTTP request handling
6. **DTO** - Data transfer objects per layer
7. **Entity Validation** - Business rules in entities
8. **Singleton** - Logger with sync.Once
9. **Builder/Args Struct** - Flexible constructors
10. **Dependency Injection** - Constructor-based via Wire
11. **Strategy** - Pluggable repository implementations
12. **Graceful Shutdown** - Signal-based clean exit

---

## Layer Interactions

### Request Processing Pipeline

```
HTTP Request
    |
Chi Router Middleware (RequestID, RealIP, Logger, Recoverer)
    |
AuthHandler.SignUp()
  - Parse JSON (request.FromJSON)
  - Validate request (SignUpRequest.Validate)
  - Convert to domain DTO (SignUpInput)
    |
SignUpUseCase.Execute()
  - Hash password (bcrypt)
  - Create entity
  - Validate entity (User.Validate)
  - Call repository
    |
UserRepository.Create()
  - Generate UUID
  - Normalize email
  - Return User
    |
AuthHandler serializes response
    |
HTTP 201 Created with JSON
```

### Dependency Direction

External -> Bootstrap -> Infrastructure -> Domain

Domain depends on nothing. One-way flow ensures domain independence.

---

## Entry Points and Initialization

### Primary: cmd/server/main.go

1. config.Load() - Load environment configuration
2. signal.NotifyContext() - Setup OS signal handling
3. bootstrap.CreateServerContainer() - Initialize DI container
   - Executes Wire-generated InitializeContainer()
   - Creates all services in dependency order
4. bootstrap.StartRestAPI() - Start HTTP server
   - Create http.Server with router
   - Listen in background goroutine
   - Wait for shutdown signal or error
5. On shutdown signal
   - Call server.Shutdown() with 10s timeout
   - Gracefully terminate in-flight requests
6. Container.Close() - Cleanup

### Initialization Order

1. Configuration loading
2. DI container creation (Wire)
3. Router initialization with all middleware
4. HTTP server startup
5. Ready for requests

---

## Code Organization Best Practices

### File Naming
- Snake_case: user_entity.go, sign_up.go, sign_up_request.go
- One primary type per file

### Package Naming
- Lowercase, single word when possible
- No underscores in package names
- Named after directory

### Type/Function Naming
- Exported: PascalCase
- Unexported: camelCase
- Constants: ALL_CAPS

### Import Organization
1. Standard library
2. Third-party packages
3. Internal packages
(Separated by blank lines)

### Error Handling
- Check immediately after operation
- Return error as last value
- Wrap with context: fmt.Errorf("context: %w", err)
- No silent failures

### Interface Design
- Small, focused interfaces
- Single responsibility principle
- Named by consumer, not implementer

### Constructor Pattern
Simple: NewType(dep Type) *Type
Complex: NewType(args NewTypeArgs) *Type
(Args struct allows easy parameter addition)

### Validation Strategy
- API Layer: Request format, data types, lengths
- Domain Layer: Business rules, constraints

---

## Architecture Quality Attributes

### Maintainability - Excellent
- Clear separation of concerns
- Single responsibility per component
- Consistent naming throughout
- File structure mirrors layers

### Testability - Excellent
- DI enables mocking
- Domain has zero external dependencies
- Repository interface for test doubles
- Framework-independent business logic

### Scalability - Good
- Modular structure for new features
- Use case pattern for business logic
- Repository pattern for data sources
- Handler pattern for endpoints

### Flexibility - Excellent
- Repository pattern abstracts storage
- Wire DI supports adding/removing dependencies
- DTO pattern isolates layers
- Interface contracts enable substitution

### Extensibility - Excellent
- Clear extension points
- Established patterns
- Minimal changes to existing code
- New features in isolated packages

### Robustness - Good
- Graceful shutdown handling
- Error handling throughout
- Request body size limits
- Unknown field rejection
- Password hashing (security)

---

## Summary

This is a well-architected Go application demonstrating:

### Architecture
- Clean Architecture with hexagonal patterns
- Proper layering with dependency inversion
- Framework-independent domain logic
- Interface-based loose coupling

### Key Strengths
1. Clear separation of concerns
2. Type-safe dependency injection (Wire)
3. Testable business logic
4. Easy to extend and maintain
5. Production-ready structure
6. Consistent code organization

### Design Decisions
- Wire over manual DI for type safety
- Repository pattern for data abstraction
- Use case pattern for business logic
- Chi router for lightweight HTTP routing
- Zap for structured logging
- Two-layer validation (API + Domain)

### Production Ready
- Graceful shutdown
- Error handling
- Configuration management
- Structured logging
- Password hashing
- Request validation

### Best For
- Reference implementation for Go projects
- Learning Clean Architecture
- Understanding DI patterns
- Production application baseline
- Team best practices guide

This application is an excellent example of how to build maintainable, testable, and scalable Go applications with proper architectural separation and design patterns.

