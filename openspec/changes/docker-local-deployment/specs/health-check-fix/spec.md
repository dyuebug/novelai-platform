## ADDED Requirements

### Requirement: curl installation in AI service container
Python AI service Dockerfile SHALL include curl package installation to support health check commands.

#### Scenario: Install curl during image build
- **WHEN** AI service Docker image is built
- **THEN** Dockerfile SHALL execute apt-get update && apt-get install -y curl

#### Scenario: Minimize image size impact
- **WHEN** curl is installed
- **THEN** apt cache SHALL be cleaned in same RUN layer to minimize image size increase

### Requirement: Health check endpoint availability
AI service SHALL expose /health endpoint that returns 200 OK when service is ready to accept requests.

#### Scenario: Health check returns success when ready
- **WHEN** AI service has completed initialization and is ready
- **THEN** GET /health SHALL return HTTP 200 with response body {"status": "healthy"}

#### Scenario: Health check returns failure during startup
- **WHEN** AI service is still initializing
- **THEN** GET /health SHALL return HTTP 503 with response body {"status": "starting"}

### Requirement: Docker health check configuration
Docker Compose SHALL configure health check for AI service using curl command.

#### Scenario: Health check command execution
- **WHEN** Docker health check runs
- **THEN** system SHALL execute "curl -f http://localhost:8001/health"

#### Scenario: Health check interval and timeout
- **WHEN** health check is configured
- **THEN** interval SHALL be 30 seconds, timeout SHALL be 10 seconds, retries SHALL be 3

#### Scenario: Service marked unhealthy after retries
- **WHEN** health check fails 3 consecutive times
- **THEN** Docker SHALL mark container as unhealthy and dependent services SHALL NOT start

### Requirement: Health check dependency chain
Gateway service SHALL depend on AI service health check passing before starting.

#### Scenario: Gateway waits for AI service health
- **WHEN** docker compose up is executed
- **THEN** gateway container SHALL NOT start until ai-service health check returns healthy

#### Scenario: Gateway starts after AI service ready
- **WHEN** ai-service health check passes
- **THEN** gateway container SHALL start within 5 seconds

### Requirement: Health check logging
AI service SHALL log health check requests at DEBUG level to avoid log spam.

#### Scenario: Health check requests logged at debug level
- **WHEN** health check endpoint is called
- **THEN** log entry SHALL use DEBUG level with format "Health check: OK"

#### Scenario: Health check failures logged at warning level
- **WHEN** health check endpoint returns non-200 status
- **THEN** log entry SHALL use WARNING level with reason for failure
