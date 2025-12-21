# Request Flows - FinFlow API

## 1. Create User (POST /users)

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant Handler
    participant Validator
    participant UserService
    participant HashService
    participant Repository
    participant Database

    Client->>CORS: POST /users {first_name, last_name, email, password}
    CORS->>CORS: Add CORS headers
    CORS->>Handler: Request
    Handler->>Handler: Verify POST method
    Handler->>Handler: Decode JSON
    Handler->>Validator: Validate fields
    Validator->>Validator: Validate first_name (required, min 2 chars)
    Validator->>Validator: Validate last_name (required, min 2 chars)
    Validator->>Validator: Validate email (required, valid format)
    Validator->>Validator: Validate password (required, min 8 chars)
    alt Validation fails
        Validator-->>Handler: Validation error
        Handler-->>Client: 400 Bad Request
    else Validation OK
        Validator-->>Handler: OK
        Handler->>UserService: Create(request)
        UserService->>HashService: Hash(password)
        HashService-->>UserService: hashedPassword
        UserService->>UserService: Generate UUID
        UserService->>UserService: NewUser(id, firstName, lastName, email, hashedPassword, systemUser)
        UserService->>Repository: Create(user)
        Repository->>Database: INSERT INTO users
        Database-->>Repository: OK
        Repository-->>UserService: OK
        UserService-->>Handler: OK
        Handler-->>Client: 200 Success
    end
```

## 2. Login (POST /auth/login)

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant AuthHandler
    participant Repository
    participant HashService
    participant JWTService
    participant Database

    Client->>CORS: POST /auth/login {email, password}
    CORS->>CORS: Add CORS headers
    CORS->>AuthHandler: Request
    AuthHandler->>AuthHandler: Verify POST method
    AuthHandler->>AuthHandler: Decode JSON
    AuthHandler->>Repository: GetByEmail(email)
    Repository->>Database: SELECT * FROM users WHERE email = ?
    alt User not found
        Database-->>Repository: No rows
        Repository-->>AuthHandler: Error
        AuthHandler-->>Client: 401 Unauthorized
    else User found
        Database-->>Repository: User data
        Repository-->>AuthHandler: User
        AuthHandler->>HashService: Verify(password, user.Password)
        alt Password incorrect
            HashService-->>AuthHandler: false
            AuthHandler-->>Client: 401 Unauthorized
        else Password correct
            HashService-->>AuthHandler: true
            AuthHandler->>JWTService: GenerateToken(userID)
            JWTService-->>AuthHandler: token
            AuthHandler->>AuthHandler: Create LoginResponse {token, user}
            AuthHandler-->>Client: 200 OK {token, user}
        end
    end
```

## 3. Get User (GET /users/{id})

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant AuthMiddleware
    participant Handler
    participant UserService
    participant Repository
    participant Database

    Client->>CORS: GET /users/123<br/>Header: Authorization: Bearer token
    CORS->>CORS: Add CORS headers
    CORS->>AuthMiddleware: Request
    AuthMiddleware->>AuthMiddleware: Extract token from header
    AuthMiddleware->>AuthMiddleware: Validate JWT token
    alt Invalid token
        AuthMiddleware-->>Client: 401 Unauthorized
    else Valid token
        AuthMiddleware->>AuthMiddleware: Extract userID from token
        AuthMiddleware->>AuthMiddleware: Save userID in context
        AuthMiddleware->>Handler: Request with context
        Handler->>Handler: Extract ID from URL (/users/123)
        Handler->>Handler: Get userID from context
        Handler->>Handler: Verify authorization (userID == id)
        alt Not authorized
            Handler-->>Client: 403 Forbidden
        else Authorized
            Handler->>UserService: GetByID(id)
            UserService->>Repository: GetByID(id)
            Repository->>Database: SELECT * FROM users WHERE id = ?
            alt User not found
                Database-->>Repository: No rows
                Repository-->>UserService: Error
                UserService-->>Handler: Error
                Handler-->>Client: 404 Not Found
            else User found
                Database-->>Repository: User data
                Repository-->>UserService: User
                UserService-->>Handler: UserResponse
                Handler->>Handler: Convert to UserResponse
                Handler-->>Client: 200 OK {id, first_name, last_name, email}
            end
        end
    end
```

## 4. Update User (PUT /users/{id})

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant AuthMiddleware
    participant Handler
    participant UserService
    participant Repository
    participant Database

    Client->>CORS: PUT /users/123<br/>Header: Authorization: Bearer token<br/>Body: {first_name, last_name, email}
    CORS->>CORS: Add CORS headers
    CORS->>AuthMiddleware: Request
    AuthMiddleware->>AuthMiddleware: Validate JWT token
    alt Invalid token
        AuthMiddleware-->>Client: 401 Unauthorized
    else Valid token
        AuthMiddleware->>AuthMiddleware: Save userID in context
        AuthMiddleware->>Handler: Request with context
        Handler->>Handler: Extract ID from URL
        Handler->>Handler: Get userID from context
        Handler->>Handler: Verify authorization (userID == id)
        alt Not authorized
            Handler-->>Client: 403 Forbidden
        else Authorized
            Handler->>Handler: Decode JSON body
            Handler->>UserService: Update(request)
            UserService->>Repository: GetByID(id)
            Repository->>Database: SELECT * FROM users WHERE id = ?
            Database-->>Repository: User data
            Repository-->>UserService: User
            UserService->>UserService: Update fields
            UserService->>UserService: UpdateModified(systemUser)
            UserService->>Repository: Update(user)
            Repository->>Database: UPDATE users SET ...
            Database-->>Repository: OK
            Repository-->>UserService: OK
            UserService-->>Handler: OK
            Handler-->>Client: 200 Success
        end
    end
```

## 5. Delete User (DELETE /users/{id})

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant AuthMiddleware
    participant Handler
    participant UserService
    participant Repository
    participant Database

    Client->>CORS: DELETE /users/123<br/>Header: Authorization: Bearer token
    CORS->>CORS: Add CORS headers
    CORS->>AuthMiddleware: Request
    AuthMiddleware->>AuthMiddleware: Validate JWT token
    alt Invalid token
        AuthMiddleware-->>Client: 401 Unauthorized
    else Valid token
        AuthMiddleware->>AuthMiddleware: Save userID in context
        AuthMiddleware->>Handler: Request with context
        Handler->>Handler: Extract ID from URL
        Handler->>Handler: Get userID from context
        Handler->>Handler: Verify authorization (userID == id)
        alt Not authorized
            Handler-->>Client: 403 Forbidden
        else Authorized
            Handler->>UserService: Delete(id)
            UserService->>Repository: Delete(id)
            Repository->>Database: DELETE FROM users WHERE id = ?
            alt User not found
                Database-->>Repository: 0 rows affected
                Repository-->>UserService: Error
                UserService-->>Handler: Error
                Handler-->>Client: 404 Not Found
            else User deleted
                Database-->>Repository: 1 row affected
                Repository-->>UserService: OK
                UserService-->>Handler: OK
                Handler-->>Client: 200 Success
            end
        end
    end
```

## 6. List Users (GET /users)

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant AuthMiddleware
    participant Handler
    participant UserService
    participant Repository
    participant Database

    Client->>CORS: GET /users<br/>Header: Authorization: Bearer token
    CORS->>CORS: Add CORS headers
    CORS->>AuthMiddleware: Request
    AuthMiddleware->>AuthMiddleware: Validate JWT token
    alt Invalid token
        AuthMiddleware-->>Client: 401 Unauthorized
    else Valid token
        AuthMiddleware->>Handler: Request
        Handler->>UserService: List()
        UserService->>Repository: List()
        Repository->>Database: SELECT * FROM users ORDER BY created_at DESC
        Database-->>Repository: Users array
        Repository-->>UserService: []User
        UserService->>UserService: Convert to []UserResponse
        UserService-->>Handler: []UserResponse
        Handler-->>Client: 200 OK [{id, first_name, last_name, email}, ...]
    end
```

## 7. Sync User from Clerk (POST /users/sync)

```mermaid
sequenceDiagram
    participant Client
    participant CORS
    participant ClerkMiddleware
    participant Handler
    participant UserService
    participant Repository
    participant Database

    Client->>CORS: POST /users/sync<br/>Header: Authorization: Bearer clerk_token
    CORS->>CORS: Add CORS headers
    CORS->>ClerkMiddleware: Request
    ClerkMiddleware->>ClerkMiddleware: Extract token from header
    ClerkMiddleware->>ClerkMiddleware: Parse JWT claims
    alt Invalid token
        ClerkMiddleware-->>Client: 401 Unauthorized
    else Valid token
        ClerkMiddleware->>ClerkMiddleware: Extract claims (sub, given_name, family_name, email)
        ClerkMiddleware->>ClerkMiddleware: Save authID, firstName, lastName, email in context
        ClerkMiddleware->>Handler: Request with context
        Handler->>Handler: Get authID, firstName, lastName, email from context
        Handler->>UserService: SyncByAuthID(authID, firstName, lastName, email)
        UserService->>Repository: GetByAuthID(authID)
        Repository->>Database: SELECT * FROM users WHERE auth_id = ?
        alt User not found
            Database-->>Repository: No rows
            Repository-->>UserService: Error (user not found)
            UserService->>UserService: Generate UUID
            UserService->>UserService: NewUserWithAuthID(id, authID, firstName, lastName, email, "", systemUser)
            UserService->>Repository: Create(newUser)
            Repository->>Database: INSERT INTO users (id, auth_id, first_name, last_name, email, ...)
            Database-->>Repository: OK
            Repository-->>UserService: OK
            UserService-->>Handler: UserResponse
            Handler-->>Client: 200 OK {id, first_name, last_name, email}
        else User found
            Database-->>Repository: User data
            Repository-->>UserService: User
            UserService-->>Handler: UserResponse
            Handler-->>Client: 200 OK {id, first_name, last_name, email}
        end
    end
```

## General Architecture Diagram

```mermaid
graph TB
    Client[Frontend/Client]

    subgraph HTTP["HTTP Layer"]
        CORS[CORS Middleware]
        Auth[Auth Middleware]
        ClerkAuth[Clerk Auth Middleware]
        Routes[Routes]
        Handlers[HTTP Handlers]
    end

    subgraph APP["Application Layer"]
        UserService[UserService]
        AuthHandler[AuthHandler]
        Validator[Validator]
    end

    subgraph DOMAIN["Domain Layer"]
        User[User Entity]
        Repository[Repository Interface]
    end

    subgraph INFRA["Infrastructure Layer"]
        PostgresRepo[PostgreSQL Repository]
        HashService[Hash Service]
        JWTService[JWT Service]
        Database[(PostgreSQL)]
    end

    Client -->|HTTP Request| CORS
    CORS -->|Add CORS headers| Auth
    CORS -->|Add CORS headers| ClerkAuth
    Auth -->|Validate JWT| Routes
    ClerkAuth -->|Validate Clerk Token| Routes
    Routes -->|Route to| Handlers
    Handlers -->|Business Logic| UserService
    Handlers -->|Authentication| AuthHandler
    Handlers -->|Validation| Validator
    UserService -->|Use| Repository
    AuthHandler -->|Use| Repository
    Repository -->|Implement| PostgresRepo
    UserService -->|Hash passwords| HashService
    AuthHandler -->|Generate tokens| JWTService
    PostgresRepo -->|Query| Database
```

## Endpoints Summary

| Method | Route         | Authentication | Description                            |
| ------ | ------------- | -------------- | -------------------------------------- |
| POST   | `/users`      | ❌ No          | Create user (registration)             |
| POST   | `/auth/login` | ❌ No          | Login and get token                    |
| POST   | `/users/sync` | ✅ Clerk Token | Sync user from Clerk (create/retrieve) |
| GET    | `/users`      | ✅ JWT Token   | List all users                         |
| GET    | `/users/{id}` | ✅ JWT Token   | Get user by ID                         |
| PUT    | `/users/{id}` | ✅ JWT Token   | Update user                            |
| DELETE | `/users/{id}` | ✅ JWT Token   | Delete user                            |
