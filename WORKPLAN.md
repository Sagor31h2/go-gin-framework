# Gin Framework Learning Work Plan

## Domain
Users own notes. Notes can be shared with other users.

## Data Model
```
User ──< Note (owner)
Note ──< NoteShare >── User (shared_with)
```

## Routes

```
POST   /auth/register
POST   /auth/login

GET    /api/v1/notes                      ← my notes + shared with me
POST   /api/v1/notes
GET    /api/v1/notes/:id
PUT    /api/v1/notes/:id                  ← owner only
DELETE /api/v1/notes/:id                  ← owner only

POST   /api/v1/notes/:id/share            ← share with user
DELETE /api/v1/notes/:id/share/:uid       ← revoke share
GET    /api/v1/notes/:id/shares           ← list who has access
```

---

## Phase 1 — User Auth
- `User` model (id, email, password_hash, timestamps)
- `POST /auth/register` — bcrypt hash password
- `POST /auth/login` — verify password, return JWT

Concepts:
- GORM model basics, `AutoMigrate`
- `golang.org/x/crypto/bcrypt`
- `golang-jwt/jwt` — sign + return token

---

## Phase 2 — Notes CRUD (auth-scoped)
- Full CRUD on `/api/v1/notes`
- All queries scoped to `user_id` from JWT

Concepts:
- Path params: `c.Param("id")`
- JWT middleware: extract user, `c.Set("user", ...)`
- GORM: `Where("user_id = ?", uid)`, `First`, `Save`, `Delete`
- HTTP status codes (200, 201, 204, 404)

---

## Phase 3 — Request Handling
- `GET /api/v1/notes?search=foo&limit=10&page=1`

Concepts:
- `c.Query`, `c.DefaultQuery`, `ShouldBindQuery`
- Struct validation tags: `binding:"required,min=3"`
- go-playground/validator custom rules

---

## Phase 4 — Middleware
- JWT auth middleware (protect `/api/v1` group)
- Ownership guard middleware (note belongs to user)
- Custom request logger

Concepts:
- `c.Next()`, `c.Abort()`, `c.AbortWithStatusJSON`
- `c.Set` / `c.Get` — pass data middleware → handler
- Middleware scoping per route group

---

## Phase 5 — Note Sharing
- `NoteShare` join table (note_id, shared_with_user_id, permission)
- `POST /api/v1/notes/:id/share`
- `DELETE /api/v1/notes/:id/share/:uid`
- `GET /api/v1/notes/:id/shares`
- `GET /api/v1/notes` returns owned + shared notes

Concepts:
- GORM `HasMany`, `BelongsTo`, `Many2Many`
- `Preload`, joins
- Scoped queries across two relations

---

## Phase 6 — Route Organization
- Version prefix `/api/v1`
- Groups: `auth` (public), `notes` (JWT-protected)
- Router setup in own package

Concepts:
- `router.Group`
- Nested groups with different middleware stacks

---

## Phase 7 — Central Error Handling
- Custom error types (NotFound, Forbidden, Validation)
- Single error handler middleware
- Consistent JSON error shape

Concepts:
- `c.Error()`
- Error middleware pattern
- Avoid scattered `ctx.JSON` error calls

---

## Phase 8 — File Attachments
- `POST /api/v1/notes/:id/attachment`
- Serve uploaded files statically

Concepts:
- `c.FormFile`, `c.SaveUploadedFile`
- `router.Static`

---

## Phase 9 — Testing
- Unit test service layer (mock DB)
- Integration test handlers

Concepts:
- `gin.SetMode(gin.TestMode)`
- `httptest.NewRecorder`
- Table-driven tests

---

## Phase 10 — Config & Production Readiness
- Env-based config (`godotenv` or `viper`)
- Graceful shutdown
- Switch SQLite → PostgreSQL

Concepts:
- `context` with timeout
- `signal.NotifyContext`
- `gin.SetMode(gin.ReleaseMode)`
- GORM Postgres driver: `gorm.io/driver/postgres`
