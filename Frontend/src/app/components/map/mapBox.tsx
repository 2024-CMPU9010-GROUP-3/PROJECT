import React, { useRef, useEffect } from 'react';
import mapboxgl from 'mapbox-gl';

// Set Mapbox access token here
mapboxgl.accessToken = 'pk.eyJ1IjoieXVhbnNodW9kdSIsImEiOiJjbTFqenRldXIxNW04MmlzOHN5dDY4Zjl3In0.OFkHtyPHr1oSu9DrH6QpYQ';

const MapBox: React.FC = () => {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const lng = -6.2603; // Dublin longitude
  const lat = 53.3498; // Dublin latitude
  const zoom = 12; // Initial zoom level

  useEffect(() => {
    if (map.current) return; // Initialize map only once
    if (mapContainer.current) {
      map.current = new mapboxgl.Map({
        container: mapContainer.current,
        style: 'mapbox://styles/mapbox/satellite-streets-v12', // Satellite with streets
        center: [lng, lat],
        zoom: zoom,
      });
    }
  }, []);

  return <div ref={mapContainer} style={{ width: '100%', height: '500px' }} />;
};

export default MapBox;