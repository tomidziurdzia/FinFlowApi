# FinFlow API

API de gestión financiera construida con Go y Domain-Driven Design (DDD).

## 📋 Tabla de Contenidos

- [Arquitectura](#arquitectura)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Bounded Contexts](#bounded-contexts)
- [Capas de la Aplicación](#capas-de-la-aplicación)
- [Uso](#uso)
- [Desarrollo](#desarrollo)

## 🏗️ Arquitectura

Este proyecto sigue los principios de **Domain-Driven Design (DDD)** con una arquitectura por **Bounded Contexts**.

### Principios Aplicados

- **Separación de responsabilidades**: Cada capa tiene una responsabilidad clara
- **Independencia de bounded contexts**: Cada dominio es independiente
- **CQRS**: Separación de Commands (escritura) y Queries (lectura)
- **Dependency Inversion**: Las capas internas no dependen de las externas

## 📁 Estructura del Proyecto

```
FinFlowApi/
├── cmd/
│   └── api/
│       └── main.go              # Punto de entrada (wiring + startup)
│
├── internal/
│   ├── users/                   # BOUNDED CONTEXT: Users
│   │   ├── domain/              # Lógica de negocio
│   │   │   ├── user.go          # Entidad User
│   │   │   ├── repository.go    # Interfaz del repositorio
│   │   │   └── errors.go        # Errores del dominio
│   │   │
│   │   ├── application/         # Casos de uso (CQRS)
│   │   │   ├── contracts/       # DTOs/Contracts
│   │   │   │   ├── commands/  # Request DTOs (Create, Update, Delete)
│   │   │   │   └── queries/    # Request/Response DTOs (Get, List)
│   │   │   ├── commands/        # Handlers de escritura
│   │   │   └── queries/         # Handlers de lectura
│   │   │
│   │   ├── interfaces/http/     # HTTP handlers
│   │   │   ├── handlers.go
│   │   │   ├── routes.go
│   │   │   └── dto_*.go         # DTOs HTTP
│   │   │
│   │   └── infrastructure/      # Implementaciones
│   │       └── persistence/
│   │           ├── memory/      # Repositorio en memoria
│   │           └── postgres/   # Repositorio PostgreSQL (futuro)
│   │
│   ├── shared/                  # Código compartido entre BCs
│   │   ├── domain/              # Entity base (patrón Entity)
│   │   ├── cqrs/                # Interfaces base CQRS
│   │   ├── errors/              # Errores comunes
│   │   ├── interface/           # Interfaces compartidas (JWT, time)
│   │   ├── http/                # Base handler
│   │   └── middleware/          # Middleware compartido
│   │
│   ├── interfaces/http/         # HTTP compartido
│   │   ├── server.go            # Servidor con graceful shutdown
│   │   ├── routes.go            # Orquestador de rutas
│   │   └── health.go            # Health check
│   │
│   ├── infrastructure/          # Servicios técnicos compartidos
│   │   ├── config/             # Configuración
│   │   ├── db/                 # Base de datos
│   │   ├── hash/               # Hashing
│   │   ├── jwt/                # JWT
│   │   └── time_service/       # Servicio de tiempo
│   │
│   └── bootstrap/              # Wiring de dependencias
│       └── wiring.go           # Construcción de dependencias
│
└── go.mod
```

## 🎯 Bounded Contexts

### Users (Implementado)

El bounded context de **Users** maneja toda la lógica relacionada con usuarios:

- **Domain**: Entidad `User` con campos básicos
- **Application**: CRUD completo (Create, Read, Update, Delete, List)
- **Infrastructure**: Repositorio en memoria (listo para PostgreSQL)
- **Interfaces**: Handlers HTTP para exponer la API

### Futuros Bounded Contexts

- **Transactions**: Gestión de transacciones financieras
- **Accounts**: Gestión de cuentas bancarias
- **Categories**: Categorización de transacciones
- **Budgets**: Presupuestos y límites

Cada bounded context seguirá la misma estructura que `users/`.

## 🧩 Capas de la Aplicación

### 1. Domain (Dominio)

**Responsabilidad**: Lógica de negocio pura, entidades, value objects.

```go
// internal/users/domain/user.go
type User struct {
    domain.Entity  // Embedding de Entity base
    AuthID    string
    FirstName string
    LastName  string
    Email     string
    Password  string
}

func NewUser(...) *User  // Constructor
```

**Características**:

- No depende de otras capas
- Solo contiene lógica de negocio
- Define interfaces de repositorios (no implementaciones)

### 2. Application (Aplicación)

**Responsabilidad**: Casos de uso, orquestación, CQRS.

#### Contracts (DTOs)

```go
// internal/users/application/contracts/commands/create_user_request.go
type CreateUserRequest struct {
    AuthID    string
    FirstName string
    LastName  string
    Email     string
    Password  string
}
```

#### Handlers (Implementación)

```go
// internal/users/application/commands/create_user_handler.go
type CreateUserHandler struct {
    repository domain.UserRepository
}

func (h *CreateUserHandler) Handle(req commands.CreateUserRequest) error {
    user := domain.NewUser(...)
    return h.repository.Create(user)
}
```

**Características**:

- Usa contracts (DTOs) para entrada/salida
- Orquesta el dominio y la infraestructura
- Separa Commands (escritura) y Queries (lectura)

### 3. Infrastructure (Infraestructura)

**Responsabilidad**: Implementaciones técnicas (DB, servicios externos).

```go
// internal/users/infrastructure/persistence/memory/user_repository.go
type Repository struct {
    users map[string]*domain.User
}

func (r *Repository) Create(user *domain.User) error {
    // Implementación en memoria
}
```

**Características**:

- Implementa interfaces definidas en Domain/Application
- Maneja detalles técnicos (DB, APIs externas)
- Puede tener múltiples implementaciones (memory, postgres)

### 4. Interfaces (Presentación)

**Responsabilidad**: Entrada/salida (HTTP, gRPC, CLI).

```go
// internal/users/interfaces/http/handlers.go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    json.NewDecoder(r.Body).Decode(&req)

    cmd := commands.CreateUserRequest{...}
    handler.Handle(cmd)
}
```

**Características**:

- Convierte HTTP Request → Application Contract
- Convierte Application Response → HTTP Response
- Maneja errores HTTP

## 🔄 Flujo de una Petición

```
1. HTTP Request
   ↓
2. HTTP Handler (interfaces/http)
   - Convierte Request → Contract
   ↓
3. Application Handler (application/commands o queries)
   - Valida
   - Usa Domain (NewUser, etc.)
   - Llama Repository
   ↓
4. Domain (domain/)
   - Lógica de negocio
   ↓
5. Infrastructure (infrastructure/persistence)
   - Persiste en DB/memoria
   ↓
6. Response
   - Domain → Application → HTTP → Client
```

## 📝 Ejemplo de Uso

### Crear un Usuario

```go
// 1. Crear handler
handler := commands.NewCreateUserHandler(repository)

// 2. Usar contract
req := commands.CreateUserRequest{
    AuthID:    "auth-123",
    FirstName: "John",
    LastName:  "Doe",
    Email:     "john@example.com",
    Password:  "password123",
}

// 3. Ejecutar
err := handler.Handle(req)
```

### Obtener un Usuario

```go
// 1. Crear handler
handler := queries.NewGetUserHandler(repository)

// 2. Usar contract
req := queries.GetUserRequest{
    UserID: "user-123",
}

// 3. Ejecutar
user, err := handler.Handle(req)
```

## 🛠️ Desarrollo

### Requisitos

- Go 1.25.5 o superior
- PostgreSQL 12 o superior

### Configuración

1. **Crear archivo `.env`** basado en `.env.example`:

```bash
cp .env.example .env
```

2. **Configurar variables de entorno** en `.env`:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=finflow
DB_SSLMODE=disable
```

3. **Crear la base de datos**:

```bash
createdb finflow
# O usando psql:
psql -U postgres -c "CREATE DATABASE finflow;"
```

4. **Ejecutar migraciones**:

```bash
psql -U postgres -d finflow -f internal/infrastructure/db/migrations/001_create_users_table.sql
```

### Compilar

```bash
go build ./cmd/api
```

### Ejecutar

```bash
./api
# O con variables de entorno explícitas:
PORT=8080 DB_HOST=localhost DB_USER=postgres DB_PASSWORD=password DB_NAME=finflow ./api
```

### Tests

```bash
go test ./...
```

### Estructura de Tests

Los tests se colocan en el mismo paquete con el sufijo `_test.go`:

```
domain/
├── user.go
└── user_test.go    # Test del mismo paquete
```

## 🎨 Patrones Implementados

### 1. Entity Base (Embedding)

```go
// shared/domain/entity.go
type Entity struct {
    ID        string
    CreatedAt time.Time
    ModifiedAt time.Time
    CreatedBy string
    ModifiedBy string
}

// users/domain/user.go
type User struct {
    domain.Entity  // Embedding - similar a herencia
    // ... campos específicos
}
```

### 2. CQRS (Command Query Responsibility Segregation)

- **Commands**: Modifican estado (Create, Update, Delete)
- **Queries**: Solo leen (Get, List)

### 3. Repository Pattern

- Interfaz en Domain
- Implementación en Infrastructure
- Permite cambiar la persistencia sin afectar el dominio

### 4. Dependency Injection

- Handlers reciben dependencias por constructor
- Facilita testing y mantenimiento

## 🚀 Próximos Pasos

1. ✅ Repositorio PostgreSQL implementado
2. Conectar handlers HTTP con application layer
3. Agregar validaciones
4. Implementar autenticación/autorización
5. Agregar más bounded contexts (Transactions, Accounts, etc.)
6. Implementar sistema de migraciones automático

## 📚 Referencias

- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Go Best Practices](https://go.dev/doc/effective_go)

## 📄 Licencia

[Tu licencia aquí]
