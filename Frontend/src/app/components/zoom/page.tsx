"use client";

import React from 'react';
import { useMap } from 'react-map-gl';

const ZoomControls = () => {
  const { current: map } = useMap();

  const zoomIn = () => {
    if (map) {
      const currentZoom = map.getZoom();
      map.easeTo({ zoom: currentZoom + 1, duration: 300 });
    }
  };

  const zoomOut = () => {
    if (map) {
      const currentZoom = map.getZoom();
      map.easeTo({ zoom: currentZoom - 1, duration: 300 });
    }
  };

  return (
    <div className="absolute top-4 right-4 z-10 flex flex-col space-y-2">
      <button
        onClick={zoomIn}
        className="bg-white p-3 rounded-full shadow hover:bg-gray-100 focus:outline-none"
        aria-label="Zoom in"
      >
        +
      </button>
      <button
        onClick={zoomOut}
        className="bg-white p-3 rounded-full shadow hover:bg-gray-100 focus:outline-none"
        aria-label="Zoom out"
      >
        -
      </button>
    </div>
  );
};

export default ZoomControls;