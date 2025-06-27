-- name: CreateFile :one
INSERT INTO files (filename, content, embedding)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFile :one
SELECT * FROM files WHERE id = $1;

-- name: UpdateFile :one
UPDATE files
  SET filename = $2, content = $3, embedding = $4
WHERE id = $1
RETURNING *;

-- name: DeleteFile :exec
DELETE FROM files WHERE id = $1;

-- name: SearchFiles :many
SELECT *, embedding <-> $1 AS distance
FROM files
ORDER BY embedding <-> $1
LIMIT $2;
