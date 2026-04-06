-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: EnsureDefaultUser :exec
INSERT IGNORE INTO users (id, name, email, password_hash)
VALUES (1, 'default', 'default@beanmemo.local', 'n/a');
