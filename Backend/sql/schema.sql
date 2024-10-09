CREATE TYPE point_type AS ENUM ('placeholder1', 'placeholder2');

CREATE TABLE points (
  Id BIGSERIAL PRIMARY KEY,
  LongLat geometry(Point, 4326) NOT NULL,
  Type point_type NOT NULL,
  Details JSONB
);

CREATE TABLE logins (
  Id UUID PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  Username VARCHAR(64) UNIQUE NOT NULL,
  Email VARCHAR(64) UNIQUE NOT NULL,
  PasswordHash VARCHAR(72) NOT NULL
);

CREATE TABLE user_details (
  Id UUID PRIMARY KEY NOT NULL REFERENCES logins (Id) ON DELETE CASCADE,
  RegisterDate TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
  FirstName VARCHAR(64) NOT NULL,
  LastName VARCHAR(64) NOT NULL,
  ProfilePicture VARCHAR(512),
  LastLoggedIn TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);