-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3,
    $4
);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeToken :exec
  UPDATE refresh_tokens
  SET
      revoked_at = Now(),
      updated_at = NOW()
  WHERE token = $1;
