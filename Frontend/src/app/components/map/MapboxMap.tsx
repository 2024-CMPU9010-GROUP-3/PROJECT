"use client";

import React, { Fragment, useEffect, useMemo, useState } from "react";
import { FaLocationDot } from "react-icons/fa6";
import DeckGL from "@deck.gl/react";
import "mapbox-gl/dist/mapbox-gl.css";
import { MapClickEvent } from "@/lib/interfaces/types";
import { Coordinates } from "@/lib/interfaces/types";
import Map, { Layer, LayerProps, Marker, Source } from "react-map-gl";
import { lightingEffect, INITIAL_VIEW_STATE } from "@/lib/mapconfig";
// import { FloatingDock } from "@/components/ui/floating-dock";
// import { IconHome } from "@tabler/icons-react";
import { Slider } from "@/components/ui/slider";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";
import { Grid } from "react-loader-spinner";
import { Badge } from "@/components/ui/badge";

type SliderProps = React.ComponentProps<typeof Slider>;
type TListOfPlaces = {
  name: string;
  spaces: number;
  address: string;
};

const LocationAggregatorMap = ({ className, ...props }: SliderProps) => {
  const [mapBoxApiKey, setMapBoxApiKey] = useState<string>("");
  const [isMarkerVisible, setIsMarkerVisible] = useState<boolean>(false);
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });
  // const [markerSquareSize, setMarkerSquareSize] = useState<number>(0.027);
  const [currentPositionCords, setCurrentPositionCords] = useState<Coordinates>(
    { latitude: 0, longitude: 0 }
  );
  const [sliderValue, setSliderValue] = useState<number>(40);
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

  useEffect(() => {
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

    setCircleCoordinates(coordinates);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sliderValue, markerCoords]);

  useEffect(() => {
    // Fetch Mapbox API key from the server
    const fetchMapboxKey = async () => {
      const response = await fetch("/api/mapbox-token");
      const data = await response.json();
      setMapBoxApiKey(data.mapBoxApiKey);
    };
    fetchMapboxKey();
  }, []);

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

  const listOfPlaces: TListOfPlaces[] = [
    {
      name: "Arnots Car Park",
      spaces: 10,
      address: "Best Car Parks Ltd Arnotts 12 Henry Street Dublin D01 C3Y9",
    },
    {
      name: "Jervis Street Car Park",
      spaces: 90,
      address:
        "Jervis Shopping Centre, Jervis St, North City, Dublin 1, D01 X868",
    },
    {
      name: "Ilac Shopping Centre Car Park",
      spaces: 50,
      address: "Ilac Centre, Parnell St, North City, Dublin 1, D01 W861",
    },
    {
      name: "Park Rite Drury Street Car Park",
      spaces: 72,
      address: "Drury St, Dublin 2, D02 V586",
    },
    {
      name: "Q-Park The Spire",
      spaces: 35,
      address: "Irish Life Mall, Abbey St Lwr, North City, Dublin 1, D01 E9X0",
    },
    {
      name: "St Stephen's Green Car Park",
      spaces: 48,
      address: "St Stephen's Green Shopping Centre, Dublin 2, D02 XY88",
    },
    {
      name: "Trinity Street Car Park",
      spaces: 35,
      address: "Trinity St, Dublin 2, D02 R274",
    },
    {
      name: "Park Rite IFSC Car Park",
      spaces: 56,
      address: "IFSC, Commons St, North Dock, Dublin 1, D01 DA06",
    },
    {
      name: "Thomas Street Car Park",
      spaces: 49,
      address: "Thomas St, Dublin 8, D08 K6Y9",
    },
    {
      name: "Park Rite Parnell Street Car Park",
      spaces: 34,
      address: "Parnell St, Rotunda, Dublin 1, D01 EK28",
    },
  ];

  // Handle map click event
  const handleMapClick = (event: unknown) => {
    const mapClickEvent = event as MapClickEvent; // Type assertion
    if (mapClickEvent.coordinate) {
      const [longitude, latitude] = mapClickEvent.coordinate;
      setCoordinates({ latitude, longitude });
      setIsMarkerVisible(true);
    }
  };

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
    <div>
      {mapBoxApiKey ? (
        <DeckGL
          effects={[lightingEffect]}
          initialViewState={INITIAL_VIEW_STATE}
          controller={true}
          onClick={handleMapClick}
        >
          <Map
            mapboxAccessToken={mapBoxApiKey}
            mapStyle="mapbox://styles/mapbox/streets-v8"
            antialias={true}
          >
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
          </Map>
        </DeckGL>
      ) : (
        <div className="absolute top-1/2 right-1/2">
          <Grid
            visible={true}
            height="80"
            width="80"
            color="#ffa15a"
            ariaLabel="grid-loading"
            radius="12.5"
            wrapperStyle={{}}
            wrapperClass="grid-wrapper"
          />
        </div>
      )}
      {isMarkerVisible && (
        <div className="absolute right-24 top-1/3 bg-white p-5 rounded-xl max-h-[400px] max-w-[450px] overflow-scroll">
          <div>
            <div className="text-2xl font-bold mb-5">Dublin Car Park</div>
            <div className="space-y-4">
              {listOfPlaces?.map((place, index) => (
                <Fragment key={index}>
                  <div className="p-2 border-2 border-gray-200 rounded-xl">
                    <div className="flex align-middle justify-between my-2">
                      <div className="text-xl font-medium">{place?.name}</div>
                      <Badge
                        variant="secondary"
                        className="text-sm rounded-full"
                      >
                        {place?.spaces} spaces
                      </Badge>
                    </div>
                    <div>Address: {place?.address}</div>
                  </div>
                </Fragment>
              ))}
            </div>
          </div>
        </div>
      )}
      {mapBoxApiKey ? (
        <>
          <div className="absolute bg-transparent top-20 right-32 scale-125 transition-all">
            <div className="2xl:min-w-[200px] 2xl:max-w-[300px] bg-white rounded-xl ">
              <div className="px-2 py-4 space-y-2">
                <div className="space-y-1">
                  <label>Distance</label>
                  <Slider
                    onValueChange={(value) => setSliderValue(value[0])}
                    defaultValue={[sliderValue]}
                    max={100}
                    step={1}
                    className={cn("w-[60%]", className)}
                    {...props}
                  />
                </div>
                <div className="">
                  <span className="p-1">{sliderValue * 100} meters</span>
                </div>
                <div className="flex justify-center gap-5">
                  <div className="w-2/4">
                    <label>Long</label>
                    <Input />
                  </div>
                  <div className="w-2/4">
                    <label>Lat</label>
                    <Input />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </>
      ) : (
        <></>
      )}
    </div>
  );
};

export default LocationAggregatorMap;
