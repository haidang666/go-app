# Go Application Architecture Analysis

## 1. Architecture Pattern

This application follows **Clean Architecture** (also known as Hexagonal Architecture or Ports & Adapters), with clear separation of concerns across multiple layers:

- **Presentation Layer**: HTTP handlers and API request/response objects
- **Application Layer**: Use cases that orchestrate business logic
- **Domain Layer**: Business entities, contracts (interfaces), and DTOs
- **Infrastructure Layer**: Repository implementations and HTTP infrastructure
- **Configuration & Bootstrap Layer**: Dependency injection and application initialization

This design ensures:
- Independent testability at each layer
- Framework/technology agnostic domain logic
- Clear dependency flow (inward dependencies only)
- Easy to extend and maintain

---

## 2. Directory Structure & Purpose

```
go-app/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/                     # API request objects (DTOs for external contracts)
│   │   └── auth/
│   │       └── sign_up_request.go
│   ├── bootstrap/               # Dependency injection & initialization
│   │   ├── container.go
│   │   ├── server.go
│   │   ├── wire.go              # Wire provider setup
│   │   └── wire_gen.go          # Auto-generated wire code
│   ├── config/                  # Application configuration
│   │   └── config.go
│   ├── domain/                  # Core business logic (no external dependencies)
│   │   ├── contract/            # Interfaces (ports)
│   │   │   └── user_repository.go
│   │   ├── dto/                 # Domain transfer objects
│   │   │   └── sign_up_input.go
│   │   ├── entity/              # Business entities
│   │   │   └── user_entity.go
│   │   └── use_case/            # Application use cases
│   │       └── auth/
│   │           └── sign_up.go
│   └── infrastructure/          # External dependencies & implementations
│       ├── http/
│       │   ├── handlers/        # HTTP request handlers
│       │   │   └── auth/
│       │   │       ├── handler.go
│       │   │       └── routes.go
│       │   └── router/          # Router setup
│       │       └── router.go
│       └── repository/          # Repository implementations
│           └── user_repository.go
└── pkg/                         # Shared utility packages (no business logic)
    ├── http/
    │   └── request/             # HTTP request/response helpers
    │       └── parser.go
    ├── jwt/
    │   └── jwt.go
    └── logger/
        └── logger.go
```

### Layer Responsibilities

| Directory | Purpose |
|-----------|---------|
| **cmd/** | Entry point; initializes the application |
| **internal/api/** | API-specific DTOs for request/response contracts |
| **internal/bootstrap/** | Dependency injection container and wire configuration |
| **internal/config/** | Environment configuration loading |
| **internal/domain/** | Pure business logic; no external dependencies |
| **internal/infrastructure/** | Concrete implementations of domain contracts |
| **pkg/** | Reusable utilities; available for internal or external use |

---

## 3. Key Components & Their Responsibilities

### 3.1 Entry Point: `cmd/server/main.go`

- Loads configuration from environment variables
- Initializes the dependency injection container
- Starts the HTTP server with graceful shutdown handling
- Coordinates signal handling (SIGTERM, SIGINT)

```go
func main() {
    cfg, err := config.Load()           // Load config
    c, err := bootstrap.CreateServerContainer()  // Initialize DI
    bootstrap.StartRestAPI(ctx, cfg, c.Router)   // Start server
}
```

### 3.2 Domain Layer

#### Entities: `domain/entity/user_entity.go`
- **Responsibility**: Represent core business objects
- **Features**:
  - UUID-based user identification
  - Email and hashed password storage
  - Timestamp tracking (created_at, updated_at)
  - Built-in validation using `go-playground/validator`

#### Contracts: `domain/contract/user_repository.go`
- **Responsibility**: Define interfaces for data persistence (ports)
- **Design**: Repository pattern with context support for cancellation
- **Purpose**: Decouple domain logic from database implementation

```go
type UserRepository interface {
    Create(ctx context.Context, u *entity.User) (*entity.User, error)
}
```

#### DTOs: `domain/dto/sign_up_input.go`
- **Responsibility**: Transfer data between layers without exposing entities
- **Purpose**: Prevents tight coupling between API and domain

#### Use Cases: `domain/use_case/auth/sign_up.go`
- **Responsibility**: Orchestrate business logic for user registration
- **Flow**:
  1. Hash password using bcrypt
  2. Create user entity
  3. Validate entity
  4. Persist via repository
  5. Return created user

### 3.3 Infrastructure Layer

#### Repository Implementation: `infrastructure/repository/user_repository.go`
- **Responsibility**: Concrete implementation of `UserRepository` interface
- **Current State**: In-memory storage (generates UUID, normalizes email)
- **Future**: Can be swapped with database implementation (PostgreSQL, MongoDB, etc.)
- **Principle**: Adheres to the Dependency Inversion Principle

#### HTTP Handlers: `infrastructure/http/handlers/auth/handler.go`
- **Responsibility**: Handle HTTP requests and responses
- **Flow**:
  1. Parse and validate HTTP request
  2. Convert API DTO to domain DTO
  3. Execute use case
  4. Return HTTP response

#### Router: `infrastructure/http/router/router.go`
- **Responsibility**: Configure HTTP routes and middleware
- **Middleware**: Request ID, Real IP, Logger, Recoverer (via chi)
- **Routes**: Version-prefixed API endpoints (/api/v1)

### 3.4 Bootstrap Layer

#### Wire Configuration: `bootstrap/wire.go`
- **Tool**: Google Wire for compile-time dependency injection
- **Responsibility**: Define providers and build graph
- **Providers**:
  - `ProvideUserRepository`: Creates repository instance
  - `ProvideSignUpUseCase`: Injects repository into use case
  - `ProvideAuthHandler`: Injects use case into handler
  - `ProvideRouter`: Injects handler into router
  - `ProvideContainer`: Wraps router in container

#### Auto-Generated Code: `bootstrap/wire_gen.go`
- Auto-generated by `wire` tool at build time
- Implements `CreateServerContainer()` function
- Ensures compile-time correctness of dependency graph

#### Server Bootstrap: `bootstrap/server.go`
- **Responsibility**: Start HTTP server and manage lifecycle
- **Features**:
  - Graceful shutdown with 10-second timeout
  - Error channel monitoring
  - Context-based cancellation support

### 3.5 Configuration: `internal/config/config.go`
- **Tool**: `kelseyhightower/envconfig` for environment binding
- **Loads**: Application and database configuration from environment
- **Pattern**: Struct-based configuration with tags and defaults

---

## 4. Data Flow & Dependency Injection

### Request Flow

```
HTTP Request
    ↓
Router (/api/v1/auth/sign-up)
    ↓
AuthHandler.SignUp()
    ↓ Parses & validates request
API Request DTO (sign_up_request.go)
    ↓ Converts to domain DTO
Domain DTO (dto/sign_up_input.go)
    ↓ Executes use case
SignUpUseCase.Execute()
    ↓ Creates entity & hashes password
Domain Entity (entity/user_entity.go)
    ↓ Persists via repository
UserRepository.Create()
    ↓ Returns created user
AuthHandler returns HTTP Response (201 Created)
```

### Dependency Injection Flow

```
main.go
    ↓
config.Load()
    ↓
bootstrap.CreateServerContainer() [Wire-generated]
    ↓
    ├─ UserRepository (concrete impl)
    │   ↓
    ├─ SignUpUseCase (depends on UserRepository)
    │   ↓
    ├─ AuthHandler (depends on SignUpUseCase)
    │   ↓
    └─ Router (depends on AuthHandler)
        ↓
    Container (holds Router)
```

**Benefits of Wire-based DI:**
- Zero runtime overhead (compile-time wiring)
- Type-safe dependency resolution
- Clear visualization of component dependencies
- Compile-time error detection

---

## 5. Technology Stack

### Core Framework
- **Go**: 1.25.4
- **HTTP Router**: `chi/v5` - lightweight, composable router
- **Dependency Injection**: `google/wire` - compile-time wiring
- **Configuration**: `envconfig` - environment-based config

### Utilities
- **Validation**: `go-playground/validator/v10` - struct validation
- **Cryptography**: 
  - `golang.org/x/crypto` - bcrypt password hashing
  - `golang-jwt/jwt/v5` - JWT authentication (imported, not yet used)
- **ID Generation**: `google/uuid` - UUIDs for entities
- **Logging**: `uber/zap` - structured logging
- **Environment Loading**: `joho/godotenv` - .env file support

### Development Tools
- **Code Generation**: Wire for dependency injection
- **Testing**: `testify` (imported as indirect dependency)

---

## 6. Design Patterns Used

| Pattern | Location | Purpose |
|---------|----------|---------|
| **Repository Pattern** | `domain/contract/`, `infrastructure/repository/` | Abstract data persistence |
| **Use Case Pattern** | `domain/use_case/` | Encapsulate business logic |
| **Dependency Injection** | `bootstrap/wire.go` | Manage component dependencies |
| **DTO Pattern** | `domain/dto/`, `internal/api/` | Decouple layers |
| **Entity Pattern** | `domain/entity/` | Core business objects |
| **Factory Pattern** | Constructors (`NewSignUpUseCase`, etc.) | Create instances with dependencies |
| **Middleware Pattern** | `router.go` | Request processing pipeline |

---

## 7. Layer Interactions

### Inbound Flow (Request Processing)

```
HTTP Layer
├─ Chi Router (middleware stack)
│
Infrastructure Layer
├─ HTTP Handler (handler.go)
│  └─ Parses HTTP request
│  └─ Validates input
│  └─ Calls use case
│
Domain/Application Layer
├─ Use Case (sign_up.go)
│  └─ Implements business rules
│  └─ Calls repository
│  └─ Returns domain entity
│
Infrastructure Layer (Persistence)
├─ Repository (user_repository.go)
│  └─ Persists entity
│  └─ Returns result
│
HTTP Handler
└─ Returns HTTP response
```

### Dependency Direction

All dependencies point **inward** toward the domain:
- HTTP layer → Domain layer ✓
- Domain layer → HTTP layer ✗
- Infrastructure (Persistence) → Domain layer ✓
- Domain layer → Infrastructure ✗ (depends on abstractions only)

This ensures domain logic remains **framework-agnostic** and **testable**.

---

## 8. Entry Points & Initialization

### Primary Entry Point: `cmd/server/main.go`

**Execution Steps:**

1. **Configuration Loading**
   ```go
   cfg, err := config.Load()
   ```
   - Loads from `.env` file if present
   - Environment variables override defaults
   - Validates required fields

2. **Signal Handling**
   ```go
   ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
   ```
   - Graceful shutdown on SIGINT (Ctrl+C) or SIGTERM
   - Propagates cancellation through context

3. **Dependency Container Initialization**
   ```go
   c, err := bootstrap.CreateServerContainer()
   ```
   - Wire generates this function from provider declarations
   - Constructs entire dependency graph
   - Returns container with initialized router

4. **HTTP Server Start**
   ```go
   bootstrap.StartRestAPI(ctx, cfg, c.Router)
   ```
   - Starts Chi router on configured port
   - Monitors context for cancellation
   - Gracefully shuts down on signal with 10-second timeout

### HTTP Routes

**Base URL**: `http://localhost:8080`

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/health` | Inline | Health check endpoint |
| POST | `/api/v1/auth/sign-up` | `AuthHandler.SignUp()` | User registration |

**Request Example**:
```json
POST /api/v1/auth/sign-up
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response Example** (201 Created):
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "created_at": "2024-12-02T10:30:00Z",
  "updated_at": null
}
```

---

## 9. Current Capabilities & Future Extensions

### Current Capabilities
✓ User registration with email/password  
✓ Password hashing with bcrypt  
✓ Input validation (email format, password length)  
✓ Entity validation  
✓ Clean architecture foundation  
✓ Dependency injection setup  
✓ Configuration management  

### Extensibility Points

1. **Database Implementation**
   - Swap `UserRepository` implementation (currently in-memory)
   - Add database connection in bootstrap
   - No domain layer changes needed

2. **Authentication**
   - JWT tokens (dependency already imported)
   - Implement authentication middleware
   - Add token generation in use cases

3. **Additional Use Cases**
   - Login, password reset, email verification
   - Follow same pattern: Use Case → Repository → Entity

4. **Error Handling**
   - Create custom error types in domain/error/
   - Domain errors → HTTP status codes mapping
   - Structured error responses

5. **Middleware**
   - Add authentication middleware in router
   - Request/response logging (chi middleware available)
   - CORS, rate limiting, etc.

6. **Testing**
   - Unit test use cases with mock repositories
   - Integration tests with real repository
   - Handler tests with test HTTP utilities

---

## 10. Code Quality Standards (from AGENTS.md)

This codebase follows strict Go conventions:

- **Imports**: Standard library → third-party → internal (blank-line separated)
- **Naming**: 
  - Files: `snake_case`
  - Functions: `CamelCase`, exported = capitalized
  - Variables: `camelCase`
  - Constants: `ALL_CAPS`
- **Error Handling**: Always check errors immediately, use `fmt.Errorf()` for context
- **Formatting**: `go fmt` and `go vet` (make targets available)
- **Line Length**: ~80-100 characters (Go standard)
- **Interfaces**: Used for dependency injection and abstraction
- **Constructors**: Return pointers when modifying state

---

## 11. Summary

This Go application demonstrates **production-grade Clean Architecture**:

- **Layers are independent** - easy to test and modify
- **Dependencies flow inward** - domain logic is framework-agnostic
- **Wire-based DI** - zero-runtime-overhead, compile-time safety
- **Clear separation** - each component has a single responsibility
- **Extensible** - new features follow established patterns
- **Type-safe** - leverages Go's strong typing and interfaces

The architecture successfully balances:
- **Simplicity**: Small codebase, easy to understand
- **Extensibility**: Clear patterns for adding features
- **Maintainability**: Independent layers with minimal coupling
- **Testability**: Each layer can be tested in isolation

Perfect foundation for a scalable backend application.
