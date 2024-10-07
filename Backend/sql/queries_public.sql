-- name: GetPointsInEnvelope :many
SELECT Id, LongLat::geometry, Type FROM points
WHERE ST_Intersects(ST_MakeEnvelope(@x1::float, @y1::float, @x2::float, @y2::float, 4326), points.LongLat);

-- name: GetPointDetails :one
SELECT Details::jsonb FROM points
WHERE id = $1 LIMIT 1;

-- name: GetUserDetails :one
SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn
FROM user_details
WHERE Id = $1
LIMIT 1;