"use client";

// React core
import React, { Suspense, useEffect, useMemo, useState } from 'react';

// Third-party packages
import DeckGL from '@deck.gl/react';
import { GeoJSON } from 'geojson';
import { FaLocationDot } from 'react-icons/fa6';
import { Grid } from 'react-loader-spinner';
import Map, { Layer, LayerProps, Marker, Source } from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { Eye, EyeOff } from 'lucide-react';
import Image from 'next/image';
import { useOnborda } from "onborda";
import { useSession } from '@/app/context/SessionContext';

// Local components
import { Slider } from '@/components/ui/slider';


// Local utils and configs
import { lightingEffect, INITIAL_VIEW_STATE } from '@/lib/mapconfig';
import packageJson from '../../../../package.json';
import MapSources from './utils/MapSources';
import { cn } from '@/lib/utils';

// Types and interfaces
import {
  MapClickEvent,
  Coordinates,
  Point,
  CoordinatesForGeoJson,
  ImageConfig,
} from '@/lib/interfaces/types';

type SliderProps = React.ComponentProps<typeof Slider>;
type GeoJsonCollection =
  | "parking_meter"
  | "bike_stand"
  | "public_wifi_access_point"
  | "library"
  | "multistorey_car_parking"
  | "drinking_water_fountain"
  | "public_toilet"
  | "bike_sharing_station"
  | "parking"
  | "accessible_parking"
  | "public_bins"
  | "coach_parking";

interface Option {
  label: string;
  value: string;
}

const MultiSelectOptions: Option[] = [
  { label: "Parking Meter", value: "parking_meter" },
  { label: "Bike Stand", value: "bike_stand" },
  { label: "Public Wifi", value: "public_wifi_access_point" },
  { label: "Library", value: "library" },
  { label: "Multi Storey Car Park", value: "multistorey_car_parking" },
  { label: "Drinking Water Fountain", value: "drinking_water_fountain" },
  { label: "Public Toilet", value: "public_toilet" },
  { label: "Bike Sharing Station", value: "bike_sharing_station" },
  { label: "Parking", value: "parking" },
  { label: "Accessible Parking", value: "accessible_parking" },
  { label: "Public Bins", value: "public_bins" },
  { label: "Coach Parking", value: "coach_parking" },
];

// Array of image paths to load 
const IMAGES: ImageConfig[] = [
  { id: 'custom_parking', path: '/images/parking.png' },
  { id: 'custom_parking_meter', path: '/images/parking_meter.png' },
  { id: 'custom_bicycle', path: '/images/bicycle.png' },
  { id: 'bicycle_share', path: '/images/bicycle_share.png' },
  { id: 'custom_bicycle_share', path: '/images/bicycle_share.png' },
  { id: 'custom_accessible_parking', path: '/images/accessibleParking.png' },
  { id: 'custom_public_bins', path: '/images/bin.png' },
  { id: 'custom_public_wifi', path: '/images/wifi.png' },
  { id: 'custom_bus', path: '/images/bus.png' },
  { id: 'custom_library', path: '/images/library.png' },
  { id: 'custom_car_parks', path: '/images/car_park.png' },
  { id: 'custom_water_fountain', path: '/images/water_fountain.png' },
  { id: 'custom_toilet', path: '/images/toilet.png' },
];

// Array of icon paths to load
const iconMap: Record<string, string> = {
  parking_meter: '/images/parking_meter.png',
  bike_stand: '/images/bicycle.png',
  public_wifi_access_point: '/images/wifi.png',
  library: '/images/library.png',
  multistorey_car_parking: '/images/car_park.png',
  drinking_water_fountain: '/images/water_fountain.png',
  public_toilet: '/images/toilet.png',
  bike_sharing_station: '/images/bicycle_share.png',
  parking: '/images/parking.png',
  accessible_parking: '/images/accessibleParking.png',
  public_bins: '/images/bin.png',
  coach_parking: '/images/bus.png',
};

const LocationAggregatorMap = ({ className, ...props }: SliderProps) => {
  const [mapBoxApiKey, setMapBoxApiKey] = useState<string>("");
  const [isMarkerVisible, setIsMarkerVisible] = useState<boolean>(false);
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });
  const [pointsGeoJson, setPointsGeoJson] = useState<Record<
    string,
    GeoJSON
  > | null>(null);
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

  // Images state
  const [imagesLoaded, setImagesLoaded] = useState({
    custom_parking: false,
    custom_parking_meter: false,
    custom_bicycle: false,
    custom_bicycle_share: false,
    custom_accessible_parking: false,
    custom_public_bins: false,
    custom_public_wifi: false,
    custom_bus: false,
    custom_library: false,
    custom_car_parks: false,
    custom_water_fountain: false,
    custom_toilet: false,
  });

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

  interface MapLoadEvent {
    target: mapboxgl.Map;
  }

  // Load custom icons to the map from the ImageConfig array
  // Note, if you add more icons, you need to add them to the IMAGES array and also edit LoadImages
  const loadImages = async (map: mapboxgl.Map) => {
    const loadImage = (config: ImageConfig): Promise<void> => {
      if (map.hasImage(config.id)) return Promise.resolve();

      return new Promise((resolve, reject) => {
        map.loadImage(window.location.origin + config.path, (error, image) => {
          if (error) reject(error);
          if (image) {
            map.addImage(config.id, image);
            setImagesLoaded(prev => ({ ...prev, [config.id]: true }));
          }
          resolve();
        });
      });
    };

    try {
      await Promise.all(IMAGES.map(loadImage));
    } catch (error) {
      console.error('Error loading images:', error);
    }
  };

  const handleMapLoad = (event: MapLoadEvent) => {
    const map = event.target;

    loadImages(map).catch(error => console.error('Error loading images:', error));
  };

  const [amenitiesFilter, setAmenitiesFilter] = useState<string[]>(() =>
    MultiSelectOptions.map((option) => option.value)
  );

  const { sessionToken } = useSession()

  const handleIconClick = (value: string) => {
    setAmenitiesFilter((prev) =>
      prev.includes(value) ? prev.filter((item) => item !== value) : [...prev, value]
    );
  };

  // Handle map click event
  const handleMapClick = (event: unknown) => {
    const mapClickEvent = event as MapClickEvent; // Type assertion
    if (mapClickEvent.coordinate) {
      const [longitude, latitude] = mapClickEvent.coordinate;
      setCoordinates({ latitude, longitude });
      setIsMarkerVisible(true);
    }
  };

  function convertToGeoJson(points: Point[]): Record<string, GeoJSON> {
    const featureTypes: GeoJsonCollection[] = [
      "parking_meter",
      "bike_stand",
      "public_wifi_access_point",
      "library",
      "multistorey_car_parking",
      "drinking_water_fountain",
      "public_toilet",
      "bike_sharing_station",
      "parking",
      "accessible_parking",
      "public_bins",
      "coach_parking",
    ];

    const geoJsonCollection: Record<string, GeoJSON> = {};
    featureTypes.forEach((type: string) => {
      geoJsonCollection[type] = {
        type: "FeatureCollection",
        features: points
          ?.filter((point) => point.Type === type)
          .map((point) => ({
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
    });

    return geoJsonCollection;
  }

  const fetchPointsFromDB = async (
    longitude: number,
    latitude: number,
    sliderValue: number,
    amenitiesFilter: string[] = []
  ) => {
    const response = await fetch(
      `/api/points?long=${longitude}&lat=${latitude}&radius=${sliderValue * 100
      }&types=${amenitiesFilter.join(",")}`,
      {
        method: "GET",
        credentials: "include",
        headers: {
          authorization: "Bearer " + sessionToken,
        }
      }
    );

    const data = await response.json();
    const geoJson = convertToGeoJson(data?.response?.content);

    setPointsGeoJson(geoJson);
  };

  const { startOnborda } = useOnborda();
  const handleStartOnborda = () => {
    startOnborda("general-onboarding");
  };

  const version = packageJson.version;

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
      sliderValue,
      amenitiesFilter
    );

    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [coordinates, sliderValue, amenitiesFilter]);

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
      {/* Onboarding help button */}
      <div className="absolute bottom-[5%] left-[1%] z-[999]">
        <div>
          <button
            onClick={handleStartOnborda}
            className="mt-2 px-4 py-2 bg-white text-gray-800 rounded-full shadow-md"
            id="onboarding-step-6"
          >
            {"?"}
          </button>
        </div>
      </div>
      {/* Map Container - Taller on mobile */}
      <div
        className="
          flex-grow
          h-full
          relative
          sm:h-[70vh]
          lg:h-screen relative
        "
        id="onboarding-st5"
      >
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
              onLoad={handleMapLoad}
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
              <MapSources
                pointsGeoJson={pointsGeoJson}
                imagesLoaded={imagesLoaded}
                amenitiesFilter={amenitiesFilter}
              />
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
      <div className="
        flex-none
        p-3
        bg-gray-50 
        overflow-y-auto
        xl:p-6
        lg:h-screen relative
        lg:p-6
        sm:p-4
      ">
        <div className="space-y-3 sm:space-y-4 lg:space-y-6 max-w-lg mx-auto lg:max-w-none">
          {mapBoxApiKey ? (
            <>
              <div className="px-2 sm:px-3 lg:px-4">
                <div className="flex items-center space-x-4" id="onboarding-step-1">
                  <Image
                    src="/images/BKlogo.svg"
                    alt="BK Logo"
                    width={64}
                    height={64}
                    className="inline-block"
                  />
                  <div>
                    <h1 className="text-lg sm:text-xl lg:text-2xl font-bold text-gray-900 tracking-tight">
                      Magpie Dashboard: <span className="text-[#3e6e96]">v{version}</span>
                    </h1>
                    <span className="italic">Services at a glance</span>
                  </div>
                </div>
              </div>
              {/* Search Radius Card */}
              <div className="sticky top-0 bg-gray-50 z-10 px-2 sm:px-3 lg:px-4">
                <div
                  className="w-full bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p4"
                  id="onboarding-step-2"
                >
                  <div className="space-y-2 sm:space-y-3">
                    <div>
                      <label className="text-sm lg:text-base font-medium text-gray-700 mb-2 block">
                        Search Radius
                      </label>
                      <Slider
                        value={[sliderValueDisplay]}
                        onValueChange={(value) => setSliderValueDisplay(value[0])}
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
                    <div className="flex space-x-2">
                      <button
                        onClick={() => {
                          setSliderValueDisplay(1);
                          setSliderValue(1);
                        }}
                        className="px-2 py-1 bg-gray-200 rounded"
                      >
                        100m
                      </button>
                      <button
                        onClick={() => {
                          setSliderValueDisplay(2);
                          setSliderValue(2);
                        }}
                        className="px-2 py-1 bg-gray-200 rounded"
                      >
                        200m
                      </button>
                      <button
                        onClick={() => {
                          setSliderValueDisplay(5);
                          setSliderValue(5);
                        }}
                        className="px-2 py-1 bg-gray-200 rounded"
                      >
                        500m
                      </button>
                      <button
                        onClick={() => {
                          setSliderValueDisplay(10);
                          setSliderValue(10);
                        }}
                        className="px-2 py-1 bg-gray-200 rounded"
                      >
                        1000m
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </>
          ) : null}
          {/* Combined Data and Filter Options Card */}
          <div className="px-2 sm:px-3 lg:px-4">
            <div
              className="bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4"
              id="onboarding-step-3"
            >
              {isMarkerVisible ? (
                <Suspense
                  fallback={<div className="animate-pulse">Loading...</div>}
                >
                  <div className="overflow-x-none">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Icon
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Amenity
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Count
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Show
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {MultiSelectOptions.map((option) => (
                          <tr
                            key={option.value}
                            className={`${!amenitiesFilter.includes(option.value)
                              ? 'bg-gray-100'
                              : ''
                              }`}
                          >
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              <Image
                                src={iconMap[option.value]}
                                alt={option.label}
                                width={24}
                                height={24}
                                className={`w-6 h-6 ${!amenitiesFilter.includes(option.value) ? 'filter grayscale' : ''}`}
                              />
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {option.label}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              {amenitiesFilter.includes(option.value)
                                ? (
                                    pointsGeoJson?.[option.value] as GeoJSON.FeatureCollection
                                  )?.features?.length > 0
                                  ? (
                                      <span className="font-bold">
                                        {(pointsGeoJson?.[option.value] as GeoJSON.FeatureCollection)
                                          ?.features?.length || 0}
                                      </span>
                                    )
                                  : (
                                      (pointsGeoJson?.[option.value] as GeoJSON.FeatureCollection)
                                        ?.features?.length || 0
                                    )
                                : '-'}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
                              <button onClick={() => handleIconClick(option.value)}>
                                {amenitiesFilter.includes(option.value) ? <Eye size={16} color="#3e6e96" /> : <EyeOff size={16} color="#3e6e96" />}
                              </button>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </Suspense>
              ) : (
                <div className="relative overflow-x-auto">
                  <div className="opacity-50 pointer-events-none">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Icon
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Amenity
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Count
                          </th>
                          <th
                            scope="col"
                            className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                          >
                            Show
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {MultiSelectOptions.map((option) => (
                          <tr key={option.value} className="bg-gray-100">
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              <Image
                                src={iconMap[option.value]}
                                alt={option.label}
                                width={24}
                                height={24}
                                className="w-6 h-6 filter grayscale"
                              />
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                              {option.label}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                              -
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
                              <EyeOff size={16} color="#3e6e96" />
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default LocationAggregatorMap;
