import { LayerProps } from "react-map-gl";

// Updating all icon sizes with larger values
const commonIconLayout = {
  small: ["interpolate", ["linear"], ["zoom"], 12, 0.6, 16, 1.0] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
  medium: ["interpolate", ["linear"], ["zoom"], 12, 0.8, 16, 1.2] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
  large: ["interpolate", ["linear"], ["zoom"], 12, 1.0, 16, 1.4] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
};

export const parkingMeterLayer: LayerProps = {
  id: "parking-meters",
  type: "symbol",
  layout: {
    "icon-image": "parking", // Using Mapbox's built-in parking icon
    "icon-size": commonIconLayout.medium,
    "icon-allow-overlap": false,
    "icon-ignore-placement": false,
    visibility: "visible",
    // Optional text label
    "text-field": ["get", "name"],
    "text-size": 12,
    "text-offset": [0, 1.5],
    "text-optional": true,
  },
  paint: {
    "text-color": "#4264FB",
    "text-halo-color": "#ffffff",
    "text-halo-width": 1,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const bikeStandLayer: LayerProps = {
  id: "bike-stands",
  type: "symbol",
  layout: {
    "icon-image": "bicycle", // Using Mapbox's built-in bicycle icon
    "icon-size": commonIconLayout.medium,
    "icon-allow-overlap": false,
    "icon-ignore-placement": false,
    visibility: "visible",
    // Optional text label
    "text-field": ["get", "name"],
    "text-size": 12,
    "text-offset": [0, 1.5],
    "text-optional": true,
  },
  paint: {
    "text-color": "#FFAA00", // Contrasting color
    "text-halo-color": "#ffffff",
    "text-halo-width": 1,
  },
};

// Bike Sharing Station Layer
export const bikeSharingLayer: LayerProps = {
  id: "bike-sharing",
  type: "symbol",
  layout: {
    "icon-image": "bicycle", // Using Mapbox's built-in bicycle icon
    "icon-size": commonIconLayout.medium,
    "icon-allow-overlap": false,
    "icon-ignore-placement": false,
    visibility: "visible",
    // Optional text label
    "text-field": ["get", "name"],
    "text-size": 12,
    "text-offset": [0, 1.5],
    "text-optional": true,
  },
  paint: {
    "text-color": "#FF6B6B",
    "text-halo-color": "#ffffff",
    "text-halo-width": 1,
  },
};

// Accessible Parking Layer
export const accessibleParkingLayer: LayerProps = {
  id: "accessible-parking",
  type: "symbol",
  layout: {
    "icon-image": "disabled",
    "icon-size": commonIconLayout.large,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
};

// Public Bins Layer
export const publicBinLayer: LayerProps = {
  id: "public-bins",
  type: "symbol",
  layout: {
    "icon-image": "waste-basket",
    "icon-size": commonIconLayout.small,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
};

// Coach Parking Layer
export const coachParkingLayer: LayerProps = {
  id: "coach-parking",
  type: "symbol",
  layout: {
    "icon-image": "bus",
    "icon-size": commonIconLayout.large,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
};

export const publicWifiLayer: LayerProps = {
  id: "public-wifi",
  type: "symbol",
  layout: {
    "icon-image": "wifi",
    "icon-size": commonIconLayout.small,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
  filter: ["!=", ["get", "id"], ""],
};

export const libraryLayer: LayerProps = {
  id: "libraries",
  type: "symbol",
  layout: {
    "icon-image": "library",
    "icon-size": commonIconLayout.large,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
  filter: ["!=", ["get", "id"], ""],
};

export const carParkLayer: LayerProps = {
  id: "car-parks",
  type: "symbol",
  layout: {
    "icon-image": "car",
    "icon-size": commonIconLayout.large,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
  paint: {
    "icon-color": "#E74C3C", // Bright red
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const waterFountainLayer: LayerProps = {
  id: "water-fountains",
  type: "symbol",
  layout: {
    "icon-image": "drinking-water",
    "icon-size": commonIconLayout.small,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
  paint: {
    "icon-color": "#3498DB", // Blue
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const publicToiletLayer: LayerProps = {
  id: "public-toilets",
  type: "symbol",
  layout: {
    "icon-image": "toilet",
    "icon-size": commonIconLayout.medium,
    "icon-allow-overlap": false,
    visibility: "visible",
  },
  paint: {
    "icon-color": "#2ECC71", // Green
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};
