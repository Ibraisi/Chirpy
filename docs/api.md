# API Reference

Base URL: `http://localhost:8080`

Auth: `Authorization: Bearer <token>`

---

## Health

### GET /api/healthz
Returns `OK` if the server is running.

---

## Users

### POST /api/users
Create a new user.

**Body**
```json
{ "email": "user@example.com", "password": "secret" }
```

**Response** `201`
```json
{ "id": "uuid", "email": "user@example.com", "is_chirpy_red": false, "created_at": "", "updated_at": "" }
```

---

### PUT /api/users
Update email and password for the authenticated user.

**Auth required**

**Body**
```json
{ "email": "new@example.com", "password": "newpass" }
```

**Response** `200` — updated user object

---

## Auth

### POST /api/login
Authenticate and receive tokens.

**Body**
```json
{ "email": "user@example.com", "password": "secret" }
```

**Response** `200`
```json
{
  "id": "uuid", "email": "...", "is_chirpy_red": false,
  "token": "<jwt>",
  "refresh_token": "<refresh>"
}
```

---

### POST /api/refresh
Exchange a refresh token for a new JWT.

**Header** `Authorization: Bearer <refresh_token>`

**Response** `200`
```json
{ "token": "<new_jwt>" }
```

---

### POST /api/revoke
Revoke a refresh token.

**Header** `Authorization: Bearer <refresh_token>`

**Response** `204`

---

## Chirps

### POST /api/chirps
Create a chirp.

**Auth required**

**Body**
```json
{ "body": "Hello world" }
```

Max 400 characters. Profane words are replaced with `****`.

**Response** `201`
```json
{ "id": "uuid", "body": "Hello world", "user_id": "uuid", "created_at": "", "updated_at": "" }
```

---

### GET /api/chirps
List all chirps.

**Query params**
- `author_id=<uuid>` — filter by user
- `sort=desc` — reverse order (default asc)

**Response** `200` — array of chirp objects

---

### GET /api/chirps/{chirpID}
Get a single chirp by ID.

**Response** `200` — chirp object | `404`

---

### DELETE /api/chirps/{chirpID}
Delete a chirp. Only the author can delete their own chirp.

**Auth required**

**Response** `204` | `403` | `404`

---

## Webhooks

### POST /api/polka/webhooks
Webhook from Polka payment provider to upgrade a user to Chirpy Red.

**Header** `Authorization: ApiKey <polka_key>`

**Body**
```json
{ "event": "user.upgraded", "data": { "user_id": "uuid" } }
```

Only `user.upgraded` events are processed; all others return `204`.

**Response** `204`

---

## Admin

### GET /admin/metrics
HTML page showing total file server hit count.

### POST /admin/reset
Delete all users and reset hit counter. Dev environment only.

**Response** `200` | `403`
