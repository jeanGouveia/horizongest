# Redis Architecture

**Version:** 1.0  
**Date:** 31/07/2026  
**Status:** ✅ **IMPLEMENTED**

---

## Overview

Redis is integrated into the HorizonGest architecture as a high-performance key-value store supporting multiple use cases: caching, idempotency, distributed locks, rate limiting, and session management. This document describes the Redis infrastructure implementation, design decisions, and integration patterns.

---

## Architecture Principles

### Clean Architecture
- Redis infrastructure is isolated in `internal/infra/redis` package
- Interfaces are defined at the infrastructure layer
- Domain layer remains independent of Redis implementation
- Dependency inversion through interfaces

### DDD Alignment
- Redis is used as an infrastructure concern, not for domain entities
- PostgreSQL remains the primary database for persistent domain data
- Redis supports cross-cutting concerns (caching, locks, sessions)

### SOLID Principles
- **Single Responsibility:** Each Redis component has a single purpose (Cache, Lock, RateLimiter, Session)
- **Open/Closed:** Interfaces allow extension without modification
- **Liskov Substitution:** Implementations can be swapped (e.g., in-memory vs Redis)
- **Interface Segregation:** Small, focused interfaces for each use case
- **Dependency Inversion:** High-level modules depend on abstractions

---

## Package Structure

```
internal/infra/redis/
├── client.go              # Redis client with healthcheck and startup validation
├── cache.go               # Cache interface and Redis implementation
├── lock.go                # LockManager interface and Redis implementation
├── session.go             # SessionStore interface and Redis implementation
├── ratelimiter.go         # RateLimiter interface and Redis implementation
├── metrics.go             # Metrics interfaces and decorators
├── client_test.go         # Client unit tests
├── cache_test.go          # Cache unit tests
├── lock_test.go           # LockManager unit tests
├── ratelimiter_test.go    # RateLimiter unit tests
```

---

## Components

### 1. Redis Client (`client.go`)

**Purpose:** Provides a configured Redis connection with health checks and graceful shutdown.

**Key Features:**
- Connection pooling with configurable pool size
- Startup validation (ping, write, read, delete test)
- Health check with latency measurement
- Thread-safe configuration access
- Graceful shutdown with mutex protection

**Configuration:**
```go
type Config struct {
    Host         string
    Port         int
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
    MaxRetries   int
    DialTimeout  time.Duration
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    PoolTimeout  time.Duration
}
```

**Health Status:**
```go
type HealthStatus struct {
    Healthy    bool
    Latency    time.Duration
    Connected  bool
    DB         int
    ClientName string
}
```

---

### 2. Cache (`cache.go`)

**Purpose:** Generic caching interface for application-level caching.

**Interface:**
```go
type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    TTL(ctx context.Context, key string) (time.Duration, error)
    Invalidate(ctx context.Context, keys ...string) error
    SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)
}
```

**Implementation Details:**
- JSON serialization for values
- Uses `rediscmd.Nil` for key-not-found detection
- Atomic SetNX for conditional writes
- Batch invalidation support

**Use Cases:**
- Dashboard data caching (5-minute TTL)
- KPIs caching (5-minute TTL)
- Financial summary caching (5-minute TTL)
- Stock data caching (5-minute TTL)

---

### 3. Lock Manager (`lock.go`)

**Purpose:** Distributed lock implementation for coordinating access across multiple instances.

**Interface:**
```go
type LockManager interface {
    Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error)
    Release(ctx context.Context, key string) error
    TryAcquireWithRetry(ctx context.Context, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (bool, error)
}
```

**Implementation Details:**
- Uses Redis `SET NX EX` for atomic lock acquisition
- Lock expiration prevents deadlocks
- Retry logic with configurable attempts and delay
- Thread-safe operations

**Use Cases:**
- Preventing duplicate event processing
- Coordinating batch jobs
- Resource access synchronization

---

### 4. Session Store (`session.go`)

**Purpose:** Session management for user authentication and state.

**Interface:**
```go
type SessionStore interface {
    Get(ctx context.Context, sessionID string, dest interface{}) error
    Set(ctx context.Context, sessionID string, data interface{}, ttl time.Duration) error
    Delete(ctx context.Context, sessionID string) error
    Exists(ctx context.Context, sessionID string) (bool, error)
    Refresh(ctx context.Context, sessionID string, ttl time.Duration) error
    Clear(ctx context.Context, userID string) error
}
```

**Implementation Details:**
- JSON serialization for session data
- TTL-based session expiration
- User-level session invalidation
- Key pattern: `session:{sessionID}`

**Use Cases:**
- User authentication sessions
- Session refresh on activity
- Logout/session invalidation

---

### 5. Rate Limiter (`ratelimiter.go`)

**Purpose:** Rate limiting for API endpoints and operations.

**Interface:**
```go
type RateLimiter interface {
    Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
    Reset(ctx context.Context, key string) error
    GetRemaining(ctx context.Context, key string, limit int, window time.Duration) (int, error)
}
```

**Implementation Details:**
- Token bucket algorithm with Redis INCR
- Atomic counter with expiration
- Window-based rate limiting
- Remaining requests calculation

**Use Cases:**
- API endpoint rate limiting
- Operation throttling
- Resource protection

---

### 6. Metrics (`metrics.go`)

**Purpose:** Observability metrics for Redis operations.

**Interface:**
```go
type RedisMetrics interface {
    // Cache operations
    IncrementCacheHit()
    IncrementCacheMiss()
    RecordCacheOperation(operation string, duration time.Duration)
    
    // Lock operations
    IncrementLockAcquired()
    IncrementLockReleased()
    IncrementLockFailed()
    RecordLockOperation(operation string, duration time.Duration)
    
    // Rate limit operations
    IncrementRateLimitAllowed()
    IncrementRateLimitDenied()
    RecordRateLimitCheck(duration time.Duration)
    
    // Idempotency operations
    IncrementIdempotencyHit()
    IncrementIdempotencyMiss()
    RecordIdempotencyCheck(duration time.Duration)
    
    // Health check operations
    RecordHealthCheck(duration time.Duration, healthy bool)
}
```

**Decorators:**
- `RedisCacheWithMetrics` - Wraps Cache with metrics
- `RedisLockManagerWithMetrics` - Wraps LockManager with metrics
- `RedisRateLimiterWithMetrics` - Wraps RateLimiter with metrics
- `NoOpRedisMetrics` - No-op implementation for testing

---

## Idempotency Integration

### Consumer Framework Integration

**Interface:**
```go
type IdempotencyChecker interface {
    IsProcessed(ctx context.Context, eventID uint) (bool, error)
    MarkProcessed(ctx context.Context, eventID uint) error
}
```

**Implementations:**
1. **In-Memory (`IdempotencyStore`)** - Original map+mutex implementation
2. **Redis (`RedisIdempotencyStore`)** - Persistent Redis implementation

**Key Pattern:** `processed:event:{eventID}`

**TTL:** 24 hours (configurable)

**Atomic Operations:**
- Uses `SetNX` for atomic mark-processed
- Prevents race conditions in distributed environments

**Middleware Integration:**
```go
func IdempotencyMiddleware(store IdempotencyChecker, consumerName string) Middleware
```

---

## Configuration

### Environment Variables

```bash
# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s
REDIS_POOL_TIMEOUT=4s
```

### Docker Compose

```yaml
services:
  redis:
    image: redis:7-alpine
    container_name: horizongest-redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
```

---

## Key Patterns

### Cache Keys
- `dashboard:{user_id}` - Dashboard data
- `kpis:{organization_id}` - KPIs data
- `financial:{organization_id}` - Financial summary
- `stock:{product_id}` - Stock data

### Idempotency Keys
- `processed:event:{event_id}` - Event processing tracking

### Session Keys
- `session:{session_id}` - User session data

### Lock Keys
- `lock:{resource_id}` - Distributed locks

### Rate Limit Keys
- `ratelimit:{identifier}` - Rate limit counters

---

## Error Handling

### Redis Nil
- Uses `rediscmd.Nil` for key-not-found detection
- Wrapped errors with context
- Graceful degradation on Redis unavailability

### Connection Errors
- Configurable retry logic
- Startup validation fails fast
- Health check monitoring

### Timeout Handling
- Configurable dial, read, write timeouts
- Context propagation for cancellation
- Timeout metrics recording

---

## Testing

### Unit Tests
- Interface compliance tests
- Configuration tests
- Health check tests
- Mock-based operation tests (skipped without Redis connection)

### Integration Tests
- Requires actual Redis connection
- Docker Compose setup for local testing
- Test container isolation

### Test Coverage
- `client_test.go` - Client configuration and health checks
- `cache_test.go` - Cache operations
- `lock_test.go` - Lock operations
- `ratelimiter_test.go` - Rate limiting operations

---

## Performance Considerations

### Connection Pooling
- Configurable pool size (default: 10)
- Minimum idle connections (default: 5)
- Connection reuse for efficiency

### Batch Operations
- `Invalidate` supports multiple keys
- Pipeline operations for bulk writes
- Atomic operations for consistency

### TTL Management
- Automatic expiration
- Configurable TTL per use case
- Memory management through expiration

---

## Security

### Authentication
- Password support via configuration
- TLS support (configurable)
- Network isolation via Docker networks

### Key Patterns
- Namespaced keys prevent collisions
- No sensitive data in keys
- Key-based access control

### Data Protection
- Redis encryption at rest (optional)
- Network encryption (TLS)
- Access control lists

---

## Monitoring

### Health Checks
- Ping-based health check
- Latency measurement
- Connection status tracking

### Metrics
- Cache hit/miss ratios
- Lock acquisition/release rates
- Rate limit allow/deny counts
- Operation durations

### Logging
- Structured logging with context
- Error logging with stack traces
- Performance logging

---

## Future Enhancements

### Planned Features
- Redis Cluster support
- Sentinel for high availability
- Pub/Sub for real-time updates
- Lua scripts for complex operations
- Redis Streams for event sourcing

### Optimization Opportunities
- Connection pool tuning
- Pipeline optimization
- Compression for large values
- Sharding strategies

---

## Dependencies

### Go Packages
- `github.com/redis/go-redis/v9` - Redis client library
- `github.com/stretchr/testify` - Testing framework

### Infrastructure
- Redis 7.x (Docker image)
- Docker Compose for local development

---

## Migration Guide

### From In-Memory to Redis Idempotency

**Before:**
```go
store := framework.NewIdempotencyStore()
middleware := framework.IdempotencyMiddleware(store, "Consumer")
```

**After:**
```go
redisClient := redis.NewClient(redisConfig)
redisCache := redis.NewCache(redisClient)
redisIdempotency := framework.NewRedisIdempotencyStore(redisCache, 24*time.Hour)
middleware := framework.IdempotencyMiddleware(redisIdempotency, "Consumer")
```

---

## Troubleshooting

### Connection Issues
- Verify Redis is running: `docker-compose ps redis`
- Check logs: `docker-compose logs redis`
- Validate configuration in `.env`

### Performance Issues
- Monitor connection pool utilization
- Check for slow operations via Redis SLOWLOG
- Review TTL settings for memory usage

### Health Check Failures
- Verify network connectivity
- Check Redis server health
- Review timeout configurations

---

## References

- [Redis Documentation](https://redis.io/docs/)
- [go-redis v9 Documentation](https://redis.uptrace.dev/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [DDD](https://martinfowler.com/bliki/DomainDrivenDesign.html)
