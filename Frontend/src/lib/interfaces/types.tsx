import { ColumnDef } from "@tanstack/react-table";

export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface CoordinatesForGeoJson {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  coordinates: [number, number] | any;
}

// Interface for distance scales within the viewport
export interface DistanceScales {
  unitsPerMeter: number[];
  metersPerUnit: number[];
  unitsPerDegree: number[];
  degreesPerUnit: number[];
}

// Type for the GeoJsonCollection
export type GeoJsonCollection =
  | "parking_meter"
  | "bike_stand"
  | "public_wifi_access_point"
  | "library"
  | "multistorey_car_parking"
  | "drinking_water_fountain"
  | "public_toilet"
  | "bike_sharing_station"
  | "parking"
  | "accessible_parking"
  | "public_bins"
  | "coach_parking";

// Interface for the viewport details
export interface Viewport {
  _frustumPlanes: Record<string, unknown>; // Placeholder for frustum planes
  id: string;
  x: number;
  y: number;
  width: number;
  height: number;
  zoom: number;
  distanceScales: DistanceScales;
  focalDistance: number;
  position: [number, number, number];
  modelMatrix: number[] | null;
  isGeospatial: boolean;
  scale: number;
  center: [number, number, number];
  viewMatrixUncentered: number[];
  viewMatrix: number[];
  projectionMatrix: number[];
  viewProjectionMatrix: number[];
  viewMatrixInverse: number[];
  cameraPosition: [number, number, number];
  pixelProjectionMatrix: number[];
  pixelUnprojectionMatrix: number[];
  latitude: number;
  longitude: number;
  pitch: number;
  bearing: number;
  altitude: number;
  fovy: number;
  orthographic: boolean;
  _subViewports: null | unknown; // Placeholder for sub-viewports
  _pseudoMeters: boolean;
}

// Main interface for the click event
export interface MapClickEvent {
  color: Uint8Array | null; // Updated to match DeckGL's expected type
  layer: string | null;
  viewport: Viewport;
  index: number;
  picked: boolean;
  x: number;
  y: number;
  pixel: [number, number]; // Pixel coordinates
  coordinate: [number, number]; // Geographic coordinates (longitude, latitude)
  pixelRatio: number;
}

export interface MapHoverEvent {
  color: Uint8Array | null; // Updated to match DeckGL's expected type
  layer: string | null;
  viewport: Viewport;
  index: number;
  picked: boolean;
  x: number;
  y: number;
  pixel: [number, number]; // Pixel coordinates
  coordinate: [number, number]; // Geographic coordinates (longitude, latitude)
  pixelRatio: number;
}

interface Location {
  coordinates: [number, number];
  type: string;
}
export interface ImageConfig {
  id: string;
  path: string;
}
interface Amenities {
  parkingMeters?: number;
  bikeStand?: number;
  publicWiFi?: number;
  library?: number;
  multiStoreyCar?: number;
  drinkingWater?: number;
  publicToilet?: number;
  bikeStationSharing?: number;
  parking?: number;
  accessibleParking?: number;
  publicBins?: number;
  coachingPark?: number;
}
export interface LocationItem {
  id: number;
  date: string;
  amenities: Amenities;
  location: string;
  coordinates: Coordinates;
}

export interface LocationData {
  id: number;
  datecreated: string;
  amenitytypes: string[];
  longlat: {
    type: string;
    coordinates: number[];
  };
  radius: number;
}

export interface Data {
  Id: number;
  Longlat: Location;
  Type: string;
}

export interface GeoJson {
  type: "FeatureCollection";
  features: {
    type: "Feature";
    geometry: {
      type: "Point";
      coordinates: [number, number];
    };
    properties: {
      Id: number;
      Type: string;
    };
  }[];
}

export type Point = {
  Id: string;
  Type: string;
  Longlat: Coordinates | CoordinatesForGeoJson;
};

export interface ImageConfig {
  value: string;
  id: string;
  path: string;
}

export interface DataTableProps<TData> {
  columns: ColumnDef<TData>[];
  data: TData[];
  rowSelection?: Record<string, boolean>;
  setRowSelection: (value: Record<string, boolean>) => void;
}
