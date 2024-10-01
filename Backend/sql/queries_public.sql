-- name: GetPointByRadius :one
SELECT Id, LongLat, Type FROM points
WHERE ST_Contains(
  ST_MakeEnvelope(
    ST_X($1::GEOMETRY) - $2::DECIMAL,
    ST_Y($1::GEOMETRY) - $2::DECIMAL,
    ST_X($1::GEOMETRY) + $2::DECIMAL,
    ST_Y($1::GEOMETRY) + $2::DECIMAL,
    ST_SRID($1::GEOMETRY)
  ), points.LongLat
);

-- name: GetPointDetails :one
SELECT Details FROM points
WHERE id = $1 LIMIT 1;
