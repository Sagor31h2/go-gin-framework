# Gin Note App

A REST API for managing and sharing notes, built with Go, Gin, GORM, and SQLite.

## Features

- **Authentication**: JWT-based auth with bcrypt password hashing.
- **Notes CRUD**: Create, read, update, and delete notes scoped to authenticated users.
- **Search & Pagination**: Filter notes by title and paginate results (`?search=`, `?page=`, `?limit=`).
- **Sharing**: Share notes with other users and manage access permissions.
- **File Attachments**: Upload and serve files attached to notes.
- **Swagger Docs**: Interactive API documentation with auto-launch in browser.
- **Error Handling & Shutdown**: Centralized JSON error responses and graceful server shutdown.

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database & ORM**: SQLite with [GORM](https://gorm.io/)
- **Auth**: JWT (`golang-jwt/jwt/v5`) and bcrypt (`golang.org/x/crypto`)
- **Docs**: [Swag](https://github.com/swaggo/swag)
- **Config**: [godotenv](https://github.com/joho/godotenv)

## Project Structure

```text
.
├── controllers/        # HTTP handlers for auth, notes, and sharing
├── docs/               # Swagger documentation files
├── internal/
│   ├── config/         # App configuration and environment loading
│   ├── database/       # Database connection setup
│   ├── errors/         # Application error types and definitions
│   └── models/         # GORM data models
├── middleware/         # Auth and error handling middleware
├── routes/             # Route groups and URL definitions
├── services/           # Business logic and unit tests
├── tests/
│   └── integration/    # HTTP integration tests
├── uploads/            # Directory for uploaded attachments
├── .env                # Local environment variables (ignored by git)
├── go.mod              # Go module definition
├── main.go             # Application entrypoint
└── README.md
```

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Git](https://git-scm.com/)

## Getting Started

### 1. Clone the repository

```bash
git clone git@github.com:Sagor31h2/go-gin-framework.git
cd go-gin-framework
```

### 2. Configure environment

Create a `.env` file in the root directory (optional, sensible defaults will be used if omitted):

```env
PORT=8000
DB_PATH=test.db
JWT_SECRET=supersecret
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Run the server

```bash
go run main.go
```

The API will listen on `http://localhost:8000` and automatically open the Swagger docs at `http://localhost:8000/swagger/index.html`.

## API Endpoints

Base path: `/api/v1`

### Auth

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `POST` | `/auth/register` | No | Register a new user (`email`, `password`) |
| `POST` | `/auth/login` | No | Authenticate user and get JWT token |

### Notes

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `GET` | `/notes` | Yes | List notes (supports `?search=`, `?page=`, `?limit=`) |
| `POST` | `/notes` | Yes | Create a new note (`title`, `status`) |
| `GET` | `/notes/:id` | Yes | Get a note by ID |
| `PUT` | `/notes/:id` | Yes | Update a note |
| `DELETE` | `/notes/:id` | Yes | Delete a note |
| `POST` | `/notes/:id/attachments` | Yes | Upload file attachment (`multipart/form-data`) |

### Sharing

| Method | Endpoint | Auth | Description |
|---|---|:---:|---|
| `POST` | `/notes/:id/share` | Yes | Share note with a user (`user_id`) |
| `DELETE` | `/notes/:id/share/:uid` | Yes | Revoke user access to note |
| `GET` | `/notes/:id/shares` | Yes | List users who have access |

## Testing

Run all unit and integration tests:

```bash
go test ./...
```

## Swagger Docs

To regenerate Swagger specs after changing route annotations:

```bash
swag init
```

## Contributing

Contributions are welcome. Feel free to open an issue or submit a pull request.

## License

[MIT](LICENSE)