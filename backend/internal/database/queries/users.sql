-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: EnsureDefaultUser :exec
INSERT OR IGNORE INTO users (id, name, email, password_hash)
VALUES (1, 'default', 'default@beanmemo.local', 'n/a');

-- name: UpsertUserBySub :exec
INSERT INTO users (sub, name, email, password_hash)
VALUES (?, ?, '', 'n/a')
ON CONFLICT(sub) DO UPDATE SET name = excluded.name;

-- name: GetUserBySub :one
SELECT * FROM users WHERE sub = ? LIMIT 1;
