# E-commerce Backend with Clean Architecture and DDD

This project is an e-commerce backend API built with Fiber, GORM, and follows clean architecture principles with Domain-Driven Design (DDD) practices.

## Architecture

The project follows a clean architecture pattern with domain-driven design principles separated into distinct layers:

### Domain Layer

The core of the application containing business logic and rules:

- **Entities**: Domain objects representing business concepts (Product, Order, Customer, etc.)
- **Value Objects**: Immutable objects with no identity but with meaning (Money, Address, etc.)
- **Domain Services**: Logic that doesn't naturally belong to entities or value objects
- **Repository Interfaces**: Contracts for data access

### Use Case Layer (Application Layer)

Orchestrates the flow of data between the domain and infrastructure layers:

- **Use Cases**: Application-specific business rules
- **DTOs**: Data Transfer Objects for input/output
- **Services**: Orchestrates domain operations

### Infrastructure Layer

Implements the interfaces defined in the domain layer:

- **Database**: GORM implementations of repositories
- **External Services**: Third-party integrations
- **Messaging**: Implementation of message queues or event buses

### Interface Layer

Handles HTTP requests and responses:

- **Controllers**: HTTP handlers mapping requests to use cases
- **Routes**: API endpoint definitions
- **Middleware**: Request processing middleware
- **Presenters**: Format data for presentation

## DDD Implementation

### Entity vs Value Object

- **Entities**: Objects with identity (ID) that can change over time (User, Product, Order)
- **Value Objects**: Immutable objects defined by their attributes (Money, Address)

### Aggregates and Aggregate Roots

- **Aggregate**: Cluster of domain objects treated as a single unit (Order with OrderItems)
- **Aggregate Root**: Entry point to the aggregate, maintains invariants (Order is the root of Order items)

### Repository Pattern

Provides collection-like interfaces for accessing domain objects:

- Repository interfaces defined in the domain layer
- Implementations in the infrastructure layer

### Domain Events

Notifies other parts of the system when something important happens in the domain:

- Defined in the domain layer
- Published through event bus/message broker

## Technology Stack

- **Go**: Programming language
- **Fiber**: Web framework
- **GORM**: ORM for database operations
- **MySQL**: Relational database

## Project Structure

```
├── cmd
│   └── api
│       └── main.go           # Entry point
├── domain
│   ├── common                # Shared domain logic
│   │   └── vo                # Value objects
│   ├── user                  # User domain
│   ├── product               # Product domain
│   ├── inventory             # Inventory domain
│   └── order                 # Order domain
├── usecase                   # Application use cases
│   ├── user
│   ├── product
│   ├── inventory
│   └── order
├── infrastructure
│   ├── repository            # Repository implementations
│   ├── persistence           # Database connection
│   ├── auth                  # Authentication
│   └── services              # External services
└── interface
    ├── api                   # API handlers
    ├── middleware            # HTTP middleware
    └── dto                   # Request/response objects
```

## Implementation Guide

### Domain Layer

1. Define entities with behavior and validation
2. Create value objects for immutable concepts
3. Define repository interfaces
4. Implement domain services

### Use Case Layer

1. Create use cases for each business operation
2. Define DTOs for input/output
3. Implement application services that orchestrate domain operations

### Infrastructure Layer

1. Implement repositories using GORM
2. Set up database connection
3. Integrate external services

### Interface Layer

1. Define API routes
2. Implement controllers
3. Create middleware
4. Set up dependency injection

## Getting Started

### Prerequisites

- Go 1.16 or higher
- MySQL

### Running the Application

1. Clone the repository
2. Set up environment variables
3. Run database migrations
4. Start the server

```bash
go run cmd/api/main.go
```

## Development Notes

- Follow domain-driven design principles
- Use dependency injection for loose coupling
- Write unit tests for domain logic
- Write integration tests for repositories
- Use continuous integration for quality assurance
