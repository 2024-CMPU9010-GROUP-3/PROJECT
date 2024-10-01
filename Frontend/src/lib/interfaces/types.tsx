export interface Coordinates {
    latitude: number;
    longitude: number;
}

// Interface for distance scales within the viewport
export interface DistanceScales {
    unitsPerMeter: number[];
    metersPerUnit: number[];
    unitsPerDegree: number[];
    degreesPerUnit: number[];
}

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