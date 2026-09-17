# TaskPulse API

Task management REST API built with Go, Gin, PostgreSQL, and JWT authentication.

## Features

- ✅ User authentication (register, login, JWT)
- ✅ Task CRUD operations with ownership checks
- ✅ Pagination and filtering
- ✅ Rate limiting (brute-force protection)
- ✅ CORS support
- ✅ Health check endpoint
- ✅ Structured logging with Zap
- ✅ Unit tests
- ✅ Swagger/OpenAPI docs

## Tech Stack

- **Language:** Go 1.27
- **Framework:** Gin
- **Database:** PostgreSQL
- **Authentication:** JWT + bcrypt
- **Containerization:** Docker & Docker Compose

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.27+

### Installation

```bash
git clone https://github.com/yourusername/TaskPulse.git
cd TaskPulse
```

### Running Locally

```bash
docker compose up --build
```

- API: `http://localhost:8080`
- Adminer (DB UI): `http://localhost:8081`

### Environment Variables

Create `.env`:

```
DATABASE_URL=postgres://postgres:[REDACTED]@db:5432/taskpulse_db?sslmode=disable
JWT_SECRET=your-secret-key
```

## API Endpoints

### Auth

- `POST /api/v1/auth/register` — Register user
- `POST /api/v1/auth/login` — Login & get JWT token

### Users

- `GET /api/v1/users/me` — Get current user (requires auth)
- `GET /api/v1/users` — Get all users

### Tasks

- `GET /api/v1/tasks` — Get user's tasks (pagination + filtering)
- `POST /api/v1/tasks` — Create task
- `GET /api/v1/tasks/:id` — Get specific task
- `PUT /api/v1/tasks/:id` — Update task
- `DELETE /api/v1/tasks/:id` — Delete task

### Monitoring

- `GET /health` — Health check

## Example Requests

### Register

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Create Task

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Task","description":"Description","priority":"high"}'
```

### Get Tasks with Pagination

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  "http://localhost:8080/api/v1/tasks?page=1&limit=10&status=pending"
```

## Running Tests

```bash
go test ./internal/handlers -v
```

## Project Structure

```
cmd/api/
├── main.go

internal/
├── handlers/
│   ├── auth.go
│   ├── task.go
│   ├── health.go
│   ├── helpers.go
│   └── *_test.go
├── models/
├── middleware/
└── logger/

internal/migrations/
```

## Security Features

- JWT token-based authentication
- bcrypt password hashing
- Rate limiting on auth endpoints
- Ownership checks on resources
- CORS middleware
- SQL injection prevention

## Future Improvements

- [ ] Comprehensive test coverage (>80%)
- [ ] WebSocket notifications
- [ ] Email verification
- [ ] Refresh tokens
- [ ] Role-based access control

## License

MIT
