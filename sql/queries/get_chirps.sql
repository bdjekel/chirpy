-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY created_at;