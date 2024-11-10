import { LayerProps } from "react-map-gl";

// Update commonIconLayout to have single standard size
const commonIconLayout = {
  standard: ["interpolate", ["linear"], ["zoom"], 12, 0.8, 16, 1.2] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
};

// Modern color palette
const colors = {
  primary: "#2563EB", // Royal Blue
  secondary: "#059669", // Emerald
  accent: "#D97706", // Amber
  danger: "#DC2626", // Red
  info: "#6366F1", // Indigo
  neutral: "#4B5563", // Gray
  water: "#0EA5E9", // Sky Blue
};

export const parkingClusterStyles = {
  symbol: {
    id: "clusters",
    type: "symbol" as const,
    layout: {
      "icon-image": "parking-garage",
      "icon-size": commonIconLayout.standard,
      "icon-allow-overlap": true,
    },
  },
  count: {
    id: "cluster-count",
    type: "symbol" as const,
    layout: {
      "text-field": "{point_count_abbreviated}",
      "text-font": ["DIN Offc Pro Medium", "Arial Unicode MS Bold"],
      "text-size": 12,
    },
  },
  unclustered: {
    id: "unclustered-point",
    type: "symbol" as const,
    layout: {
      "icon-image": "parking-garage",
      "icon-size": commonIconLayout.standard,
      "icon-allow-overlap": true,
    },
  },
};

export const parkingMeterLayer: LayerProps = {
  id: "parking-meters",
  type: "symbol",
  layout: {
    "icon-image": "parking",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    "icon-ignore-placement": false,
    visibility: "visible",
    "text-field": ["get", "name"],
    "text-size": 11,
    "text-offset": [0, 1.2],
    "text-optional": true,
    "text-font": ["DIN Pro Medium", "Arial Unicode MS Bold"],
  },
  paint: {
    "text-color": colors.primary,
    "text-halo-color": "#ffffff",
    "text-halo-width": 2,
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const bikeStandLayer: LayerProps = {
  id: "bike-stands",
  type: "symbol",
  layout: {
    "icon-image": "bicycle",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    "icon-ignore-placement": false,
    visibility: "visible",
    "text-field": ["get", "name"],
    "text-size": 11,
    "text-offset": [0, 1.2],
    "text-optional": true,
    "text-font": ["DIN Pro Medium", "Arial Unicode MS Bold"],
  },
  paint: {
    "text-color": colors.accent,
    "text-halo-color": "#ffffff",
    "text-halo-width": 2,
    "icon-opacity": 0.9,
  },
};

export const bikeSharingLayer: LayerProps = {
  id: "bike-sharing",
  type: "symbol",
  layout: {
    "icon-image": "bicycle-share",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    "icon-ignore-placement": false,
    visibility: "visible",
    "text-field": ["get", "name"],
    "text-size": 11,
    "text-offset": [0, 1.2],
    "text-optional": true,
    "text-font": ["DIN Pro Medium", "Arial Unicode MS Bold"],
  },
  paint: {
    "text-color": colors.secondary,
    "text-halo-color": "#ffffff",
    "text-halo-width": 2,
    "icon-opacity": 0.9,
  },
};

export const accessibleParkingLayer: LayerProps = {
  id: "accessible-parking",
  type: "symbol",
  layout: {
    "icon-image": "disabled",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
    "text-font": ["DIN Pro Medium", "Arial Unicode MS Bold"],
  },
  paint: {
    "icon-opacity": 0.9,
    "icon-color": colors.primary,
  },
};

export const publicBinLayer: LayerProps = {
  id: "public-bins",
  type: "symbol",
  layout: {
    "icon-image": "waste-basket",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-opacity": 0.8,
    "icon-color": colors.neutral,
  },
};

export const coachParkingLayer: LayerProps = {
  id: "coach-parking",
  type: "symbol",
  layout: {
    "icon-image": "bus",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-opacity": 0.9,
    "icon-color": colors.primary,
  },
};

export const publicWifiLayer: LayerProps = {
  id: "public-wifi",
  type: "symbol",
  layout: {
    "icon-image": "wifi",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-opacity": 0.8,
    "icon-color": colors.info,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const libraryLayer: LayerProps = {
  id: "libraries",
  type: "symbol",
  layout: {
    "icon-image": "library",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-opacity": 0.9,
    "icon-color": colors.secondary,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const carParkLayer: LayerProps = {
  id: "car-parks",
  type: "symbol",
  layout: {
    "icon-image": "car",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-color": colors.primary,
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const waterFountainLayer: LayerProps = {
  id: "water-fountains",
  type: "symbol",
  layout: {
    "icon-image": "drinking-water",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-color": colors.water,
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};

export const publicToiletLayer: LayerProps = {
  id: "public-toilets",
  type: "symbol",
  layout: {
    "icon-image": "toilet",
    "icon-size": commonIconLayout.standard,
    "icon-allow-overlap": true,
    visibility: "visible",
  },
  paint: {
    "icon-color": colors.info,
    "icon-opacity": 0.9,
  },
  filter: ["!=", ["get", "id"], ""],
};
