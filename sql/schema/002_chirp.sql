-- +goose Up
CREATE TABLE chirps(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    body text,
    user_id UUID ,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id)
        REFERENCES users(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE chirps;
