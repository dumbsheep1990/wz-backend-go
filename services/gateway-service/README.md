# Gateway Service

The Gateway Service acts as the central entry point for all client requests in the 万知 backend system. It routes requests to the appropriate microservices, enforces authentication, performs rate limiting, and provides service discovery and health monitoring.

## Architecture

The Gateway Service follows a Domain-Driven Design (DDD) architecture with the following layers:

### Domain Layer
- **Entities**: `Route` and `Service` represent the core domain concepts
- **Value Objects**: Strongly typed values like `RouteID`, `Path`, and `ServiceName`
- **Domain Services**: `GatewayDomainService` implements domain logic like route matching and service health checking

### Application Layer
- **Application Services**: `GatewayApplicationService` orchestrates domain services and repositories
- **DTOs**: Data Transfer Objects for input/output across the application boundary

### Infrastructure Layer
- **Persistence**:
  - MySQL repositories for persistent storage
  - Redis repositories for caching and rate limiting
  - Repository factories and facades to combine storage with caching

### Interface Layer
- **HTTP Controllers**: Handle API requests and transform between HTTP and application layer
- **Middleware**: Process requests for routing, authentication, rate limiting, and proxying

## Repository Implementation

### Repository Factory Pattern

The Gateway Service uses a repository factory pattern to manage the creation and lifecycle of different repository implementations:

```go
repoFactory := gateway.NewRepositoryFactory(dbConn, redisClient)
routeRepo := repoFactory.GetRouteRepository()
serviceRepo := repoFactory.GetServiceRepository()
rateLimiterRepo := repoFactory.GetRateLimiterRepository()
```

This provides a single entry point for obtaining repository instances and hides the complexity of choosing between MySQL or Redis implementations.

### Repository Facades

Repository facades combine MySQL and Redis repositories to provide caching with persistent storage:

1. **RouteRepositoryFacade**:
   - First tries to read from Redis cache
   - Falls back to MySQL if cache miss
   - Automatically updates cache when data changes
   - Handles cache invalidation

2. **ServiceRepositoryFacade**:
   - Similar cache-first approach for services
   - Special handling for service health status updates
   - Supports finding active and healthy services efficiently

### Redis Rate Limiter

The Redis-based rate limiter uses Lua scripts for atomic operations to enforce request limits based on:
- IP address
- API token
- Service name

## API Endpoints

### Management API

- **Services**:
  - `POST /api/v1/management/services` - Register a new service
  - `GET /api/v1/management/services` - List all services
  - `GET /api/v1/management/services/{name}` - Get service details
  - `PUT /api/v1/management/services/{name}` - Update a service
  - `DELETE /api/v1/management/services/{name}` - Delete a service
  - `GET /api/v1/management/services/{name}/health` - Check service health
  - `POST /api/v1/management/services/{name}/routes` - Add route to service
  - `GET /api/v1/management/services/{name}/routes` - Get service routes
  - `POST /api/v1/management/services/wz-categories` - Create 万知 category routes

- **Routes**:
  - `POST /api/v1/management/routes` - Register a new route
  - `GET /api/v1/management/routes` - List all routes
  - `GET /api/v1/management/routes/{id}` - Get route details
  - `PUT /api/v1/management/routes/{id}` - Update a route
  - `DELETE /api/v1/management/routes/{id}` - Delete a route

### Gateway Functionality

The gateway itself handles all other routes by:
1. Finding the matching route definition
2. Applying authentication if required
3. Enforcing rate limits
4. Proxying the request to the target service

## Configuration

Configuration is loaded from `configs/gateway.yaml` and includes:

- Server settings (port, environment)
- Service definitions
- Database connection settings
- Redis connection settings
- Telemetry settings

## Getting Started

### Prerequisites

- Go 1.16+
- MySQL 5.7+
- Redis 6.0+

### Running the Service

1. Make sure MySQL and Redis are running
2. Update `configs/gateway.yaml` with proper connection settings
3. Build and run the service:

```bash
cd services/gateway-service
go build
./gateway-service -f configs/gateway.yaml
```

## Special Features for 万知 Website

The Gateway Service supports automatically creating routes for the 21 "入同" categories of the 万知 website:

1. 同用
2. 同好
3. 同购
4. 同年
5. 同游
6. 同在
7. 同市
8. 同企
9. 同亲
10. 同班
11. 同师
12. 同业
13. 同网
14. 同工
15. 同务
16. 同艺
17. 同玩
18. 同闲
19. 同拍
20. 同乡
21. 同学

These routes can be automatically created using the `POST /api/v1/management/services/wz-categories` endpoint.
