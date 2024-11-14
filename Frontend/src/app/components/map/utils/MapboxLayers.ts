import { LayerProps } from "react-map-gl";

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

const createZoomBasedLayers = (
  baseId: string,
  iconImage: string,
  color: string,
  showLabel: boolean = false
): Record<string, LayerProps> => {
  return {
    far: {
      id: `${baseId}-far`,
      type: "symbol",
      layout: {
        "icon-image": iconImage,
        "icon-size": 0.3,
        "icon-allow-overlap": false,
        visibility: "visible",
      },
      paint: {
        "icon-opacity": 0.7,
        "icon-color": color,
      },
      filter: ["all", ["!=", ["get", "id"], ""], ["<=", ["zoom"], 14]],
    },
    medium: {
      id: `${baseId}-medium`,
      type: "symbol",
      layout: {
        "icon-image": iconImage,
        "icon-size": 0.3,
        "icon-allow-overlap": false,
        visibility: "visible",
      },
      paint: {
        "icon-opacity": 0.8,
        "icon-color": color,
      },
      filter: [
        "all",
        ["!=", ["get", "id"], ""],
        [">", ["zoom"], 14],
        ["<=", ["zoom"], 16],
      ],
    },
    close: {
      id: `${baseId}-close`,
      type: "symbol",
      layout: {
        "icon-image": iconImage,
        "icon-size": 0.3,
        "icon-allow-overlap": false,
        visibility: "visible",
        ...(showLabel && {
          "text-field": ["get", "name"],
          "text-size": 11,
          "text-offset": [0, 1.2],
          "text-optional": true,
          "text-font": ["DIN Pro Medium", "Arial Unicode MS Bold"],
        }),
      },
      paint: {
        "icon-opacity": 0.9,
        "icon-color": color,
        ...(showLabel && {
          "text-color": color,
          "text-halo-color": "#ffffff",
          "text-halo-width": 2,
        }),
      },
      filter: ["all", ["!=", ["get", "id"], ""], [">", ["zoom"], 16]],
    },
  };
};

export const parkingClusterStyles = createZoomBasedLayers(
  "parking-clusters",
  "custom_parking",
  colors.primary
);

export const parkingMeterLayers = createZoomBasedLayers(
  "parking-meters",
  "custom_parking_meter",
  colors.primary
);
export const bikeStandLayers = createZoomBasedLayers(
  "bike-stands",
  "custom_bicycle",
  colors.accent,
  true
);
export const bikeSharingLayers = createZoomBasedLayers(
  "bike-sharing",
  "custom_bicycle_share",
  colors.secondary,
  true
);
export const accessibleParkingLayers = createZoomBasedLayers(
  "accessible-parking",
  "custom_accessible_parking",
  colors.primary
);
export const publicBinLayers = createZoomBasedLayers(
  "public-bins",
  "custom_public_bins",
  colors.neutral
);
export const coachParkingLayers = createZoomBasedLayers(
  "coach-parking",
  "custom_bus",
  colors.primary
);
export const publicWifiLayers = createZoomBasedLayers(
  "public-wifi",
  "custom_public_wifi",
  colors.info
);
export const libraryLayers = createZoomBasedLayers(
  "libraries",
  "custom_library",
  colors.secondary
);
export const carParkLayers = createZoomBasedLayers(
  "car-parks",
  "custom_car_parks",
  colors.primary
);
export const waterFountainLayers = createZoomBasedLayers(
  "water-fountains",
  "custom_water_fountain",
  colors.water
);
export const publicToiletLayers = createZoomBasedLayers(
  "public-toilets",
  "custom_toilet",
  colors.info
);
