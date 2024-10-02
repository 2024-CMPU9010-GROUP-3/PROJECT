-- name: CreatePoint :one
INSERT INTO points (
  LongLat, Type, Details
) VALUES (
  $1, $2, $3
) RETURNING Id;

-- name: UpdatePoint :exec
UPDATE points
SET LongLat = $2,
Type = $3,
Details = $4
WHERE Id = $1;

-- name: DeletePoint :exec
DELETE FROM points
WHERE id = $1;
