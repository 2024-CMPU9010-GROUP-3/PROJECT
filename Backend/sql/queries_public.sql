-- name: GetPointsInEnvelope :many
SELECT Id, LongLat::geometry, Type FROM points
WHERE ST_Intersects(ST_MakeEnvelope(@x1::float, @y1::float, @x2::float, @y2::float, 4326), points.LongLat);

-- name: GetPointDetails :one
SELECT Details::jsonb FROM points
WHERE id = $1 LIMIT 1;

-- name: GetLogin :one
SELECT Id, Username, Email, PasswordHash
FROM logins
WHERE Id = $1
LIMIT 1;

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

-- name: UpdateLogin :exec
UPDATE logins
SET 
  Username = COALESCE($2, Username),
  Email = COALESCE($3, Email),
  PasswordHash = COALESCE(NULLIF(@PasswordHash::VARCHAR(72), ''), PasswordHash)
WHERE Id = $1;

-- name: UpdateUserDetails :exec
UPDATE user_details
SET
  FirstName = COALESCE(NULLIF(@FirstName::VARCHAR(64), ''), FirstName),
  LastName = COALESCE(NULLIF(@LastName::VARCHAR(64), ''), LastName),
  ProfilePicture = COALESCE(NULLIF(@ProfilePicture::VARCHAR(512), ''), ProfilePicture)
WHERE Id = $1;