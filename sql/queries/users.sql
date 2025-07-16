-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3

)
RETURNING *;
-- name: DeleteAllUsers :exec
DELETE FROM users;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
FROM users
WHERE email = $1;


-- name: UpdateUserCredentials :exec
UPDATE users
SET 
    email = $1,
    hashed_password = $2,
    updated_at = Now()
WHERE id = $3;

-- name: UpgradeUserToChirpyRed :exec
UPDATE users
SET 
    is_chirpy_red = $1
WHERE ID = $2;