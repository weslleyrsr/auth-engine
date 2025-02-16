# SSO Project Plan

Below is a **step‐by‐step project plan** laid out as if you had **Jira cards**. You can download or copy/paste this as a `.md` file (e.g., `PROJECT_PLAN.md`) in your repository. Each card includes a description, steps, and file/package guidelines.

---

## Card 1: Project Initialization

**Description**  
Initialize a new Go module and create a basic folder structure to host an SSO project with a layered architecture.

**Steps to Complete**
1. Create an empty GitHub repository (e.g., `github.com/YourOrg/sso-app`).
2. Clone it locally and open it in your IDE.
3. In the root folder, run:
   ```bash
   go mod init github.com/YourOrg/sso-app
   ```
4. Create the initial folder structure:
   ```
   .
   ├── cmd/
   │   └── sso-server/
   │       └── main.go
   ├── internal/
   │   ├── user/
   │   └── token/
   ├── pkg/
   └── go.mod
   ```
5. Add a simple "Hello World" in `cmd/sso-server/main.go`:
   ```go
   package main

   import "fmt"

   func main() {
       fmt.Println("SSO Server starting...")
   }
   ```
6. Commit and push to GitHub.

**File Names / Locations**
- `cmd/sso-server/main.go`: The main entry point.
- `internal/user` and `internal/token`: Packages for your domain logic.

---

## Card 2: Define Domain Models

**Description**  
Create the core domain models (e.g., User, TokenPair) in a “pure” domain package. These should not contain persistence details (DB tags, etc.).

**Steps to Complete**
1. Inside `internal/user/`, create a folder `domain/` and add a file `user.go`.
2. Define a `User` struct with relevant fields:
   ```go
   package domain

   import "github.com/google/uuid"

   type User struct {
       ID       uuid.UUID
       Email    string
       Password string
       Name     string
       ImageURL string
       Website  string
   }
   ```
3. Inside `internal/token/`, create `domain/token.go`.
4. Define a `TokenPair` struct:
   ```go
   package domain

   type TokenPair struct {
       AccessToken  string
       RefreshToken string
   }
   ```
5. Make sure to import `github.com/google/uuid` (or another UUID library) if needed.

**File Names / Locations**
- `internal/user/domain/user.go`
- `internal/token/domain/token.go`

---

## Card 3: User Repository Interface

**Description**  
Create an interface that describes the persistence operations for users. This decouples your business logic from the data store implementation.

**Steps to Complete**
1. In `internal/user/`, create a file `repository.go` or a subfolder `repository/repository.go`.
2. Define the repository interface:
   ```go
   package user

   import (
       "github.com/google/uuid"
       "github.com/YourOrg/sso-app/internal/user/domain"
   )

   type Repository interface {
       Create(user *domain.User) error
       GetByID(id uuid.UUID) (*domain.User, error)
       GetByEmail(email string) (*domain.User, error)
       Update(user *domain.User) error
       // Add more methods if needed (Delete, etc.)
   }
   ```
3. Add any error constants or sentinel errors if you want to standardize them.

**File Names / Locations**
- `internal/user/repository.go`  
  or
- `internal/user/repository/repository.go`

---

## Card 4: Postgres Implementation of User Repository

**Description**  
Implement the `Repository` interface using Postgres as the data store. You’ll likely use `database/sql`, `pgx`, or an ORM like GORM.

**Steps to Complete**
1. In `internal/user/repository/`, create `postgres.go`.
2. Define a struct with a Postgres connection field:
   ```go
   package repository

   import (
       "database/sql"

       "github.com/google/uuid"
       "github.com/YourOrg/sso-app/internal/user"
       "github.com/YourOrg/sso-app/internal/user/domain"
   )

   type PostgresRepository struct {
       db *sql.DB
   }

   func NewPostgresRepository(db *sql.DB) user.Repository {
       return &PostgresRepository{db: db}
   }
   ```
3. Implement the `Create`, `GetByID`, `GetByEmail`, `Update` methods:
   ```go
   func (r *PostgresRepository) Create(u *domain.User) error {
       // INSERT statement
       return nil
   }

   func (r *PostgresRepository) GetByID(id uuid.UUID) (*domain.User, error) {
       // SELECT statement
       return nil, nil
   }

   func (r *PostgresRepository) GetByEmail(email string) (*domain.User, error) {
       // SELECT statement
       return nil, nil
   }

   func (r *PostgresRepository) Update(u *domain.User) error {
       // UPDATE statement
       return nil
   }
   ```
4. For now, just stub out or write minimal code. You can refine queries in a later card.

**File Names / Locations**
- `internal/user/repository/postgres.go`

---

## Card 5: Token Service (JWT Generation)

**Description**  
Implement a service responsible for creating Access and Refresh tokens. This service might also handle invalidation logic if needed.

**Steps to Complete**
1. In `internal/token/service/`, create `token_service.go`.
2. Define an interface if you like:
   ```go
   package service

   import (
       "github.com/YourOrg/sso-app/internal/token/domain"
       "github.com/YourOrg/sso-app/internal/user/domain"
   )

   type Service interface {
       CreateTokenPair(user *domain.User) (*domain.TokenPair, error)
       // Optionally: Invalidate, Validate, etc.
   }
   ```
3. Create a struct that holds signing secrets/config:
   ```go
   type jwtService struct {
       secretKey     []byte
       refreshSecret []byte // optional if you want a separate key
   }

   func NewJWTService(secretKey, refreshSecret []byte) Service {
       return &jwtService{secretKey, refreshSecret}
   }
   ```
4. Implement `CreateTokenPair` to:
   - Create an **AccessToken** (JWT) with a short expiration (e.g., 15m).
   - Create a **RefreshToken** (could be JWT or random string) with longer expiration.
   - Return them in a `TokenPair`.

**File Names / Locations**
- `internal/token/service/token_service.go`

---

## Card 6: User Service

**Description**  
Create a business logic layer for user operations (Signup, Login, Get, Update). This service orchestrates the repository and token service.

**Steps to Complete**
1. In `internal/user/service/`, create `user_service.go`.
2. Define an interface:
   ```go
   package service

   import (
       "github.com/google/uuid"
       "github.com/YourOrg/sso-app/internal/user/domain"
       "github.com/YourOrg/sso-app/internal/token/domain"
   )

   type Service interface {
       Signup(u *domain.User) (*domain.TokenPair, error)
       Login(email, password string) (*domain.TokenPair, error)
       Get(uid uuid.UUID) (*domain.User, error)
       Update(u *domain.User) error
   }
   ```
3. Create a struct that includes a `user.Repository` and a `token.Service`:
   ```go
   import (
       "github.com/YourOrg/sso-app/internal/user"
       "github.com/YourOrg/sso-app/internal/token/service"
   )

   type userService struct {
       repo     user.Repository
       tokenSvc service.Service
   }

   func NewUserService(repo user.Repository, tokenSvc service.Service) Service {
       return &userService{repo: repo, tokenSvc: tokenSvc}
   }
   ```
4. Implement `Signup`:
   - Hash the password (e.g., via `bcrypt`).
   - Create user in DB.
   - Generate token pair from `tokenSvc`.
5. Implement `Login`:
   - Fetch user by email.
   - Compare hashed password.
   - Generate token pair.
6. Implement `Get` and `Update` using repository methods.

**File Names / Locations**
- `internal/user/service/user_service.go`

---

## Card 7: HTTP Handlers (Gin)

**Description**  
Create HTTP endpoints to expose the user operations (signup, login, me, etc.). These handlers will parse/validate input, call the user service, and return JSON responses.

**Steps to Complete**
1. In `internal/user/transport/http/`, create `user_handler.go`.
2. Create a struct:
   ```go
   package http

   import (
       "net/http"

       "github.com/gin-gonic/gin"
       userService "github.com/YourOrg/sso-app/internal/user/service"
   )

   type UserHandler struct {
       userSvc userService.Service
   }

   func NewUserHandler(router *gin.Engine, us userService.Service) {
       h := &UserHandler{userSvc: us}
       router.POST("/signup", h.Signup)
       router.POST("/login", h.Login)
       router.GET("/me", h.Me)
   }
   ```
3. Implement `Signup`:
   ```go
   func (h *UserHandler) Signup(c *gin.Context) {
       var req struct {
           Email    string `json:"email"`
           Password string `json:"password"`
           Name     string `json:"name"`
       }
       if err := c.ShouldBindJSON(&req); err != nil {
           c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
           return
       }

       // Create a domain.User
       // Call h.userSvc.Signup
       // Return TokenPair as JSON
   }
   ```
4. Implement `Login` similarly.
5. Implement `Me`:
   - Extract user info from the `Authorization` header (JWT).
   - Validate token, get user ID, call `userSvc.Get`.
   - Return user data.

**File Names / Locations**
- `internal/user/transport/http/user_handler.go`

---

## Card 8: Main Application Wiring

**Description**  
Wire everything together in `cmd/sso-server/main.go`: database connection, services, and routes.

**Steps to Complete**
1. In `cmd/sso-server/main.go`, add a database connection (Postgres) and wire your repository/service:
   ```go
   package main

   import (
       "database/sql"
       "log"

       "github.com/gin-gonic/gin"
       _ "github.com/lib/pq"

       "github.com/YourOrg/sso-app/internal/user/repository"
       userService "github.com/YourOrg/sso-app/internal/user/service"
       tokenService "github.com/YourOrg/sso-app/internal/token/service"
       userHttp "github.com/YourOrg/sso-app/internal/user/transport/http"
   )

   func main() {
       db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
       if err != nil {
           log.Fatal(err)
       }
       defer db.Close()

       userRepo := repository.NewPostgresRepository(db)
       tSvc := tokenService.NewJWTService([]byte("access-secret"), []byte("refresh-secret"))
       uSvc := userService.NewUserService(userRepo, tSvc)

       r := gin.Default()
       userHttp.NewUserHandler(r, uSvc)

       if err := r.Run(":8080"); err != nil {
           log.Fatal(err)
       }
   }
   ```
2. Ensure you have environment variables for secrets and DB credentials (could use Viper or similar).

**File Names / Locations**
- `cmd/sso-server/main.go`

---

## Card 9: Redis Integration (Refresh Token Storage / Session)

**Description**  
Enhance token invalidation: store refresh tokens or sessions in Redis so you can quickly invalidate tokens or log out users from all devices.

**Steps to Complete**
1. Add a new package in `internal/token/repository/redis` for storing tokens in Redis (optional, but recommended for scalability).
2. Modify `token.Service` to store each RefreshToken in Redis after creation (e.g., store `userID -> list of refresh tokens`).
3. Add an invalidation function to remove the refresh token from Redis on logout.
4. Update `userService.Logout` or create a new method for token invalidation.

**File Names / Locations**
- `internal/token/repository/redis/redis_repository.go`

---

## Card 10: Docker & docker-compose

**Description**  
Containerize the SSO server along with Postgres and Redis for local development and easy deployment.

**Steps to Complete**
1. In the root directory, create a `Dockerfile`:
   ```dockerfile
   FROM golang:1.20-alpine
   WORKDIR /app
   COPY . .
   RUN go mod download
   RUN go build -o sso-server ./cmd/sso-server
   EXPOSE 8080
   CMD ["./sso-server"]
   ```
2. Create a `docker-compose.yml` to run the service plus Postgres and Redis:
   ```yaml
   version: "3.8"
   services:
     sso-server:
       build: .
       ports:
         - "8080:8080"
       depends_on:
         - postgres
         - redis

     postgres:
       image: postgres:15
       environment:
         POSTGRES_PASSWORD: yourpassword
       ports:
         - "5432:5432"

     redis:
       image: redis:7
       ports:
         - "6379:6379"
   ```
3. Test by running:
   ```bash
   docker-compose up --build
   ```
4. Verify the server logs to ensure successful startup.

---

## Card 11: Unit & Integration Tests

**Description**  
Ensure correctness through automated tests. Write both unit tests (mock repositories) and integration tests (with a real DB).

**Steps to Complete**
1. **Unit tests**:
   - Use the standard `testing` package.
   - Mock the repository and verify service logic (e.g., `internal/user/service/user_service_test.go`).
2. **Integration tests**:
   - Spin up Postgres and Redis in a test environment (could use Docker test containers).
   - Insert a user, call `Signup`, verify DB records, etc.
3. Run:
   ```bash
   go test ./...
   ```
   to test all packages.

**File Names / Locations**
- `internal/user/service/user_service_test.go`
- `internal/token/service/token_service_test.go`
- Possibly a `test/integration/` folder for integration tests.

---

## Card 12: Security Hardening & Observability

**Description**  
Add final touches for production readiness: environment variables for secrets, structured logging, metrics, etc.

**Steps to Complete**
1. **Secrets**: Retrieve from environment or a secure vault (e.g., `ACCESS_SECRET`, `REFRESH_SECRET`).
2. **Logging**: Use a logging library like `uber-go/zap` or `sirupsen/logrus`.
3. **Metrics**: Add a Prometheus endpoint or use a Gin middleware (e.g., `ginprometheus`).
4. **TLS**: Ensure you run the server behind HTTPS in production.
5. Document how to configure environment variables or secrets.

---

## (Optional) Additional Cards

- **gRPC Endpoints**: If you plan to expose gRPC, create a `transport/grpc` package.
- **OAuth2 Provider**: Integrate with external IdPs or become an IdP yourself, handling OAuth2 flows.
- **CI/CD**: Set up GitHub Actions or another CI pipeline to build, test, and push Docker images automatically.

---

### Final Thoughts

By following these 12 cards, you’ll build a **scalable SSO application** in Go that is easy to test, maintain, and scale. Customize each step to suit your organization’s needs, and enjoy a clean, modular design that’s ready for production.
