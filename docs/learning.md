# Learning Outcomes

What I learned implementing each part of Chirpy.

---

## Routing — `net/http` ServeMux

Go 1.22 added method + path pattern matching directly in `ServeMux` (`GET /api/chirps/{chirpID}`), eliminating the need for a router library. Path values are extracted with `r.PathValue("chirpID")`.

---

## Middleware

Middleware is just a function that wraps an `http.HandlerFunc` and returns an `http.Handler`. Chaining is done manually at registration time.

- **Logging** — wraps every handler, logs method + path
- **HitsMiddleware** — counts file server requests using `atomic.Int32`
- **AuthMiddleware** — validates JWT, injects user ID into context

Key insight: storing values in `context.Context` (with a typed key) is how you pass auth state through the handler chain without changing function signatures.

---

## Password Hashing — `POST /api/users`, `POST /api/login`

`bcrypt.GenerateFromPassword` hashes passwords before storing. `bcrypt.CompareHashAndPassword` validates on login. Plain passwords never touch the database.

---

## JWT Authentication — `POST /api/login`, `PUT /api/users`, etc.

JWTs are signed with HMAC-SHA256. The token encodes the user ID as the subject and has a 1-hour expiry. The server validates the signature and expiry on every authenticated request. The secret key lives in `.env`, never in code.

---

## Refresh Tokens — `POST /api/refresh`, `POST /api/revoke`

JWTs are short-lived (1h). Refresh tokens are long-lived (60 days), stored in the database with an expiry and a `revoked_at` column. On refresh, the server checks the token exists, is not expired, and is not revoked, then issues a new JWT.

Logout = revoking the refresh token.

---

## Database with sqlc + postgres

SQL queries are written by hand in `.sql` files, and `sqlc` generates type-safe Go functions from them. No ORM. The generated code in `internal/database/` is not edited directly.

Migrations are plain `.sql` files in `sql/schema/` run in order.

---

## File Server + Hit Counter — `GET /app/`

`http.FileServer` serves static files. Wrapping it in a middleware that increments an `atomic.Int32` on each request demonstrates how middleware works with non-`HandlerFunc` handlers using `http.StripPrefix`.

---

## Query Parameters — `GET /api/chirps`

`r.URL.Query().Get("key")` reads query params. Used for `author_id` filtering and `sort=desc` to reverse results with `slices.Reverse`.

---

## Webhook + API Key Auth — `POST /api/polka/webhooks`

External services authenticate with an API key in the `Authorization: ApiKey <key>` header instead of a JWT. The handler ignores unknown events (`204`) and only acts on `user.upgraded` — a pattern for forward-compatible webhook handling.

---

## Input Validation & Sanitization — `POST /api/chirps`

Enforcing a max body length (400 chars) and replacing profane words with `****` before storing. Simple string processing with `strings.Fields` and a lookup map.
