import { AmbientLight, PointLight, LightingEffect } from "@deck.gl/core";

export const ambientLight: AmbientLight = new AmbientLight({
  color: [255, 255, 255],
  intensity: 1.0,
});

export const pointLight1: PointLight = new PointLight({
  color: [255, 255, 255],
  intensity: 0.8,
  position: [-0.144528, 49.739968, 80000],
});

export const pointLight2: PointLight = new PointLight({
  color: [255, 255, 255],
  intensity: 0.8,
  position: [-3.807751, 54.104682, 8000],
});

export const lightingEffect: LightingEffect = new LightingEffect({
  ambientLight,
  pointLight1,
  pointLight2,
});

interface Material {
  ambient: number;
  diffuse: number;
  shininess: number;
  specularColor: [number, number, number];
}

export const material: Material = {
  ambient: 0.64,
  diffuse: 0.6,
  shininess: 32,
  specularColor: [51, 51, 51],
};

interface ViewState {
  longitude: number;
  latitude: number;
  zoom: number;
  minZoom: number;
  maxZoom: number;
  pitch: number;
  bearing: number;
}

export const INITIAL_VIEW_STATE: ViewState = {
  latitude: 53.3498,
  longitude: -6.2674862,
  zoom: 11,
  minZoom: 2,
  maxZoom: 100,
  pitch: 0,
  bearing: 0,
};

export const colorRange: [number, number, number][] = [
  [1, 152, 189],
  [73, 227, 206],
  [216, 254, 181],
  [254, 237, 177],
  [254, 173, 84],
  [209, 55, 78],
];
