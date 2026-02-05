## ADDED Requirements

### Requirement: Independent initialization container
Docker Compose SHALL include a dedicated init-db service that runs database migrations automatically before application services start.

#### Scenario: Init container runs before application services
- **WHEN** docker compose up is executed
- **THEN** init-db service SHALL start after postgres is healthy and complete before gateway/ai-service start

#### Scenario: Migration success allows app startup
- **WHEN** init-db service completes successfully (exit code 0)
- **THEN** gateway and ai-service containers SHALL be allowed to start

#### Scenario: Migration failure blocks app startup
- **WHEN** init-db service fails (exit code non-zero)
- **THEN** gateway and ai-service containers SHALL NOT start and docker compose SHALL report error

### Requirement: Automatic migration execution
Init-db service SHALL automatically execute all pending database migrations using golang-migrate tool.

#### Scenario: Apply all pending migrations
- **WHEN** init-db container starts and database has migrations pending
- **THEN** system SHALL execute all migrations in order from current version to latest

#### Scenario: No migrations needed
- **WHEN** init-db container starts and database is already at latest version
- **THEN** system SHALL log "No change" and exit successfully without error

#### Scenario: Idempotent migration execution
- **WHEN** init-db service runs multiple times
- **THEN** already-applied migrations SHALL be skipped and only new migrations SHALL execute

### Requirement: PostgreSQL extension creation
Init-db service SHALL ensure required PostgreSQL extensions (uuid-ossp, pgvector) are created before running migrations.

#### Scenario: Create extensions if missing
- **WHEN** init-db runs and extensions do not exist
- **THEN** system SHALL execute CREATE EXTENSION IF NOT EXISTS for uuid-ossp and vector

#### Scenario: Extensions already exist
- **WHEN** init-db runs and extensions already exist
- **THEN** system SHALL skip extension creation without error

### Requirement: Database connection validation
Init-db service SHALL validate database connectivity before attempting migrations.

#### Scenario: Wait for database readiness
- **WHEN** init-db starts but postgres is not yet accepting connections
- **THEN** system SHALL retry connection every 2 seconds for up to 30 seconds

#### Scenario: Database connection timeout
- **WHEN** postgres does not become ready within 30 seconds
- **THEN** init-db SHALL exit with error code 1 and log "Database connection timeout"

### Requirement: Migration logging
Init-db service SHALL log all migration operations with timestamps and version numbers.

#### Scenario: Log successful migration
- **WHEN** migration 001_create_users.up.sql executes successfully
- **THEN** system SHALL log "Applied migration 001_create_users.up.sql at [timestamp]"

#### Scenario: Log migration failure with details
- **WHEN** migration fails due to SQL error
- **THEN** system SHALL log error message, SQL statement, and line number before exiting

### Requirement: Container cleanup
Init-db container SHALL automatically remove itself after completion to avoid cluttering container list.

#### Scenario: Remove container on success
- **WHEN** init-db completes successfully
- **THEN** container SHALL be automatically removed (docker compose --rm behavior)

#### Scenario: Preserve container on failure
- **WHEN** init-db fails
- **THEN** container SHALL remain for debugging with logs accessible via docker compose logs init-db
