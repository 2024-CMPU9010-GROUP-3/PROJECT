import { LayerProps } from "react-map-gl";

// Update commonIconLayout to have single standard size
const commonIconLayout = {
  small: ["interpolate", ["linear"], ["zoom"], 10, 0.5, 14, 0.8] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
  standard: ["interpolate", ["linear"], ["zoom"], 12, 0.8, 16, 1.2] as [
    string,
    [string],
    [string],
    number,
    number,
    number,
    number
  ],
  large: ["interpolate", ["linear"], ["zoom"], 14, 1.2, 18, 1.6] as [
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

export const parkingMeterLayers: Record<string, LayerProps> = {
  far: {
    id: "parking-meters-far",
    type: "symbol",
    layout: {
      "icon-image": "parking",
      "icon-size": commonIconLayout.small,
      "icon-allow-overlap": true,
      "icon-ignore-placement": true,
      visibility: "visible",
    },
    paint: {
      "icon-opacity": 0.7,
    },
    filter: ["all", ["!=", ["get", "id"], ""], ["<=", ["zoom"], 14]],
  },
  medium: {
    id: "parking-meters-medium",
    type: "symbol",
    layout: {
      "icon-image": "parking",
      "icon-size": commonIconLayout.standard,
      "icon-allow-overlap": true,
      "icon-ignore-placement": false,
      visibility: "visible",
    },
    paint: {
      "icon-opacity": 0.8,
    },
    filter: [
      "all",
      ["!=", ["get", "id"], ""],
      [">", ["zoom"], 14],
      ["<=", ["zoom"], 16],
    ],
  },
  close: {
    id: "parking-meters-close",
    type: "symbol",
    layout: {
      "icon-image": "parking",
      "icon-size": commonIconLayout.large,
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
    filter: ["all", ["!=", ["get", "id"], ""], [">", ["zoom"], 16]],
  },
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
        "icon-size": commonIconLayout.small,
        "icon-allow-overlap": true,
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
        "icon-size": commonIconLayout.standard,
        "icon-allow-overlap": true,
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
        "icon-size": commonIconLayout.large,
        "icon-allow-overlap": true,
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

export const bikeStandLayers = createZoomBasedLayers(
  "bike-stands",
  "bicycle",
  colors.accent,
  true
);
export const bikeSharingLayers = createZoomBasedLayers(
  "bike-sharing",
  "bicycle-share",
  colors.secondary,
  true
);
export const accessibleParkingLayers = createZoomBasedLayers(
  "accessible-parking",
  "disabled",
  colors.primary
);
export const publicBinLayers = createZoomBasedLayers(
  "public-bins",
  "waste-basket",
  colors.neutral
);
export const coachParkingLayers = createZoomBasedLayers(
  "coach-parking",
  "bus",
  colors.primary
);
export const publicWifiLayers = createZoomBasedLayers(
  "public-wifi",
  "wifi",
  colors.info
);
export const libraryLayers = createZoomBasedLayers(
  "libraries",
  "library",
  colors.secondary
);
export const carParkLayers = createZoomBasedLayers(
  "car-parks",
  "car",
  colors.primary
);
export const waterFountainLayers = createZoomBasedLayers(
  "water-fountains",
  "drinking-water",
  colors.water
);
export const publicToiletLayers = createZoomBasedLayers(
  "public-toilets",
  "toilet",
  colors.info
);
