# DrGo Backend API

A Go-based REST API using the Gin web framework for the DrGo medical appointment platform.

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. Navigate to the backend directory:
```bash
cd backend
```

2. Download dependencies:
```bash
go mod download
```

## Running the Server

```bash
go run cmd/api/main.go
```
The API will be available at `http://localhost:8080`


## Testing the API

### Using PowerShell (Windows):

```powershell
# Test health check
(Invoke-WebRequest -Uri http://localhost:8080/health -UseBasicParsing).Content
```

### Using curl:

```bash
# Test health check
curl http://localhost:8080/health
```

## Building for Production

Build the binary:
```bash
go build -o drgo-api cmd/api/main.go
```

Run the binary:
```bash
./drgo-api
```

## Project Structure

This project follows a clean, modular architecture with clear separation of concerns for a medical appointment platform (DrGo). The structure is designed to support a 3-domain architecture: public website, doctor dashboard, and patient dashboard.

### Directory Organization

#### 1. **Entry Points** (`cmd/`)
Contains the main application entry points:

- **`cmd/api/main.go`** - Main API server entry point
  - Initializes database connection with environment-based configuration
  - Sets up HTTP server with graceful shutdown
  - Configures CORS and error handling middleware
  - Starts the Gin web server on configurable port (default: 8080)

- **`cmd/seed/main.go`** - Database seeding utility
  - Seeds the database with sample data for development
  - Creates fake doctor profiles with various specialties
  - Uses `gofakeit` library for generating realistic test data
  - Currently seeds 50 doctors with random specialties and details

#### 2. **Internal Package** (`internal/`)
All private application code that cannot be imported by external projects.

##### **Domain Layer** (`internal/domain/`)
Contains business entities and interfaces:

- **`doctor.go`** - Base interface for all domain entities
  - Defines `ModelEnt` interface with `GetId()` methods
  - Ensures consistent entity behavior across the application

- **`medical/`** - Medical domain entities
  - **`doctor_entity.go`** - Doctor entity definition
  - **`specialty_entity.go`** - Medical specialty entity definition

##### **Repository Layer** (`internal/repository/`)
Data access layer implementing repository pattern:

- **`medical/doctor_repository.go`** - Doctor data access
  - Implements `DoctorRepository` interface
  - Uses `go-sqlbuilder` for SQL query generation
  - Supports filtering, pagination, and complex queries
  - Handles database connection and error management


##### **API Layer** (`internal/api/`)
HTTP handlers organized by domain and functionality:

- **`patient-panel/medical/`** - Doctor listing for patients
  - **`doctor_handler.go`** - HTTP handlers for doctor search and listing
  - **`doctor_dto.go`** - Data Transfer Objects for API requests/responses
  - Supports search by specialty, name, and pagination
  - Returns paginated results with metadata

##### **Middleware** (`internal/middleware/`)
HTTP middleware components:

- **`error_handler.go`** - Centralized error handling
  - Handles validation errors with detailed field-level messages
  - Provides consistent error response format
  - Supports different error types (validation, internal server errors)
  - Custom error messages for different validation rules

##### **Router** (`internal/router/`)
Route configuration and setup:

- **`router.go`** - Main router setup
  - Configures Gin router with middleware
  - Sets up API routes with versioning (`/api/`)
  - Includes health check endpoints
  - Organizes routes by domain (public, doctor, patient)

- **`patient-panel/router.go`** - Patient panel routes
  - Sets up patient-specific routes
  - Initializes repositories and handlers
  - Registers doctor-related endpoints for patient access

##### **Database** (`internal/database/`)
Database connection and migration management:

- **`connection.go`** - Database connection management
  - PostgreSQL connection with connection pooling
  - Environment-based configuration
  - Connection health checks and graceful shutdown
  - Configurable connection limits and timeouts

- **`migrations/`** - Database schema migrations
  - SQL-based migrations with proper indexing
  - Supports incremental schema updates

##### **Pagination** (`internal/pagination/`)
Pagination utilities for API responses:

- **`paginator.go`** - Base pagination interface
- **`offset.go`** - Offset-based pagination implementation
- Supports different pagination strategies based on use case

##### **Filter** (`internal/filter/`)
Query filtering system for dynamic database queries:

- **`filter.go`** - Base filter interface and composite filters
- **`medical/`** - Medical domain filters
  - **`doctor_filter.go`** - Doctor-specific filters (search, specialty)
  - **`specialty_filter.go`** - Specialty-specific filters
- Enables dynamic query building with multiple filter conditions

##### **Query Builder** (`internal/query_builder/`)
SQL query construction utilities:

- **`query_builder.go`** - Base query builder interface
- **`medical/`** - Medical domain query builders
  - **`doctor_query_builder.go`** - Doctor-specific query building
- Supports filtering, ordering, and pagination

##### **Documentation** (`docs/`)
Project documentation and planning:

- **`ai-talks/`** - AI-generated documentation and planning
  - **`api/`** - API documentation (doctor, patient, public APIs)
  - **`er/`** - Entity relationship diagrams
  - **`mvp/`** - MVP planning documents
  - **`product/`** - Product description and requirements

#### 4. **Configuration Files**
- **`go.mod`** - Go module dependencies
  - Gin web framework
  - PostgreSQL driver (`lib/pq`)
  - JWT authentication
  - SQL query builder (`go-sqlbuilder`)
  - Validation library
  - UUID generation
  - Environment variable management

- **`go.sum`** - Dependency checksums for reproducible builds

- **`run.ps1`** - Windows PowerShell script for running the application

### Architecture Principles

#### **Dependency Flow:**
```
Handler → Service → Repository → Database
              ↓
          Provider (SMS, AI, Payment)
```

#### **Key Design Patterns:**
- **Repository Pattern** - Clean separation between business logic and data access
- **Interface Segregation** - Small, focused interfaces for each component
- **Dependency Injection** - Loose coupling between layers
- **Query Builder Pattern** - Dynamic SQL generation with type safety

#### **Database Strategy:**
- **No ORM** - Direct SQL queries using `go-sqlbuilder`
- **SQL Injection Protection** - Parameterized queries only
- **Connection Pooling** - Efficient database connection management
- **Migration System** - Version-controlled schema changes

#### **API Design:**
- **RESTful Endpoints** - Standard HTTP methods and status codes
- **Pagination Support** - Both cursor and offset-based pagination
- **Filtering & Search** - Dynamic query building with multiple filters
- **Error Handling** - Consistent error response format
- **Validation** - Input validation with detailed error messages

#### **Security Considerations:**
- **JWT Authentication** - Token-based authentication system
- **Input Validation** - Comprehensive request validation
- **SQL Injection Prevention** - Parameterized queries only
- **CORS Configuration** - Cross-origin request handling

### File Naming Conventions

- **Repositories:** `{entity}_repository.go`
- **Handlers:** `{entity}_handler.go`
- **DTOs:** `{entity}_dto.go`
- **Entities:** `{entity}.go` or `{entity}_entity.go`
- **Filters:** `{entity}_filter.go`
- **Query Builders:** `{entity}_query_builder.go`
- **Migrations:** `{number}_{description}.sql`
- **Routers:** `router.go` (in domain-specific directories)
