BEGIN;

ALTER TYPE point_type ADD VALUE 'coach_parking';
ALTER TYPE point_type ADD VALUE 'bike_sharing_station';
ALTER TYPE point_type ADD VALUE 'bike_stand';
ALTER TYPE point_type ADD VALUE 'drinking_water_fountain';
ALTER TYPE point_type ADD VALUE 'public_toilet';
ALTER TYPE point_type ADD VALUE 'accessible_parking';
ALTER TYPE point_type ADD VALUE 'public_wifi_access_point';
ALTER TYPE point_type ADD VALUE 'library';
ALTER TYPE point_type ADD VALUE 'multistorey_car_parking';
ALTER TYPE point_type ADD VALUE 'parking_meter';
ALTER TYPE point_type ADD VALUE 'public_bins';

COMMIT;

BEGIN;

-- extract types if previously stored in json by down migration
UPDATE points
SET 
    type = (details->>'_type')::point_type,  -- Explicitly cast to point_type
    details = details - '_type'
WHERE details ? '_type'
  AND details->>'_type' IN (
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

COMMIT;