## ADDED Requirements

### Requirement: Runtime API URL injection
Frontend container SHALL support dynamic API URL configuration at container startup time, not build time.

#### Scenario: Inject API URL on container start
- **WHEN** container starts with VITE_API_BASE_URL environment variable set to "http://gateway:8080"
- **THEN** Nginx configuration SHALL be updated to proxy API requests to that URL before serving traffic

#### Scenario: Default API URL when env var missing
- **WHEN** VITE_API_BASE_URL environment variable is not set
- **THEN** system SHALL use default value "http://localhost:8080"

### Requirement: Entrypoint script for configuration substitution
Frontend Dockerfile SHALL include an entrypoint script that performs environment variable substitution using envsubst.

#### Scenario: Replace placeholders in Nginx config
- **WHEN** entrypoint script runs with VITE_API_BASE_URL="http://gateway:8080"
- **THEN** all occurrences of ${VITE_API_BASE_URL} in Nginx config SHALL be replaced with actual value

#### Scenario: Preserve original config template
- **WHEN** entrypoint performs substitution
- **THEN** original template file SHALL remain unchanged and output SHALL be written to active config location

### Requirement: Nginx reverse proxy configuration
Frontend Nginx SHALL proxy /api requests to the backend gateway service.

#### Scenario: Proxy API requests to gateway
- **WHEN** browser requests /api/projects
- **THEN** Nginx SHALL forward request to ${VITE_API_BASE_URL}/api/projects

#### Scenario: Serve static files directly
- **WHEN** browser requests /assets/main.js
- **THEN** Nginx SHALL serve file directly from /usr/share/nginx/html without proxying

### Requirement: CORS handling
Frontend Nginx configuration SHALL handle CORS headers for API requests when frontend and backend are on different origins.

#### Scenario: Add CORS headers for cross-origin requests
- **WHEN** API request comes from different origin than gateway
- **THEN** Nginx SHALL add appropriate Access-Control-Allow-Origin headers

#### Scenario: Same-origin requests skip CORS
- **WHEN** frontend and gateway are accessed via same domain
- **THEN** CORS headers SHALL NOT be added (same-origin policy applies)

### Requirement: Configuration validation
Entrypoint script SHALL validate that VITE_API_BASE_URL is a valid HTTP/HTTPS URL before starting Nginx.

#### Scenario: Valid URL format
- **WHEN** VITE_API_BASE_URL is "http://gateway:8080" or "https://api.example.com"
- **THEN** validation SHALL pass and Nginx SHALL start

#### Scenario: Invalid URL format
- **WHEN** VITE_API_BASE_URL is "not-a-url" or empty string
- **THEN** entrypoint SHALL log error and exit with code 1
