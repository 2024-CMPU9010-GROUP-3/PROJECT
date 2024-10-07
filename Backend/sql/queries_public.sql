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

-- name: CreateUser :one
INSERT INTO logins (
  Username, Email, PasswordHash
) VALUES (
  $1, $2, $3
) RETURNING Id;

-- name: CreateUserDetails :one
INSERT INTO user_details (
  Id, FirstName, LastName, ProfilePicture
) VALUES (
  $1, $2, $3, $4
) RETURNING Id;