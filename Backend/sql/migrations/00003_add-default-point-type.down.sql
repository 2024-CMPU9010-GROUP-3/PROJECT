-- WARNING: THIS MIGRATION WILL RESULT IN DATA LOSS IF EXECUTED

-- The only other option would be to set the value to "parking" which would falsify the data
DELETE FROM points WHERE type = 'unknown';

ALTER TYPE point_type RENAME TO point_type_old;
CREATE TYPE point_type AS ENUM('parking');
ALTER TABLE points ALTER COLUMN "type" TYPE point_type USING type::text::point_type;
DROP TYPE point_type_old;