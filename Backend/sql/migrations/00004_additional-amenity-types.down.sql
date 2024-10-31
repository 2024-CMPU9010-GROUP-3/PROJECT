BEGIN;

-- store types in json to prevent data loss as much as possible
UPDATE points
SET 
    details = jsonb_set(details, '{_type}', to_jsonb(type::text), true),
    type = 'unknown'
WHERE type IN (
    'coach_parking',
    'bike_sharing_station',
    'bike_stand',
    'drinking_water_fountain',
    'public_toilet',
    'accessible_parking',
    'public_wifi_access_point',
    'library',
    'multistorey_car_parking',
    'parking_meter',
    'public_bins'
);

ALTER TYPE point_type RENAME TO point_type_old;

CREATE TYPE point_type AS ENUM ('parking', 'unknown');

ALTER TABLE points
ALTER COLUMN type TYPE point_type USING type::text::point_type;

DROP TYPE point_type_old;

COMMIT;