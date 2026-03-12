-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps  ORDER BY created_at ASC;

-- name: GetAllChirpsByUser :many
SELECT * FROM chirps where user_id = $1 ORDER BY created_at ASC;

-- name: GetAllChirpsByID :one
SELECT * FROM chirps where id = $1 ORDER BY created_at ASC;
