CREATE TABLE history_amenity_counts (
  HistoryEntryId BIGSERIAL REFERENCES location_history(Id) ON DELETE CASCADE,
  Type point_type NOT NULL,
  Count int NOT NULL,
  CONSTRAINT PK_history_amenity_counts PRIMARY KEY (HistoryEntryId, Type)
);

ALTER TABLE location_history
DROP COLUMN IF EXISTS AmenityTypes;