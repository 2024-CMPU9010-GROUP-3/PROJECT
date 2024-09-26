"use client";

import React, { useState } from "react";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import Map, { Marker } from "react-map-gl";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { HexagonLayer } from "@deck.gl/aggregation-layers";
import DeckGL from "@deck.gl/react";
import "mapbox-gl/dist/mapbox-gl.css";

import {
  lightingEffect,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  material,
  INITIAL_VIEW_STATE,
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  colorRange,
} from "@/lib/mapconfig";

interface Coordinates {
  latitude: number;
  longitude: number;
}

const LocationAggregatorMap = () => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });
  // const getCurrentLocation: Promise<Coordinates> = () => {
  //   return new Promise((resolve, reject) => {
  //     // Check if geolocation is supported
  //     if (!navigator.geolocation) {
  //       reject(new Error("Geolocation is not supported by this browser."));
  //       return;
  //     }

  //     // Use the Geolocation API to get the user's position
  //     navigator.geolocation.getCurrentPosition(
  //       (position: GeolocationPosition) => {
  //         const coordinates: Coordinates = {
  //           latitude: position.coords.latitude,
  //           longitude: position.coords.longitude,
  //         };
  //         resolve(coordinates);
  //       },
  //       (error: GeolocationPositionError) => {
  //         reject(error);
  //       }
  //     );
  //   });
  // };

  return (
    <div>
      <DeckGL
        effects={[lightingEffect]}
        initialViewState={INITIAL_VIEW_STATE}
        controller={true}
      >
        <Map
          mapboxAccessToken={process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN}
          mapStyle="mapbox://styles/mapbox/streets-v8"
          antialias={true}
        >
          {/* <Marker
            coordinates={[coordinates?.longitude, coordinates?.longitude]}
            anchor="bottom"
          >
            <img src="path_to_marker_image" alt="User Location" />
          </Marker> */}
        </Map>
      </DeckGL>
    </div>
  );
};

export default LocationAggregatorMap;
