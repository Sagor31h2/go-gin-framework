# Progress

Last updated: 2026-05-21

## Status

| Phase | Title | Status |
|-------|-------|--------|
| 1 | User Auth | ✅ Done |
| 2 | Notes CRUD (auth-scoped) | ✅ Done |
| 3 | Request Handling | ✅ Done |
| 4 | Middleware | ✅ Done |
| 5 | Note Sharing | ✅ Done |
| 6 | Route Organization | ✅ Done |
| 7 | Central Error Handling | ✅ Done |
| 8 | File Attachments | ⬜ Not started |
| 9 | Testing | ⬜ Not started |
| 10 | Config & Production Readiness | ⬜ Not started |

## Phase 1 — Done

- `User` model: id, email, password_hash, created_at, updated_at
- `POST /auth/register` — bcrypt hash, return user (201)
- `POST /auth/login` — verify password, return JWT (24h expiry)
- JWT secret from `JWT_SECRET` env, fallback `"secret"`
- `AutoMigrate` for User + Note in main.go

Files:
- `internal/models/userModel.go`
- `services/authService.go`
- `controllers/auth.go`
- `main.go` (updated)

## Phase 2 — Done

Tasks:
- [x] Add `user_id` to `Note` model
- [x] `GET /notes/` — auth scoped
- [x] `POST /notes/` — auth scoped
- [x] `GET /notes/:id` — auth scoped
- [x] `PUT /notes/:id` — auth scoped
- [x] `DELETE /notes/:id` — auth scoped
- [x] Scoping all queries to user_id

## Phase 3 — Done

Tasks:
- [x] Implement search query param `?search=foo`
- [x] Implement pagination `?limit=10&page=1`
- [x] Request handling logic in controller/service

## Phase 4 — Done

Tasks:
- [x] JWT auth middleware
- [ ] Ownership guard middleware (redundant with scoped queries?)
- [ ] Custom request logger

## Phase 6 — Done (Refactor)

Tasks:
- [x] Centralized route organization in `routes/routes.go`
- [x] Version prefix `/api/v1`
- [x] Cleaned up `main.go` initialization

## Next Up

→ Phase 8 (File Attachments)

## Phase 7 — Done

- `AppError` type with `Status` + `Message` in `internal/errors/errors.go`
- Helpers: `NotFound`, `Forbidden`, `BadRequest`, `Conflict`, `Internal`, `Unauthorized`
- `ErrorHandler` middleware: reads `c.Errors.Last()`, type-asserts to `*AppError`, writes JSON
- All services return typed `*AppError` instead of raw errors
- All controllers use `c.Error(err)` + `ctx.Abort()` — no scattered `ctx.JSON` error calls
- `ErrorHandler` registered globally in `routes.go`

Files:
- `internal/errors/errors.go`
- `middleware/errorHandler.go`
- `services/notesService.go` (updated)
- `services/authService.go` (updated)
- `controllers/notes.go` (updated)
- `controllers/auth.go` (updated)
- `controllers/share.go` (updated)
- `routes/routes.go` (updated)

## Phase 5 — Done

- `NoteShare` model: note_id, shared_with_user, created_at
- `POST /api/v1/notes/:id/share` — owner shares with user
- `DELETE /api/v1/notes/:id/share/:uid` — owner revokes share
- `GET /api/v1/notes/:id/shares` — owner lists who has access
- `GET /api/v1/notes` now returns owned + shared notes (subquery)

Files:
- `internal/models/noteShareModel.go`
- `services/notesService.go` (updated)
- `controllers/share.go`
- `routes/routes.go` (updated)
- `main.go` (updated AutoMigrate)
