CREATE TYPE point_type AS ENUM ('placeholder1', 'placeholder2');

CREATE TABLE points (
  Id BIGSERIAL PRIMARY KEY,
  LongLat geometry(Point, 4326) NOT NULL,
  Type point_type,
  Details JSONB
);

CREATE TABLE logins (
  Id UUID DEFAULT (uuid_generate_v4()),
  Username VARCHAR(64),
  Email VARCHAR(64),
  PasswordHash CHAR(72)
);

CREATE TABLE user_details (
  Id UUID PRIMARY KEY REFERENCES logins (Id),
  RegisterDate TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
  FirstName VARCHAR(64),
  LastName VARCHAR(64),
  ProfilePicture VARCHAR(512),
  LastLoggedIn TIMESTAMP
);
