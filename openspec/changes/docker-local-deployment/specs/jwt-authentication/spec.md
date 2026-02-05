## MODIFIED Requirements

### Requirement: HMAC-based JWT authentication
Go Gateway SHALL use HMAC-SHA256 algorithm for JWT token signing and verification, replacing the previous RS256 public-key cryptography approach.

#### Scenario: Generate JWT token with HMAC
- **WHEN** user successfully authenticates
- **THEN** system SHALL generate JWT token signed with HMAC-SHA256 using JWT_SECRET

#### Scenario: Verify JWT token with HMAC
- **WHEN** request includes JWT token in Authorization header
- **THEN** system SHALL verify token signature using HMAC-SHA256 with same JWT_SECRET

#### Scenario: Reject token with invalid signature
- **WHEN** JWT token signature does not match expected HMAC
- **THEN** system SHALL return HTTP 401 Unauthorized with error "Invalid token signature"

### Requirement: JWT secret from environment variable
Go Gateway SHALL read JWT secret from JWT_SECRET environment variable instead of loading RSA private key from file.

#### Scenario: Load JWT secret on startup
- **WHEN** gateway service starts
- **THEN** system SHALL read JWT_SECRET environment variable and use it for all JWT operations

#### Scenario: Fail fast when JWT secret missing
- **WHEN** JWT_SECRET environment variable is not set and config.yaml does not contain auth.jwt_secret
- **THEN** system SHALL exit with error code 1 and log "JWT_SECRET not configured"

### Requirement: Remove RSA key file dependencies
Go Gateway SHALL NOT require private.pem or public.pem files for JWT operations.

#### Scenario: Start without key files
- **WHEN** gateway starts with JWT_SECRET set but no keys/ directory
- **THEN** system SHALL start successfully and generate/verify tokens using HMAC

#### Scenario: Ignore key file paths in config
- **WHEN** config.yaml contains jwt_private_key_path or jwt_public_key_path
- **THEN** system SHALL ignore these settings and use JWT_SECRET instead

### Requirement: Token payload compatibility
JWT token payload structure SHALL remain unchanged to maintain compatibility with existing clients.

#### Scenario: Token contains same claims
- **WHEN** new HMAC token is generated
- **THEN** payload SHALL contain same claims (user_id, email, exp, iat) as previous RS256 tokens

#### Scenario: Existing tokens become invalid
- **WHEN** system switches from RS256 to HMAC
- **THEN** all previously issued RS256 tokens SHALL be rejected with "Invalid token signature" error

**Reason**: Simplify deployment by removing key file management complexity and align with docker-compose environment variable configuration pattern.

**Migration**:
1. Generate strong JWT_SECRET (minimum 32 characters): `openssl rand -base64 32`
2. Add JWT_SECRET to .env file
3. Remove JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH from docker-compose.yml
4. Remove keys/ directory volume mount
5. All users will need to re-authenticate after deployment (existing tokens invalid)

### Requirement: JWT secret strength validation
Go Gateway SHALL validate JWT_SECRET meets minimum security requirements.

#### Scenario: Warn on weak secret
- **WHEN** JWT_SECRET is shorter than 32 characters
- **THEN** system SHALL log WARNING "JWT_SECRET is shorter than recommended 32 characters" but continue operation

#### Scenario: Accept strong secret
- **WHEN** JWT_SECRET is 32 or more characters
- **THEN** system SHALL start without warnings related to secret strength
