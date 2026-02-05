## ADDED Requirements

### Requirement: Environment variable configuration priority
Go Gateway SHALL prioritize environment variables over config.yaml when both are present. The configuration loading order MUST be: environment variables first, then config.yaml as fallback.

#### Scenario: Environment variable overrides config.yaml
- **WHEN** both DATABASE_URL environment variable and config.yaml database.host are set
- **THEN** system SHALL use DATABASE_URL and ignore config.yaml database settings

#### Scenario: Fallback to config.yaml when env var missing
- **WHEN** DATABASE_URL environment variable is not set but config.yaml exists
- **THEN** system SHALL load database configuration from config.yaml

### Requirement: Database connection from environment
Go Gateway SHALL support DATABASE_URL environment variable in standard PostgreSQL connection string format.

#### Scenario: Parse DATABASE_URL correctly
- **WHEN** DATABASE_URL is set to "postgres://user:pass@host:5432/dbname?sslmode=disable"
- **THEN** system SHALL extract host, port, user, password, database name, and SSL mode

#### Scenario: Invalid DATABASE_URL format
- **WHEN** DATABASE_URL has invalid format
- **THEN** system SHALL log error and fail to start with clear error message

### Requirement: gRPC service address from environment
Go Gateway SHALL support GRPC_AI_SERVICE_ADDR environment variable for AI service connection.

#### Scenario: Connect to AI service via environment variable
- **WHEN** GRPC_AI_SERVICE_ADDR is set to "ai-service:50051"
- **THEN** system SHALL establish gRPC connection to that address

#### Scenario: Fallback to config.yaml for gRPC address
- **WHEN** GRPC_AI_SERVICE_ADDR is not set
- **THEN** system SHALL use grpc.ai_service_addr from config.yaml

### Requirement: JWT secret from environment
Go Gateway SHALL support JWT_SECRET environment variable for HMAC token signing.

#### Scenario: Use JWT_SECRET for token generation
- **WHEN** JWT_SECRET environment variable is set
- **THEN** system SHALL use it for HMAC-SHA256 token signing and verification

#### Scenario: Minimum JWT secret length
- **WHEN** JWT_SECRET is shorter than 32 characters
- **THEN** system SHALL log warning about weak secret but continue operation

### Requirement: Configuration validation on startup
Go Gateway SHALL validate all required configuration values on startup and fail fast if critical values are missing.

#### Scenario: Missing database configuration
- **WHEN** neither DATABASE_URL nor config.yaml database settings are available
- **THEN** system SHALL exit with error code 1 and log "Database configuration missing"

#### Scenario: Missing JWT secret
- **WHEN** neither JWT_SECRET nor config.yaml auth.jwt_secret are available
- **THEN** system SHALL exit with error code 1 and log "JWT secret not configured"
