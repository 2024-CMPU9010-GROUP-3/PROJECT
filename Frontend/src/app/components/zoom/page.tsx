"use client";

import React from 'react';
import mapboxgl from 'mapbox-gl';
import type { Map } from 'mapbox-gl';

interface CustomNavigationControlProps {
  mapInstance: Map | null;
}

const CustomNavigationControl: React.FC<CustomNavigationControlProps> = ({ mapInstance }) => {
  const handleZoomIn = () => {
    if (mapInstance) {
      const currentZoom = mapInstance.getZoom();
      mapInstance.easeTo({
        zoom: Math.min(currentZoom + 1, mapInstance.getMaxZoom()),
        duration: 300
      });
    }
  };

  const handleZoomOut = () => {
    if (mapInstance) {
      const currentZoom = mapInstance.getZoom();
      mapInstance.easeTo({
        zoom: Math.max(currentZoom - 1, mapInstance.getMinZoom()),
        duration: 300
      });
    }
  };

  return (
    <div className="absolute top-4 right-4 z-[9999] flex flex-col space-y-2">
      <button
        onClick={handleZoomIn}
        className="bg-white p-3 rounded-full shadow hover:bg-gray-100 focus:outline-none"
        aria-label="zoomin"
      >
        +
      </button>
      <button
        onClick={handleZoomOut}
        className="bg-white p-3 rounded-full shadow hover:bg-gray-100 focus:outline-none"
        aria-label="zoomout"
      >
        -
      </button>
    </div>
  );
};

export default CustomNavigationControl;