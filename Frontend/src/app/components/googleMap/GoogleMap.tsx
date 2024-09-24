import React, { useEffect, useRef } from 'react';

const GoogleMap = () => {
  const mapRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (mapRef.current) {
      const map = new google.maps.Map(mapRef.current, {
        center: { lat: -34.397, lng: 150.644 },
        zoom: 8,
      });
    }
  }, []);

  return (
    <div ref={mapRef} style={{ width: '100%', height: '500px' }}></div>
  );
};

export default GoogleMap;
