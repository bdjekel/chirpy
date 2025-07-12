# chirpy

## Project Overview

**Chirpy** is a full-featured Twitter-style microblogging API built in Go, implementing a complete social media backend with user authentication, content management, and real-time features. The application serves as a robust REST API that handles user registration, authentication, post creation, content moderation, and membership management.


## Table of Contents

- [Detailed Technical Documentation](#detailed-technical-documentation)
  - [Architecture & Design](#architecture--design)
  - [Core Features](#core-features)
    - [1. User Management System](#1-user-management-system)
    - [2. Content Management](#2-content-management)
    - [3. Security Features](#3-security-features)
    - [4. Monitoring & Administration](#4-monitoring--administration)
  - [Database Schema](#database-schema)
  - [API Endpoints](#api-endpoints)
    - [Authentication Endpoints](#authentication-endpoints)


## Detailed Technical Documentation

### Architecture & Design

The application follows a clean, modular architecture with clear separation of concerns:

- **Main Application Layer**: Handles HTTP routing, middleware, and request/response processing
- **Business Logic Layer**: Manages chirp creation, user operations, and authentication flows
- **Data Access Layer**: Uses SQLC for type-safe database operations with PostgreSQL
- **Authentication Layer**: Implements JWT-based authentication with refresh token rotation

### Core Features

#### 1. User Management System
- **User Registration**: Secure user creation with bcrypt password hashing
- **Authentication**: JWT-based login with access and refresh token system
- **Credential Updates**: Secure password and email update functionality
- **Membership Tiers**: Premium user upgrade system via webhook integration

#### 2. Content Management
- **Chirp Creation**: 140-character limit posts with profanity filtering
- **Content Retrieval**: Paginated chirp retrieval with author filtering and sorting
- **Content Moderation**: Automatic profanity detection and replacement
- **Content Ownership**: User-specific chirp deletion with authorization checks

#### 3. Security Features
- **Password Security**: bcrypt hashing with configurable cost
- **Token Management**: JWT access tokens with refresh token rotation
- **Authorization**: Bearer token validation for protected endpoints
- **API Key Protection**: Webhook authentication for external integrations

#### 4. Monitoring & Administration
- **Metrics Tracking**: Hit counter middleware for request monitoring
- **Health Checks**: Readiness endpoint for service health monitoring
- **Admin Controls**: Database reset and metrics viewing capabilities

### Database Schema

The application uses PostgreSQL with the following core tables:

- **users**: User accounts with email, hashed passwords, and membership status
- **chirps**: User posts with content, timestamps, and foreign key relationships
- **refresh_tokens**: Secure token storage with expiration and revocation tracking

### API Endpoints

#### Authentication Endpoints
- `POST /api/users` - User registration
- `POST /api/login` - User authentication
- `POST /api/refresh` - Token refresh
- `POST /api/revoke` - Token revocation
- `PUT /api/users` - Credential updates

#### Content Endpoints
- `GET /api/chirps` - Retrieve chirps (with optional filtering)
- `GET /api/chirps/{id}` - Get specific chirp
- `POST /api/chirps` - Create new chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp (authorized users only)

#### Webhook Endpoints
- `POST /api/polka/webhooks` - Membership upgrade processing

#### Administrative Endpoints
- `GET /api/healthz` - Health check
- `GET /admin/metrics` - View hit metrics
- `POST /admin/reset` - Reset application state

### Technical Implementation Details

#### Authentication Flow
1. User login validates credentials against bcrypt-hashed passwords
2. JWT access tokens are generated with 1-hour expiration
3. Refresh tokens are stored in database with expiration tracking
4. Token rotation prevents long-term session hijacking

#### Content Moderation
- Profanity filtering uses O(1) map lookups for efficiency
- Configurable word replacement with asterisk masking
- Case-insensitive matching for comprehensive filtering

#### Database Operations
- SQLC generates type-safe Go code from SQL queries
- Prepared statements prevent SQL injection
- Foreign key constraints ensure data integrity
- Cascade deletes maintain referential integrity

#### Error Handling
- Structured error responses with appropriate HTTP status codes
- Comprehensive logging for debugging and monitoring
- Graceful error recovery with user-friendly messages

## Technology Stack & Skills

### Languages
- **Go (Golang)** - Primary application language
- **SQL** - Database queries and schema management
- **HTML** - Basic frontend interface

### Frameworks & Libraries
- **Standard HTTP Library** - Go's built-in HTTP server and routing
- **SQLC** - Type-safe SQL code generation
- **JWT (golang-jwt)** - JSON Web Token implementation
- **bcrypt (golang.org/x/crypto)** - Password hashing
- **UUID (google/uuid)** - Unique identifier generation
- **godotenv** - Environment variable management
- **lib/pq** - PostgreSQL driver

### Design Patterns
- **Repository Pattern** - Database access abstraction through SQLC
- **Middleware Pattern** - Request/response processing pipeline
- **Dependency Injection** - Configuration struct passing
- **Factory Pattern** - Database connection and query creation
- **Strategy Pattern** - Different authentication methods (JWT, API Key)

### Algorithms
- **bcrypt Hashing** - Secure password hashing with salt
- **JWT Token Generation** - HMAC-SHA256 signing
- **Profanity Filtering** - O(1) hash map lookups for word replacement
- **Token Rotation** - Secure session management

### Data Structures
- **Maps** - Profanity filtering and configuration storage
- **Slices** - Chirp collections and sorting
- **Structs** - Request/response models and database entities
- **Atomic Integers** - Thread-safe hit counter
- **UUID** - Unique identifier generation and storage

### Database & Infrastructure
- **PostgreSQL** - Primary database with ACID compliance
- **SQL Migrations** - Version-controlled schema management
- **Foreign Key Constraints** - Referential integrity
- **Indexed Queries** - Optimized data retrieval

### Security Practices
- **Password Hashing** - bcrypt with configurable cost
- **JWT Authentication** - Stateless token-based auth
- **Refresh Token Rotation** - Session security
- **Input Validation** - Request parameter sanitization
- **SQL Injection Prevention** - Prepared statements via SQLC
- **CORS Handling** - Cross-origin request management

### Development Tools & Practices
- **Environment Configuration** - .env file management
- **Error Handling** - Structured error responses
- **Logging** - Application monitoring and debugging
- **Health Checks** - Service availability monitoring
- **Metrics Collection** - Request tracking and analytics

## OpenAPI Specification

```yaml
openapi: 3.1.1
info:
  title: Chirpy API
  description: A Twitter-style microblogging API with user authentication, content management, and real-time features
  version: 1.0.0
  contact:
    name: Chirpy API Support
    url: https://github.com/bdjekel/chirpy
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: http://localhost:8080
    description: Development server
  - url: https://api.chirpy.com
    description: Production server

paths:
  /api/healthz:
    get:
      summary: Health check endpoint
      description: Returns the health status of the API
      operationId: getHealth
      tags:
        - Health
      responses:
        '200':
          description: API is healthy
          content:
            text/plain:
              schema:
                type: string
                example: "OK"

  /api/users:
    post:
      summary: Create a new user
      description: Register a new user account with email and password
      operationId: createUser
      tags:
        - Authentication
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
            examples:
              valid_user:
                summary: Valid user registration
                value:
                  email: "user@example.com"
                  password: "securepassword123"
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'
              examples:
                created_user:
                  summary: User created
                  value:
                    id: "123e4567-e89b-12d3-a456-426614174000"
                    created_at: "2024-01-15T10:30:00Z"
                    updated_at: "2024-01-15T10:30:00Z"
                    email: "user@example.com"
                    is_chirpy_red: false
        '400':
          description: Invalid request data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/login:
    post:
      summary: Authenticate user
      description: Login with email and password to receive access and refresh tokens
      operationId: loginUser
      tags:
        - Authentication
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
            examples:
              valid_login:
                summary: Valid login credentials
                value:
                  email: "user@example.com"
                  password: "securepassword123"
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'
              examples:
                successful_login:
                  summary: Successful login
                  value:
                    id: "123e4567-e89b-12d3-a456-426614174000"
                    created_at: "2024-01-15T10:30:00Z"
                    updated_at: "2024-01-15T10:30:00Z"
                    email: "user@example.com"
                    is_chirpy_red: false
                    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
                    refresh_token: "abc123def456ghi789"
        '401':
          description: Invalid credentials
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/refresh:
    post:
      summary: Refresh access token
      description: Use refresh token to get a new access token
      operationId: refreshToken
      tags:
        - Authentication
      security:
        - BearerAuth: []
      responses:
        '200':
          description: Token refreshed successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/RefreshResponse'
              examples:
                refreshed_token:
                  summary: Token refreshed
                  value:
                    token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        '401':
          description: Invalid or expired refresh token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/revoke:
    post:
      summary: Revoke refresh token
      description: Revoke a refresh token to invalidate it
      operationId: revokeToken
      tags:
        - Authentication
      security:
        - BearerAuth: []
      responses:
        '204':
          description: Token revoked successfully
        '401':
          description: Invalid refresh token
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/users:
    put:
      summary: Update user credentials
      description: Update user email and password
      operationId: updateUser
      tags:
        - Authentication
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateUserRequest'
            examples:
              update_credentials:
                summary: Update user credentials
                value:
                  email: "newemail@example.com"
                  password: "newsecurepassword123"
      responses:
        '200':
          description: User updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserResponse'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/chirps:
    get:
      summary: Get chirps
      description: Retrieve chirps with optional filtering by author and sorting
      operationId: getChirps
      tags:
        - Content
      parameters:
        - name: author_id
          in: query
          description: Filter chirps by author ID
          required: false
          schema:
            type: string
            format: uuid
            example: "123e4567-e89b-12d3-a456-426614174000"
        - name: sort
          in: query
          description: Sort order for chirps
          required: false
          schema:
            type: string
            enum: [asc, desc]
            default: asc
            example: "desc"
      responses:
        '200':
          description: Chirps retrieved successfully
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Chirp'
              examples:
                chirps_list:
                  summary: List of chirps
                  value:
                    - id: "123e4567-e89b-12d3-a456-426614174000"
                      created_at: "2024-01-15T10:30:00Z"
                      updated_at: "2024-01-15T10:30:00Z"
                      body: "Hello, Chirpy world!"
                      user_id: "123e4567-e89b-12d3-a456-426614174000"
                    - id: "123e4567-e89b-12d3-a456-426614174001"
                      created_at: "2024-01-15T09:15:00Z"
                      updated_at: "2024-01-15T09:15:00Z"
                      body: "This is my second chirp!"
                      user_id: "123e4567-e89b-12d3-a456-426614174000"
        '400':
          description: Invalid parameters
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    post:
      summary: Create a new chirp
      description: Create a new chirp with content (max 140 characters)
      operationId: createChirp
      tags:
        - Content
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateChirpRequest'
            examples:
              valid_chirp:
                summary: Valid chirp content
                value:
                  body: "Hello, Chirpy world! This is my first chirp."
              profanity_filtered:
                summary: Chirp with profanity (will be filtered)
                value:
                  body: "This kerfuffle is getting out of hand!"
      responses:
        '201':
          description: Chirp created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Chirp'
              examples:
                created_chirp:
                  summary: Chirp created
                  value:
                    id: "123e4567-e89b-12d3-a456-426614174000"
                    created_at: "2024-01-15T10:30:00Z"
                    updated_at: "2024-01-15T10:30:00Z"
                    body: "Hello, Chirpy world! This is my first chirp."
                    user_id: "123e4567-e89b-12d3-a456-426614174000"
        '400':
          description: Invalid chirp data (e.g., too long)
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/chirps/{id}:
    get:
      summary: Get chirp by ID
      description: Retrieve a specific chirp by its ID
      operationId: getChirpById
      tags:
        - Content
      parameters:
        - name: id
          in: path
          required: true
          description: Chirp ID
          schema:
            type: string
            format: uuid
            example: "123e4567-e89b-12d3-a456-426614174000"
      responses:
        '200':
          description: Chirp retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Chirp'
        '404':
          description: Chirp not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/chirps/{chirpID}:
    delete:
      summary: Delete chirp
      description: Delete a chirp (only by the author)
      operationId: deleteChirp
      tags:
        - Content
      security:
        - BearerAuth: []
      parameters:
        - name: chirpID
          in: path
          required: true
          description: Chirp ID to delete
          schema:
            type: string
            format: uuid
            example: "123e4567-e89b-12d3-a456-426614174000"
      responses:
        '204':
          description: Chirp deleted successfully
        '403':
          description: Forbidden - not the chirp author
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: Chirp not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /api/polka/webhooks:
    post:
      summary: Process membership upgrade webhook
      description: Handle webhook for user membership upgrades
      operationId: processMembershipUpgrade
      tags:
        - Webhooks
      security:
        - ApiKeyAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/WebhookRequest'
            examples:
              user_upgraded:
                summary: User upgraded event
                value:
                  event: "user.upgraded"
                  data:
                    user_id: "123e4567-e89b-12d3-a456-426614174000"
      responses:
        '204':
          description: Webhook processed successfully
        '401':
          description: Invalid API key
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: User not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /admin/metrics:
    get:
      summary: Get application metrics
      description: Retrieve application hit metrics (admin only)
      operationId: getMetrics
      tags:
        - Administration
      responses:
        '200':
          description: Metrics retrieved successfully
          content:
            text/html:
              schema:
                type: string
                example: |
                  <html>
                  <body>
                    <h1>Welcome, Chirpy Admin</h1>
                    <p>Chirpy has been visited 42 times!</p>
                  </body>
                  </html>

  /admin/reset:
    post:
      summary: Reset application state
      description: Reset hit counter and database (admin only)
      operationId: resetApplication
      tags:
        - Administration
      responses:
        '200':
          description: Application reset successfully
          content:
            application/json:
              schema:
                type: string
                example: "Hits reset to 0. Database reset to initial state."
        '403':
          description: Forbidden - admin privileges required
          content:
            application/json:
              schema:
                type: string
                example: "Reset requires admin privileges"

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT token for user authentication
    ApiKeyAuth:
      type: apiKey
      in: header
      name: Authorization
      description: API key for webhook authentication (format: "ApiKey <key>")

  schemas:
    CreateUserRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
          description: User's email address
          example: "user@example.com"
        password:
          type: string
          minLength: 1
          description: User's password (will be hashed)
          example: "securepassword123"

    LoginRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
          description: User's email address
          example: "user@example.com"
        password:
          type: string
          description: User's password
          example: "securepassword123"

    UpdateUserRequest:
      type: object
      required:
        - email
        - password
      properties:
        email:
          type: string
          format: email
          description: New email address
          example: "newemail@example.com"
        password:
          type: string
          minLength: 1
          description: New password (will be hashed)
          example: "newsecurepassword123"

    CreateChirpRequest:
      type: object
      required:
        - body
      properties:
        body:
          type: string
          maxLength: 140
          description: Chirp content (max 140 characters)
          example: "Hello, Chirpy world! This is my first chirp."

    WebhookRequest:
      type: object
      required:
        - event
        - data
      properties:
        event:
          type: string
          enum: [user.upgraded]
          description: Webhook event type
          example: "user.upgraded"
        data:
          type: object
          required:
            - user_id
          properties:
            user_id:
              type: string
              format: uuid
              description: User ID for the event
              example: "123e4567-e89b-12d3-a456-426614174000"

    UserResponse:
      type: object
      properties:
        id:
          type: string
          format: uuid
          description: User's unique identifier
          example: "123e4567-e89b-12d3-a456-426614174000"
        created_at:
          type: string
          format: date-time
          description: User creation timestamp
          example: "2024-01-15T10:30:00Z"
        updated_at:
          type: string
          format: date-time
          description: User last update timestamp
          example: "2024-01-15T10:30:00Z"
        email:
          type: string
          format: email
          description: User's email address
          example: "user@example.com"
        is_chirpy_red:
          type: boolean
          description: Premium membership status
          example: false

    LoginResponse:
      allOf:
        - $ref: '#/components/schemas/UserResponse'
        - type: object
          properties:
            token:
              type: string
              description: JWT access token
              example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
            refresh_token:
              type: string
              description: Refresh token for getting new access tokens
              example: "abc123def456ghi789"

    RefreshResponse:
      type: object
      properties:
        token:
          type: string
          description: New JWT access token
          example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

    Chirp:
      type: object
      properties:
        id:
          type: string
          format: uuid
          description: Chirp's unique identifier
          example: "123e4567-e89b-12d3-a456-426614174000"
        created_at:
          type: string
          format: date-time
          description: Chirp creation timestamp
          example: "2024-01-15T10:30:00Z"
        updated_at:
          type: string
          format: date-time
          description: Chirp last update timestamp
          example: "2024-01-15T10:30:00Z"
        body:
          type: string
          maxLength: 140
          description: Chirp content
          example: "Hello, Chirpy world! This is my first chirp."
        user_id:
          type: string
          format: uuid
          description: ID of the chirp author
          example: "123e4567-e89b-12d3-a456-426614174000"

    ErrorResponse:
      type: object
      properties:
        error:
          type: string
          description: Error message
          example: "Invalid request parameters"

tags:
  - name: Authentication
    description: User authentication and authorization endpoints
  - name: Content
    description: Chirp creation, retrieval, and management
  - name: Webhooks
    description: External service integration endpoints
  - name: Administration
    description: Admin-only endpoints for monitoring and management
  - name: Health
    description: Health check and monitoring endpoints