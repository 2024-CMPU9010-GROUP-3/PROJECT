CREATE TABLE location_history (
  Id BIGSERIAL PRIMARY KEY,
  UserId UUID NOT NULL REFERENCES logins (Id) ON DELETE CASCADE,
  DateCreated TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
  AmenityTypes point_type[],
  LongLat geometry(Point, 4326) NOT NULL,
  Radius int NOT NULL
);