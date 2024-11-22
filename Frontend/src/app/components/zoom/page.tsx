import React from 'react';
import { useMap } from 'react-map-gl';

const ZoomControls = () => {
  const { current: map } = useMap();

  const zoomIn = () => {
    if (map) {
      map.zoomIn();
    }
  };

  const zoomOut = () => {
    if (map) {
      map.zoomOut();
    }
  };

  return (
    <div className="absolute top-4 right-4 z-10 flex flex-col space-y-2">
      <button
        onClick={zoomIn}
        className="bg-white p-2 rounded shadow hover:bg-gray-100"
      >
        +
      </button>
      <button
        onClick={zoomOut}
        className="bg-white p-2 rounded shadow hover:bg-gray-100"
      >
        -
      </button>
    </div>
  );
};

export default ZoomControls;