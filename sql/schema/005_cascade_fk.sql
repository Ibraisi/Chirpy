-- +goose Up
ALTER TABLE refresh_tokens
DROP CONSTRAINT fk_user,
ADD CONSTRAINT fk_user
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE refresh_tokens
DROP CONSTRAINT fk_user,
ADD CONSTRAINT fk_user
  FOREIGN KEY (user_id) REFERENCES users(id);
