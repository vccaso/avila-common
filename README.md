# Avila Common Library

Shared utilities, models, and architecture patterns for the Avila microservices ecosystem.

## 📦 Core Packages

### 1. `event`
Implements the core logic for the **Hybrid Event Architecture**:
- **`relay.go`**: Reusable `OutboxRelay` worker that handles concurrency-safe polling, locking, and publishing to NATS JetStream.
- **`consumer.go`**: Reusable `ProcessIdempotentEvent` wrapper that handles duplicate checks and transaction management for reliable consumption.

### 2. `util`
- **Logging**: Structured logging (Info, Warning, Error).
- **Environment**: Helpers for parsing environment variables.

### 3. `auth`
- Shared authentication and JWT verification middleware.

### 4. `model`
- Shared domain models used across multiple services.

## ⚙️ Dependencies
- `github.com/nats-io/nats.go`: Message Bus connectivity.
- `github.com/google/uuid`: Unique identification for events and workers.
