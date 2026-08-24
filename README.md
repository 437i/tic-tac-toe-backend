# Tic-Tac-Toe Backend

A learning Go backend project implementing a REST API for PvP and PvE Tic-Tac-Toe with JWT authentication, PostgreSQL, Minimax AI, dependency injection, and automated tests.

## Features

- User registration and authentication
- JWT access and refresh tokens
- Refresh token rotation
- PvP and PvE game modes
- Choice of `X` or `O` when creating a game
- Minimax-based AI opponent
- PvP game matchmaking
- Game history
- Leaderboard
- User and game lookup by UUID
- Unit and integration tests

A lightweight browser UI is also included for manual API testing, including PvP testing in two browser tabs, leaderboard, game history, and HTTP request/response logging.

## Tech Stack

- Go
- PostgreSQL
- REST API
- JWT
- gorilla/mux
- pgx
- Squirrel
- Uber FX
- Goose
- bcrypt
- Docker / Docker Compose

## Architecture

The backend follows a layered architecture with separate HTTP, business logic, and persistence layers:

```text
HTTP
 │
 ▼
web
 │
 ▼
domain
 │
 ▼
datasource
 │
 ▼
PostgreSQL
```

Additional components:

- `auth` — authentication and JWT logic
- `di` — dependency injection and application wiring
- `migrations` — database migrations
- `cmd` — application entry point

## Project Structure

```text
.
├── auth
├── cmd
├── datasource
├── di
├── domain
├── manual-test-UI
├── migrations
├── web
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## API

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/signup` | Register a user |
| `POST` | `/login` | Login and obtain tokens |
| `POST` | `/refresh/access` | Refresh access token |
| `POST` | `/refresh/refresh` | Rotate refresh token and obtain a new token pair |

### Users

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/user/me` | Get current user |
| `GET` | `/user/{id}` | Get user by UUID |

### Games

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/game` | Create a PvP or PvE game |
| `GET` | `/game/available` | Get available PvP games |
| `GET` | `/game/history` | Get completed games |
| `GET` | `/game/{id}` | Get game by UUID |
| `POST` | `/game/{id}/join` | Join a PvP game |
| `POST` | `/game/{id}` | Make a move |

### Leaderboard

```text
GET /leaderboard?limit=10
```

Returns the top players with their win statistics.

Protected endpoints require:

```text
Authorization: Bearer <access_token>
```

## Running locally

### Requirements

- Go
- Docker
- Docker Compose
- Goose

### 1. Configure environment variables

Create a `.env` file based on `.env.example`.

```env
PG_HOST=localhost
PG_PORT=5432
PG_USER=
PG_PASSWORD=
PG_DB=

GOOSE_DRIVER=postgres
GOOSE_MIGRATION_DIR=./migrations
GOOSE_DBSTRING=

ACCESS=
ACCESS_EXP=5m
REFRESH=
REFRESH_EXP=24h
```

Do not commit real secrets to the repository.

### 2. Start PostgreSQL

```bash
make db-up
```

### 3. Apply migrations

```bash
make migrate-up
```

### 4. Start the backend

```bash
make run
```

## Testing

Run all tests with:

```bash
make test
```

The test suite covers game logic, user services, and PostgreSQL repositories.

## Manual Test UI

The repository includes a small browser UI for manual API testing.

It supports:

- authentication and token management;
- PvE and PvP games;
- interactive game board;
- PvP testing in multiple browser tabs;
- game history and available games;
- leaderboard;
- user lookup;
- HTTP request/response logging.

Start the backend first, then run the UI:

```bash
make ui
```

Open:

```text
http://127.0.0.1:3000
```

The UI connects to the Go backend through a small Python reverse proxy because the API does not currently enable CORS.
