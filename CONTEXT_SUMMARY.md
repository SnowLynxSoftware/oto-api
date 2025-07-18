# Open Trivia Online API - Context Summary

## Table of Contents

1. [Project Overview & Purpose](#project-overview--purpose)
2. [Architecture & Design Patterns](#architecture--design-patterns)
3. [Technology Stack & Dependencies](#technology-stack--dependencies)
4. [Development Environment & Setup](#development-environment--setup)
5. [Authentication & Authorization](#authentication--authorization)
6. [Database Design & Relationships](#database-design--relationships)
7. [API Design & Communication Patterns](#api-design--communication-patterns)
8. [Code Organization & File Structure](#code-organization--file-structure)
9. [Error Handling & Logging](#error-handling--logging)
10. [Testing Approach & Standards](#testing-approach--standards)
11. [Build & Deployment Process](#build--deployment-process)
12. [Integration Points with Frontend](#integration-points-with-frontend)
13. [Key Architectural Decisions](#key-architectural-decisions)
14. [Common Patterns & Best Practices](#common-patterns--best-practices)
15. [Future Development Context](#future-development-context)

---

## Project Overview & Purpose

**Open Trivia Online API** is the backend service for a multiplayer trivia game platform built in Go. The API serves as the central hub for:

- **User Management**: Registration, authentication, profile management, and role-based access control
- **Trivia Content Management**: Questions, answers, decks, and categories with full CRUD operations
- **Game Session Management**: Future multiplayer game logic and session handling
- **Administrative Operations**: Content moderation, user management, and system analytics
- **Waitlist Management**: Pre-launch user registration and notification system

The API is designed to scale horizontally and support real-time multiplayer gaming while maintaining data consistency and security.

---

## Architecture & Design Patterns

### Clean Architecture Implementation

The project follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │───▶│    Services     │───▶│  Repositories   │
│  (HTTP Layer)   │    │ (Business Logic)│    │  (Data Access)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │     Models      │    │    Database     │
│ (Cross-cutting) │    │   (Entities)    │    │   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Design Patterns Used

- **Repository Pattern**: Data access abstraction with interfaces for testability
- **Dependency Injection**: Constructor-based DI for loose coupling
- **Interface Segregation**: Small, focused interfaces for each domain
- **Command Pattern**: CLI commands for server start and migrations
- **Middleware Pattern**: Cross-cutting concerns like authentication and CORS
- **Factory Pattern**: Service and repository initialization

### Layer Responsibilities

- **Controllers**: HTTP request/response handling, input validation, authorization checks
- **Services**: Business logic, data transformation, orchestration between repositories
- **Repositories**: Data persistence, query optimization, database interaction
- **Middleware**: Authentication, logging, CORS, request preprocessing
- **Models**: Data structures, DTOs, and entity definitions

---

## Technology Stack & Dependencies

### Core Technologies

- **Go 1.24.0**: Primary language with modern Go features
- **PostgreSQL**: Primary database with ACID compliance
- **Chi Router v5**: HTTP routing and middleware
- **sqlx**: SQL extensions for Go with prepared statements
- **JWT (golang-jwt/jwt/v5)**: Token-based authentication
- **bcrypt (golang.org/x/crypto)**: Password hashing
- **Zerolog**: Structured logging with performance focus

### Key Dependencies

```go
// Core Framework
github.com/go-chi/chi/v5        // HTTP router and middleware
github.com/jmoiron/sqlx         // SQL extensions and query builder
github.com/lib/pq               // PostgreSQL driver

// Security & Authentication
github.com/golang-jwt/jwt/v5    // JWT token handling
golang.org/x/crypto             // Password hashing (bcrypt)

// External Services
github.com/sendgrid/sendgrid-go // Email notifications
github.com/morpheuszero/go-heimdall // HTTP client with retry logic

// Development & Utilities
github.com/joho/godotenv        // Environment variable loading
github.com/rs/zerolog           // Structured logging
```

### Development Tools

- **Go Modules**: Dependency management
- **Docker & Docker Compose**: Containerization and local development
- **Database Migrations**: SQL-based with versioning
- **Environment Configuration**: `.env` file with validation

---

## Development Environment & Setup

### Local Development Setup

1. **Prerequisites**: Go 1.24+, PostgreSQL 13+, Docker (optional)
2. **Configuration**: Environment variables via `.env` file
3. **Database**: Local PostgreSQL or Docker container
4. **Commands**:
   - `go run .` - Start server (default command)
   - `go run . migrate` - Run database migrations
   - `go mod tidy` - Update dependencies

### Environment Configuration

**Required Environment Variables**:

```bash
CLOUD_ENV=local                    # Environment identifier
DB_CONNECTION_STRING=postgres://   # PostgreSQL connection
AUTH_HASH_PEPPER=secret           # Password hashing salt
JWT_SECRET_KEY=secret             # JWT signing key
SENDGRID_API_KEY=key              # Email service API key
```

**Optional Configuration**:

```bash
DEBUG_MODE=true                   # Enable debug logging
CORS_ALLOWED_ORIGIN=http://...    # Frontend URL for CORS
COOKIE_DOMAIN=localhost           # Cookie domain for auth
```

### Command Line Interface

The application uses a command pattern for different operations:

- **Default/Server**: `go run .` or `go run . server`
- **Migrations**: `go run . migrate`
- **Extensible**: Easy to add new commands via the handler pattern

---

## Authentication & Authorization

### Multi-Tier Authentication System

#### User Roles & Permissions

```go
// User Types (Hierarchical)
UserTypeAdmin   = "admin"     // Full system access
UserTypeSupport = "support"   // User management access
UserTypePlayer  = "player"    // Standard user access
```

#### Authentication Flow

1. **Registration**: Email/password with verification required
2. **Login**: Basic auth header → JWT token in HTTP-only cookie
3. **Token Validation**: Middleware validates JWT on protected routes
4. **User Context**: Enriched user data for request authorization

#### Authorization Middleware

```go
type AuthorizedUserContext struct {
    Id        int    `json:"id"`
    Email     string `json:"email"`
    Username  string `json:"username"`
    IsAdmin   bool   `json:"is_admin,omitempty"`
    IsSupport bool   `json:"is_support,omitempty"`
}
```

#### Security Features

- **JWT Tokens**: Stateless authentication with expiration
- **HTTP-Only Cookies**: XSS protection for token storage
- **Password Hashing**: bcrypt with pepper for enhanced security
- **Role-Based Access**: Granular permissions per endpoint
- **Account Security**: Email verification, ban/unban functionality

### Email Verification System

- **Verification Flow**: Register → Email sent → Token verification → Account activation
- **Email Service**: SendGrid integration with templated emails
- **Token Management**: Secure verification tokens with expiration

---

## Database Design & Relationships

### Core Entities

#### Users Table

```sql
- id (Primary Key)
- email (Unique, Required)
- display_name (Required)
- user_type_key (admin/support/player)
- is_verified, is_archived, is_banned
- password_hash, created_at, modified_at
- profile fields (avatar_url, profile_text)
- ban_reason, last_login
```

#### Trivia System

```sql
# Trivia Decks
- id, name, description
- is_approved, is_archived, is_system_deck
- created_at

# Trivia Questions
- id, question, correct_answer
- tags (JSON array), is_published, is_archived
- created_at

# Wrong Answer Pool
- id, answer_text
- tags (JSON array), is_archived
- created_at
```

#### Waitlist Table

```sql
- id, email, created_at
- Simple pre-launch email collection
```

### Database Migration Strategy

- **Version-Based**: Sequential numbered migration files
- **SQL-Based**: Pure SQL for clarity and database optimization
- **Forward-Only**: No rollback migrations (by design choice)
- **Environment Agnostic**: Same migrations across all environments

### Data Relationships

- **Users**: Self-contained with role-based typing
- **Trivia Content**: Loosely coupled, tag-based categorization
- **Future Scalability**: Designed for game sessions and user statistics

---

## API Design & Communication Patterns

### RESTful API Design

#### Endpoint Structure

```
/health              - System health check
/auth/*              - Authentication endpoints
/users/*             - User management (admin/support only)
/trivia/*            - Trivia content management
/waitlist/*          - Waitlist management
```

#### HTTP Patterns

- **GET**: Data retrieval with pagination and filtering
- **POST**: Resource creation and authentication actions
- **PUT**: Full resource updates
- **PATCH**: Partial updates (toggle operations)
- **Consistent Responses**: JSON with standardized error messages

### Request/Response Patterns

#### Pagination Standard

```go
type PaginatedResponse struct {
    PageSize int   `json:"page_size"`
    Page     int   `json:"page"`
    Total    int   `json:"total"`
    Results  []any `json:"results"`
}
```

#### Error Handling

```go
// Consistent HTTP status codes
400 - Bad Request (validation errors)
401 - Unauthorized (auth required/failed)
403 - Forbidden (insufficient permissions)
404 - Not Found
500 - Internal Server Error
```

#### Query Parameters

- **Pagination**: `page`, `page_size` (default: 25)
- **Search**: `search` (text-based filtering)
- **Filters**: `status`, `user_type`, `tags`
- **Consistent naming**: snake_case for JSON, camelCase avoided

### CORS Configuration

- **Allowed Origins**: Configurable frontend URL
- **Methods**: GET, POST, PUT, PATCH, DELETE, OPTIONS
- **Credentials**: Enabled for cookie-based authentication
- **Headers**: Standard content types and authorization

---

## Code Organization & File Structure

### Project Structure

```
oto-api/
├── main.go                     # Application entry point
├── go.mod, go.sum             # Dependency management
├── Dockerfile                 # Container configuration
├── docker-compose.yml         # Local development setup
├── cmd/                       # CLI commands and handlers
│   ├── handler.go             # Command dispatcher
│   ├── server.go              # Server start command
│   └── migrate.go             # Database migration command
├── config/                    # Configuration management
│   └── app-config.go          # Environment variable handling
├── server/                    # Core application logic
│   ├── app.server.go          # Server setup and DI container
│   ├── controllers/           # HTTP request handlers
│   ├── services/              # Business logic layer
│   ├── database/              # Data access layer
│   │   ├── database.go        # Connection management
│   │   └── repositories/      # Data access implementations
│   ├── middleware/            # Cross-cutting concerns
│   ├── models/                # Data structures and DTOs
│   └── util/                  # Shared utilities
└── migrations/                # Database schema versions
    ├── 20250401_v1_initial.sql
    ├── 20250408_v1_waitlist.sql
    └── 20250409_v1_trivia_entries.sql
```

### Naming Conventions

- **Files**: snake_case.go (e.g., `auth.controller.go`)
- **Packages**: lowercase, descriptive (controllers, services, repositories)
- **Interfaces**: IServiceName pattern (e.g., `IUserService`)
- **Structs**: PascalCase (e.g., `UserController`)
- **Functions**: camelCase (private) and PascalCase (public)

### Import Organization

- **Standard Library**: First group
- **Third-Party**: Second group
- **Local Packages**: Third group
- **Consistent Ordering**: Alphabetical within groups

---

## Error Handling & Logging

### Logging Strategy

#### Structured Logging with Zerolog

```go
// Log Levels
util.LogError(err)              // Error conditions
util.LogWarning(message)        // Warning conditions
util.LogInfo(message)           // General information
util.LogDebug(message)          // Debug information

// Error Logging with Stack Traces
util.LogErrorWithStackTrace(err) // Full error context
```

#### Logging Configuration

- **Development**: Debug level, console output
- **Production**: Info level, JSON format
- **Performance**: Optimized for high-throughput logging

### Error Handling Patterns

#### Controller Error Handling

```go
// Consistent pattern across all controllers
if err != nil {
    util.LogErrorWithStackTrace(err)
    http.Error(w, "user-friendly message", http.StatusCode)
    return
}
```

#### Service Layer Error Propagation

- **Errors bubble up**: Services return detailed errors to controllers
- **Context preservation**: Full error context maintained through layers
- **User-friendly messages**: Controllers translate technical errors

#### Database Error Handling

- **Connection failures**: Panic on startup (fail-fast)
- **Query errors**: Detailed logging with SQL context
- **Transaction handling**: Proper rollback on failures

---

## Testing Approach & Standards

### Testing Strategy

#### Unit Testing Coverage

- **Services**: Business logic testing with mocked repositories
- **Repositories**: Database interaction testing
- **Utilities**: Helper function validation

#### Mock Implementation Pattern

```go
// Mock repositories for service testing
type MockUserRepository struct{}
func (m *MockUserRepository) GetUserById(id int) (*UserEntity, error) {
    // Mock implementation
}
```

#### Test Organization

- **Parallel Structure**: Tests mirror production code structure
- **Naming Convention**: `service_test.go`, `repository_test.go`
- **Table-Driven Tests**: For multiple test cases
- **Test Fixtures**: Reusable test data and helpers

### Testing Commands

```bash
go test ./... -v              # Run all tests with verbose output
go test ./... -coverprofile=coverage.out  # Generate coverage report
```

---

## Build & Deployment Process

### Local Development

```bash
# Development workflow
go mod tidy                   # Update dependencies
go run . migrate              # Update database schema
go run .                      # Start development server
```

### Production Build

```bash
# Build process
go build -o oto-api          # Compile binary
./oto-api migrate            # Run migrations
./oto-api                    # Start server
```

### Docker Configuration

```dockerfile
# Multi-stage build for optimized images
FROM golang:1.24-alpine AS builder
# Build stage...

FROM alpine:latest
# Runtime stage with minimal footprint
```

### Environment Management

- **Local**: `.env` file for development
- **Production**: Environment variables via deployment platform
- **Validation**: Required variables checked on startup

---

## Integration Points with Frontend

### API Communication Standards

#### Authentication Flow

1. **Frontend Login**: POST to `/auth/login` with Basic auth
2. **Token Storage**: JWT in HTTP-only cookie (automatic)
3. **Request Authentication**: Cookies sent automatically
4. **Token Validation**: `/auth/token` endpoint for user context

#### Data Exchange Patterns

```javascript
// Frontend request pattern
const response = await fetch("/api/endpoint", {
  method: "POST",
  credentials: "include", // Send cookies
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify(data),
});
```

#### CORS Configuration

- **Frontend URL**: Configurable allowed origin
- **Credentials**: Enabled for cookie-based auth
- **Development**: Typically `http://localhost:4200` (Angular dev server)

### Real-time Communication (Future)

- **WebSocket Support**: Planned for multiplayer game sessions
- **Event Broadcasting**: User actions, game state changes
- **Connection Management**: Authentication over WebSocket

### Admin Panel Integration

- **Role-Based Endpoints**: Admin/support role requirements
- **User Management**: Full CRUD operations via API
- **Trivia Management**: Content creation and moderation
- **Analytics Data**: Dashboard statistics and metrics

---

## Key Architectural Decisions

### Decision Log

#### Database Choice: PostgreSQL

**Decision**: Use PostgreSQL as primary database
**Rationale**:

- ACID compliance for financial/scoring data
- JSON support for flexible schemas (tags)
- Excellent Go driver support (lib/pq)
- Horizontal scaling capabilities

**Date**: Project inception
**Status**: Implemented

#### Authentication: JWT + HTTP-Only Cookies

**Decision**: JWT tokens stored in HTTP-only cookies
**Rationale**:

- XSS protection (HTTP-only prevents JavaScript access)
- CSRF protection with SameSite cookie settings
- Stateless authentication for horizontal scaling
- Automatic cookie handling by browsers

**Date**: Authentication implementation
**Status**: Implemented

#### Architecture: Clean Architecture

**Decision**: Implement clean architecture with repository pattern
**Rationale**:

- Testability through dependency injection
- Separation of concerns for maintainability
- Database abstraction for future migrations
- Business logic isolation from infrastructure

**Date**: Project structure design
**Status**: Implemented

#### CLI Pattern: Command Handler

**Decision**: Use command pattern for different application modes
**Rationale**:

- Single binary with multiple functions
- Easy addition of new commands (migrations, workers, etc.)
- Clear separation of concerns
- Simplified deployment

**Date**: CLI design phase
**Status**: Implemented

#### Error Handling: Structured Logging

**Decision**: Zerolog for structured, high-performance logging
**Rationale**:

- JSON output for log aggregation
- Zero allocation logging for performance
- Structured context for debugging
- Production-ready observability

**Date**: Logging implementation
**Status**: Implemented

---

## Common Patterns & Best Practices

### Controller Patterns

#### Standard Controller Structure

```go
type ControllerName struct {
    service        services.IServiceName
    authMiddleware middleware.IAuthMiddleware
}

func NewController(service services.IServiceName, auth middleware.IAuthMiddleware) *ControllerName {
    return &ControllerName{service: service, authMiddleware: auth}
}

func (c *ControllerName) MapController() *chi.Mux {
    r := chi.NewRouter()
    r.Get("/", c.getItems)
    r.Post("/", c.createItem)
    return r
}
```

#### Error Handling Pattern

```go
func (c *Controller) endpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Authorization check
    userContext, err := c.authMiddleware.Authorize(r, []string{"admin"})
    if err != nil {
        util.LogErrorWithStackTrace(err)
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    // 2. Input validation
    var dto SomeDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    // 3. Business logic
    result, err := c.service.DoSomething(&dto)
    if err != nil {
        util.LogErrorWithStackTrace(err)
        http.Error(w, "operation failed", http.StatusInternalServerError)
        return
    }

    // 4. Response formatting
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

### Service Patterns

#### Service Interface Design

```go
type IServiceName interface {
    GetItems(pageSize, offset int, filters string) (*models.PaginatedResponse, error)
    GetItemById(id int64) (*repositories.EntityName, error)
    CreateItem(dto *models.CreateDTO) (*repositories.EntityName, error)
    UpdateItem(dto *models.UpdateDTO, id int64) (*repositories.EntityName, error)
    ToggleItemStatus(id int64) error
}
```

### Repository Patterns

#### Repository Interface Design

```go
type IRepositoryName interface {
    GetItemsCount(filters string) (*int, error)
    GetItems(pageSize, offset int, filters string) ([]*EntityName, error)
    GetItemById(id int64) (*EntityName, error)
    CreateItem(dto *models.CreateDTO) (*EntityName, error)
    UpdateItem(dto *models.UpdateDTO, id int64) (*EntityName, error)
    ToggleItemArchived(id int64) error
}
```

### Database Interaction Patterns

#### Query Execution Pattern

```go
func (r *Repository) GetItems(pageSize, offset int, search string) ([]*Entity, error) {
    var entities []*Entity

    sql := `SELECT id, name, created_at FROM table_name WHERE 1=1`
    args := []interface{}{pageSize, offset}

    if search != "" {
        sql += ` AND name ILIKE $3`
        args = append(args, "%"+search+"%")
    }

    sql += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`

    err := r.db.DB.Select(&entities, sql, args...)
    if err != nil {
        return nil, err
    }

    return entities, nil
}
```

---

## Future Development Context

### Planned Features & Considerations

#### Multiplayer Game Sessions

- **WebSocket Integration**: Real-time communication for live games
- **Game State Management**: Redis for session state and caching
- **Matchmaking Service**: Queue management and game creation
- **Scoring System**: Real-time score calculation and leaderboards

#### Performance Optimization

- **Database Optimization**: Query optimization, indexing strategy
- **Caching Layer**: Redis for frequently accessed data
- **Connection Pooling**: Database connection optimization
- **API Rate Limiting**: Protection against abuse

#### Monitoring & Observability

- **Metrics Collection**: Prometheus integration
- **Health Checks**: Comprehensive system health endpoints
- **Distributed Tracing**: Request tracing across services
- **Error Aggregation**: Centralized error tracking

#### Security Enhancements

- **API Rate Limiting**: Request throttling and abuse prevention
- **Input Validation**: Enhanced validation library integration
- **Audit Logging**: User action tracking for compliance
- **Security Headers**: Enhanced HTTP security headers

### Extension Points

#### Adding New Entities

1. Create migration file in `/migrations/`
2. Add entity struct in `/server/models/`
3. Implement repository interface and struct
4. Create service interface and implementation
5. Add controller with route mapping
6. Update server initialization in `app.server.go`

#### Adding New Commands

1. Create command struct implementing `Command` interface
2. Add to command map in `cmd/handler.go`
3. Implement `Execute()` method with business logic

#### Adding New Middleware

1. Create middleware struct in `/server/middleware/`
2. Implement middleware function signature
3. Add to router setup in `app.server.go`

### Migration Strategy for Future Changes

#### Database Schema Changes

- **Additive Changes**: Safe migrations that don't break existing code
- **Breaking Changes**: Coordinated with frontend updates
- **Data Migration**: Separate data migration scripts when needed

#### API Versioning Strategy

- **Current**: Single version, backward compatible changes
- **Future**: URL-based versioning (`/v1/`, `/v2/`) when breaking changes needed
- **Deprecation**: Gradual deprecation with proper client notification

This context summary provides comprehensive coverage of the Open Trivia Online API architecture, patterns, and implementation decisions. It serves as a reference for future development sessions and onboarding new team members to the project.
