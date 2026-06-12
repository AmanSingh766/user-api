# user-api

A simple REST API built with Go to manage users. It stores `name` and `dob` (date of birth) and returns the user's age dynamically on fetch.

---

## Tech Stack

- **Go 1.22**
- **GoFiber v2** – HTTP framework
- **PostgreSQL** – database
- **SQLC** – type-safe SQL query generation
- **Uber Zap** – structured logging
- **go-playground/validator** – request validation
- **Docker + docker-compose** – optional containerised setup

---

## What you need installed

### Without Docker

| Tool | Version | Install |
|---|---|---|
| Go | 1.22+ | https://go.dev/dl/ |
| PostgreSQL | 14+ | https://www.postgresql.org/download/ |

### With Docker (recommended)

| Tool | Install |
|---|---|
| Docker Desktop | https://www.docker.com/products/docker-desktop/ |
| docker-compose | included with Docker Desktop |

---

## Running locally (without Docker)

**1. Clone the repo**

```bash
https://github.com/AmanSingh766/user-api.git
cd user-api
```

**2. Set up the database**

Create a PostgreSQL database and run the migration:

```bash
psql -U postgres -c "CREATE DATABASE userdb;"
psql -U postgres -d userdb -f db/migrations/001_create_users_table.up.sql
```

**3. Configure environment**

```bash
cp .env.example .env
# Edit .env with your database credentials if needed
```

**4. Install dependencies**

```bash
go mod tidy
```

**5. Start the server**

```bash
go run ./cmd/server
```

Server will start at `http://localhost:8080`

---

## Running with Docker

```bash
docker-compose up --build
```

This spins up PostgreSQL and the API together. The migration runs automatically via the `docker-entrypoint-initdb.d` volume mount.

Server available at `http://localhost:8080`

---

## API Endpoints

### POST /users – Create a user

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "dob": "1990-05-10"}'
```

Response `201`:
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10"}
```

---

### GET /users/:id – Get a user (with age)

```bash
curl http://localhost:8080/users/1
```

Response `200`:
```json
{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35}
```

---

### PUT /users/:id – Update a user

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Updated", "dob": "1991-03-15"}'
```

Response `200`:
```json
{"id": 1, "name": "Alice Updated", "dob": "1991-03-15"}
```

---

### DELETE /users/:id – Delete a user

```bash
curl -X DELETE http://localhost:8080/users/1
```

Response `204 No Content`

---

### GET /users – List all users (paginated)

```bash
curl "http://localhost:8080/users?page=1&limit=10"
```

Response `200`:
```json
{
  "data": [{"id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35}],
  "total": 1,
  "page": 1,
  "limit": 10,
  "total_pages": 1
}
```

---

## Running tests

```bash
go test ./internal/service/... -v
```

This runs the unit tests for the age calculation logic.

---

## Project Structure

```
/cmd/server/main.go         - entrypoint
/config/                    - environment config loader
/db/migrations/             - SQL migration files
/db/sqlc/                   - SQLC generated + query files
/internal/
  handler/                  - HTTP handlers (GoFiber)
  repository/               - database access layer
  service/                  - business logic + age calculation
  routes/                   - route registration
  middleware/               - requestID + request logger
  models/                   - request/response structs
  logger/                   - Zap logger setup
```

---

## Notes

- Age is calculated dynamically in Go using the `time` package — it is never stored in the database
- Every response includes an `X-Request-ID` header for tracing
- Request method, path, status code, and duration are logged to stdout via Zap
