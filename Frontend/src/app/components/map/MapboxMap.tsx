"use client";

import React, { Fragment, Suspense, useEffect, useMemo, useState } from "react";
import { FaLocationDot } from "react-icons/fa6";
import DeckGL from "@deck.gl/react";
import "mapbox-gl/dist/mapbox-gl.css";
import {
  MapClickEvent,
  Coordinates,
  Point,
  CoordinatesForGeoJson,
} from "@/lib/interfaces/types";
import Map, { Layer, LayerProps, Marker, Source } from "react-map-gl";
import { lightingEffect, INITIAL_VIEW_STATE } from "@/lib/mapconfig";
import { GeoJSON, FeatureCollection } from "geojson";
// import { FloatingDock } from "@/components/ui/floating-dock";
// import { IconHome } from "@tabler/icons-react";
import { Slider } from "@/components/ui/slider";
import { cn } from "@/lib/utils";
import { Grid } from "react-loader-spinner";
import { Badge } from "@/components/ui/badge";

type SliderProps = React.ComponentProps<typeof Slider>;

const initialGeoJson: GeoJSON = {
  type: "FeatureCollection",
  features: [],
};

const LocationAggregatorMap = ({ className, ...props }: SliderProps) => {
  const [mapBoxApiKey, setMapBoxApiKey] = useState<string>("");
  const [isMarkerVisible, setIsMarkerVisible] = useState<boolean>(false);
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });
  const [pointsGeoJson, setPointsGeoJson] =
    useState<FeatureCollection>(initialGeoJson);
  // const [markerSquareSize, setMarkerSquareSize] = useState<number>(0.027);
  const [currentPositionCords, setCurrentPositionCords] = useState<Coordinates>(
    { latitude: 0, longitude: 0 }
  );
  const [sliderValue, setSliderValue] = useState<number>(40);
  const [sliderValueDisplay, setSliderValueDisplay] = useState<number>(40);
  const markerCoords = useMemo(
    () => [coordinates?.longitude, coordinates?.latitude],
    [coordinates]
  );

  const [circleCoordinates, setCircleCoordinates] = useState<number[][]>(() => {
    const radiusInMeters = sliderValue * 100; // Convert slider value to meters
    const radiusInDegrees = radiusInMeters / 111320; // Convert meters to degrees (approximation)
    const numPoints = 64; // Number of points to define the circle
    const angleStep = (2 * Math.PI) / numPoints;
    const coordinates = [];

    for (let i = 0; i < numPoints; i++) {
      const angle = i * angleStep;
      const x =
        markerCoords[0] +
        (radiusInDegrees * Math.cos(angle)) /
          Math.cos(markerCoords[1] * (Math.PI / 180));
      const y = markerCoords[1] + radiusInDegrees * Math.sin(angle);
      coordinates.push([x, y]);
    }

    // Close the circle
    coordinates.push(coordinates[0]);

    return coordinates;
  });

  // Checkbox code begins

  // Checkbox code ends

  useEffect(() => {
    console.log("Slider Value Commit>>>>>", sliderValue);
    console.log("Slider Value Display>>>>", sliderValueDisplay);
  }, [sliderValue, sliderValueDisplay]);

  // Handle map click event
  const handleMapClick = (event: unknown) => {
    const mapClickEvent = event as MapClickEvent; // Type assertion
    if (mapClickEvent.coordinate) {
      const [longitude, latitude] = mapClickEvent.coordinate;
      setCoordinates({ latitude, longitude });
      setIsMarkerVisible(true);
    }
  };

  function convertToGeoJson(points: Point[]): GeoJSON {
    console.log("Points>>>", points);

    return {
      type: "FeatureCollection",
      features: points?.map((point) => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [
            (point?.Longlat as CoordinatesForGeoJson)?.coordinates[0],
            (point?.Longlat as CoordinatesForGeoJson)?.coordinates[1],
          ],
        },
        properties: {
          Id: point.Id,
          Type: point.Type,
        },
      })),
    };
  }

  const fetchPointsFromDB = async (
    longitude: number,
    latitude: number,
    sliderValue: number
  ) => {
    console.log("fetchPointsFromDB>>>", longitude, latitude, sliderValue);
    const response = await fetch(
      `/api/points?long=${longitude}&lat=${latitude}&radius=${
        sliderValue * 100
      }`,
      {
        method: "GET",
        credentials: "include",
      }
    );
    const data = await response.json();
    const geoJson = convertToGeoJson(data?.response?.content);

    setPointsGeoJson(geoJson as FeatureCollection);
    console.log("GeoJson>>>", geoJson);
  };

  useEffect(() => {
    const radiusInMeters = sliderValueDisplay * 100; // Convert slider value to meters
    const radiusInDegrees = radiusInMeters / 111320; // Convert meters to degrees (approximation)
    const numPoints = 64; // Number of points to define the circle
    const angleStep = (2 * Math.PI) / numPoints;
    const coordinates = [];

    for (let i = 0; i < numPoints; i++) {
      const angle = i * angleStep;
      const x =
        markerCoords[0] +
        (radiusInDegrees * Math.cos(angle)) /
          Math.cos(markerCoords[1] * (Math.PI / 180));
      const y = markerCoords[1] + radiusInDegrees * Math.sin(angle);
      coordinates.push([x, y]);
    }

    // Close the circle
    coordinates.push(coordinates[0]);

    setCircleCoordinates(coordinates);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sliderValueDisplay, markerCoords]);

  useEffect(() => {
    // Fetch Mapbox API key from the server
    const fetchMapboxKey = async () => {
      const response = await fetch("/api/mapbox-token");
      const data = await response.json();
      setMapBoxApiKey(data.mapBoxApiKey);
    };
    fetchMapboxKey();
  }, []);

  // Queries
  useEffect(() => {
    fetchPointsFromDB(
      coordinates?.longitude,
      coordinates?.latitude,
      sliderValue
    );
    console.log(coordinates, sliderValue);

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [coordinates, sliderValue]);

  useEffect(() => {
    console.log(currentPositionCords);
  }, [currentPositionCords]);

  const options = {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 0,
  };

  function success(pos: GeolocationPosition) {
    const crd = pos.coords;
    setCurrentPositionCords({
      latitude: crd.latitude,
      longitude: crd.longitude,
    });
  }

  function error(err: unknown) {
    const geolocationError = err as GeolocationPositionError;

    console.warn(
      `ERROR(${geolocationError.code}): ${geolocationError.message}`
    );
  }

  // Get current position
  useEffect(() => {
    navigator.geolocation.getCurrentPosition(success, error, options);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const layerStyle: LayerProps = {
    id: "circle",
    type: "line",
    paint: {
      "line-color": "#007cbf", // Color of the circle outline
      "line-width": 2, // Width of the circle outline
      "line-opacity": 0.8, // Opacity of the circle outline
    },
  };

  return (
    <div className="flex flex-col lg:flex-row min-h-screen">
      {/* Map Container - Taller on mobile */}
      <div className="w-full lg:w-[75%] h-[60vh] sm:h-[70vh] lg:h-screen relative">
        {mapBoxApiKey ? (
          <DeckGL
            effects={[lightingEffect]}
            initialViewState={INITIAL_VIEW_STATE}
            controller={true}
            onClick={handleMapClick}
            style={{ width: "100%", height: "100%" }}
          >
            <Map
              mapboxAccessToken={mapBoxApiKey}
              mapStyle="mapbox://styles/mapbox/streets-v12"
              antialias={true}
              style={{ width: "100%", height: "100%" }}
            >
              {/* Map content remains the same */}
              <Marker
                longitude={coordinates?.longitude}
                latitude={coordinates?.latitude}
                anchor="center"
              >
                <div>
                  <FaLocationDot size={50} color="FFA15A" />
                </div>
              </Marker>
              <Marker
                latitude={currentPositionCords?.latitude}
                longitude={currentPositionCords?.longitude}
              >
                <div>
                  <FaLocationDot size={35} color="blue" />
                </div>
              </Marker>
              <Source
                id="circle"
                type="geojson"
                data={{
                  type: "Feature",
                  geometry: {
                    type: "Polygon",
                    coordinates: [circleCoordinates],
                  },
                }}
              >
                <Layer {...layerStyle} />
              </Source>
              <Suspense fallback={<div>Loading...</div>}>
                <Source
                  id="points"
                  type="geojson"
                  data={pointsGeoJson}
                  cluster={true}
                  clusterMaxZoom={14} // Max zoom to cluster points on
                  clusterRadius={50}
                >
                  <Layer
                    id="clusters"
                    type="symbol"
                    layout={{
                      "icon-image": "parking-garage",
                      "icon-size": 1.5,
                      "icon-allow-overlap": true,
                    }}
                  />
                  <Layer
                    id="cluster-count"
                    type="symbol"
                    layout={{
                      "text-field": "{point_count_abbreviated}",
                      "text-font": [
                        "DIN Offc Pro Medium",
                        "Arial Unicode MS Bold",
                      ],
                      "text-size": 12,
                    }}
                  />
                  <Layer
                    id="unclustered-point"
                    type="symbol"
                    layout={{
                      "icon-image": "parking-garage",
                      "icon-size": 1.5,
                      "icon-allow-overlap": true,
                    }}
                  />
                </Source>
              </Suspense>
            </Map>
          </DeckGL>
        ) : (
          <div className="absolute inset-0 flex items-center justify-center">
            <Grid
              visible={true}
              height="80"
              width="80"
              color="#ffa15a"
              ariaLabel="grid-loading"
              radius="12.5"
            />
          </div>
        )}
      </div>

      {/* Sidebar - Full width on mobile, scrollable */}
      <div className="w-full lg:w-[25%] h-[40vh] sm:h-[30vh] lg:h-screen p-3 sm:p-4 lg:p-6 bg-gray-50 overflow-y-auto">
        <div className="space-y-3 sm:space-y-4 lg:space-y-6 max-w-lg mx-auto lg:max-w-none">
          {mapBoxApiKey ? (
            <>
              <div className="px-2 sm:px-3 lg:px-4">
                <h1 className="text-lg sm:text-xl lg:text-2xl font-bold text-gray-900 tracking-tight">
                  Magpie Dashboard
                </h1>
              </div>

              {/* Search Radius Card */}
              <div className="px-2 sm:px-3 lg:px-4">
                <div className="w-full bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4">
                  <div className="space-y-2 sm:space-y-3">
                    <div>
                      <label className="text-sm lg:text-base font-medium text-gray-700 mb-2 block">
                        Search Radius
                      </label>
                      <Slider
                        onValueChange={(value) =>
                          setSliderValueDisplay(value[0])
                        }
                        onValueCommit={(value) => setSliderValue(value[0])}
                        defaultValue={[sliderValue]}
                        max={100}
                        step={1}
                        className={cn("w-full touch-none", className)}
                        {...props}
                      />
                    </div>
                    <div className="text-sm lg:text-base font-medium text-gray-600">
                      {sliderValueDisplay * 100} meters
                    </div>
                  </div>
                </div>
              </div>
            </>
          ) : null}

          {/* Marker Data Card */}
          <div className="px-2 sm:px-3 lg:px-4">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4">
              <h2 className="text-base sm:text-lg lg:text-xl font-semibold text-gray-900 mb-2 sm:mb-3">
                Marker Data
              </h2>
              {isMarkerVisible ? (
                <Suspense
                  fallback={<div className="animate-pulse">Loading...</div>}
                >
                  <div className="flex items-center justify-between">
                    <span className="text-sm sm:text-base lg:text-lg text-gray-700">
                      Parking
                    </span>
                    <Badge
                      variant="secondary"
                      className="px-2 py-1 text-xs sm:text-sm rounded-full bg-gray-100"
                    >
                      {pointsGeoJson?.features?.length || 0} Spots
                    </Badge>
                  </div>
                </Suspense>
              ) : (
                <div className="text-xs sm:text-sm lg:text-base text-gray-500">
                  Place a marker on the map to view data
                </div>
              )}
            </div>
          </div>

          {/* Filter Options Card */}
          {mapBoxApiKey && (
            <div className="px-2 sm:px-3 lg:px-4">
              <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4">
                <h2 className="text-base sm:text-lg lg:text-xl font-semibold text-gray-900 mb-2 sm:mb-3">
                  Filter Options
                </h2>
                <div className="space-y-2">
                  {[
                    "Coach Parking",
                    "Bike Stand",
                    "Public Toilet",
                    "Parking Meter",
                    "Parking",
                  ].map((option, index) => (
                    <label
                      key={index}
                      className="flex items-center space-x-3 group cursor-pointer hover:bg-gray-50 p-2 rounded-md transition-colors min-h-[44px] touch-manipulation"
                    >
                      <input
                        type="checkbox"
                        id={`option${index + 1}`}
                        name={`option${index + 1}`}
                        className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 transition-colors"
                      />
                      <span className="text-sm sm:text-base text-gray-700 group-hover:text-gray-900">
                        {option}
                      </span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default LocationAggregatorMap;
