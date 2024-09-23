import { GoogleMap, LoadScript, Marker } from "@react-google-maps/api";
import { useState, useEffect } from "react";

type Location = {
  lat: number;
  lng: number;
  name: string;
};

const containerStyle = {
  width: "100%",
  height: "400px",
};

const defaultCenter: Location = {
  lat: 53.3498,
  lng: -6.2603,
  name: "Dublin"
};

function MapComponent() {
  const [locations, setLocations] = useState<Location[]>([defaultCenter]);

  useEffect(() => {
    async function fetchLocations() {
      const response = await fetch('/api/locations');
      const data: Location[] = await response.json();
      setLocations(data.length > 0 ? data : [defaultCenter]);
    }
    fetchLocations();
  }, []);

  return (
    <LoadScript googleMapsApiKey={process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY as string}>
      <GoogleMap
        mapContainerStyle={containerStyle}
        center={locations[0]}
        zoom={10}
      >
        {locations.map((location, index) => (
          <Marker key={index} position={location} title={location.name} />
        ))}
      </GoogleMap>
    </LoadScript>
  );
}

export default MapComponent;