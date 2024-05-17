-- name: ListGophers :many
SELECT g.* FROM gophers AS g
ORDER BY g.name;

-- name: GetGopher :one
SELECT g.* FROM gophers AS g
WHERE g.id = ?
LIMIT 1;

-- name: CreateGopher :execresult
INSERT INTO gophers (
    name, job
) VALUES (
    ?, ?
);

-- name: DeleteGopher :exec
DELETE FROM gophers
WHERE id = ?;