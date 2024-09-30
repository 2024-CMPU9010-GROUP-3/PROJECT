// components/GoogleMap.tsx
import { GoogleMap, LoadScript, Marker } from "@react-google-maps/api";

type LatLng = {
  lat: number;
  lng: number;
};

const containerStyle = {
  width: "100%",
  height: "400px",
};

const center: LatLng = {
  lat: 53.3498, // Default latitude (Dublin example)
  lng: -6.2603, // Default longitude
};

function MapComponent() {
  const currentLocation: LatLng = center;

  return (
    <LoadScript
      googleMapsApiKey={process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY as string}
    >
      <GoogleMap
        mapContainerStyle={containerStyle}
        center={currentLocation}
        zoom={10}
      >
        {/* Optional Marker */}
        <Marker position={currentLocation} />
      </GoogleMap>
    </LoadScript>
  );
}

export default MapComponent;
