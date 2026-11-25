-- name: CreateOTP :one
INSERT INTO otps (phone, code_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CountRecentOTPs :one
SELECT COUNT(*) FROM otps
WHERE phone = $1 AND created_at > NOW() - INTERVAL '10 minutes';

-- name: GetLatestOTP :one
SELECT * FROM otps
WHERE phone = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteOTPsByPhone :exec
DELETE FROM otps WHERE phone = $1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: CreateUser :one
INSERT INTO users (phone, role)
VALUES ($1, $2)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
