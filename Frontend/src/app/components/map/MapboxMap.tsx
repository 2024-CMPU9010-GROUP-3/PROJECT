"use client";

// React core
import React, { Suspense, useEffect, useMemo, useState } from "react";

// Third-party packages
import DeckGL from "@deck.gl/react";
import { Feature, GeoJSON } from "geojson";
import { FaLocationDot } from "react-icons/fa6";
import { Grid } from "react-loader-spinner";
import Map, { Layer, LayerProps, Marker, Popup, Source } from "react-map-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import { Eye, EyeOff, Loader2, Save, Search } from "lucide-react";
import Image from "next/image";
import { useOnborda } from "onborda";
import { useSession } from "@/app/context/SessionContext";

// Local components
import { Slider } from "@/components/ui/slider";

// Local utils and configs
import { lightingEffect, INITIAL_VIEW_STATE } from "@/lib/mapconfig";
import { isWithin20Meters, haversineDistance } from "./utils/MeasurementUtils";
import { getCookiesAccepted } from "@/lib/cookies";
import { cn } from "@/lib/utils";
import packageJson from "../../../../package.json";
import MapSources from "./utils/MapSources";

// Types and interfaces
import {
  MapClickEvent,
  Coordinates,
  Point,
  CoordinatesForGeoJson,
  ImageConfig,
  GeoJsonCollection,
  MapHoverEvent,
} from "@/lib/interfaces/types";
import { Button } from "@/components/ui/button";
import { useToast } from "@/hooks/use-toast";
import { Toaster } from "@/components/ui/toaster";
import { Input } from "@/components/ui/input"
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import {useSearchParams} from "next/navigation";

type SliderProps = React.ComponentProps<typeof Slider>;

const mapElements = [
  {
    label: "Parking Meter",
    value: "parking_meter",
    id: "custom_parking_meter",
    path: "/mapicons/parking_meter.png",
  },
  {
    label: "Bike Stand",
    value: "bike_stand",
    id: "custom_bicycle",
    path: "/mapicons/bicycle.png",
  },
  {
    label: "Public Wifi",
    value: "public_wifi_access_point",
    id: "custom_public_wifi",
    path: "/mapicons/wifi.png",
  },
  {
    label: "Library",
    value: "library",
    id: "custom_library",
    path: "/mapicons/library.png",
  },
  {
    label: "Multi Storey Car Park",
    value: "multistorey_car_parking",
    id: "custom_car_parks",
    path: "/mapicons/car_park.png",
  },
  {
    label: "Drinking Water Fountain",
    value: "drinking_water_fountain",
    id: "custom_water_fountain",
    path: "/mapicons/water_fountain.png",
  },
  {
    label: "Public Toilet",
    value: "public_toilet",
    id: "custom_toilet",
    path: "/mapicons/toilet.png",
  },
  {
    label: "Bike Sharing Station",
    value: "bike_sharing_station",
    id: "custom_bicycle_share",
    path: "/mapicons/bicycle_share.png",
  },
  {
    label: "Parking",
    value: "parking",
    id: "custom_parking",
    path: "/mapicons/parking.png",
  },
  {
    label: "Accessible Parking",
    value: "accessible_parking",
    id: "custom_accessible_parking",
    path: "/mapicons/accessibleParking.png",
  },
  {
    label: "Public Bins",
    value: "public_bins",
    id: "custom_public_bins",
    path: "/mapicons/bin.png",
  },
  {
    label: "Coach Parking",
    value: "coach_parking",
    id: "custom_bus",
    path: "/mapicons/bus.png",
  },
];

const LocationAggregatorMap = ({ className, ...props }: SliderProps) => {
  const [mapBoxApiKey, setMapBoxApiKey] = useState<string>("");
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });
  const [markerIsVisible, setMarkerIsVisible] = useState<boolean>(false);
  const [pointsGeoJson, setPointsGeoJson] = useState<Record<
    string,
    GeoJSON
  > | null>(null);
  const [currentPositionCords, setCurrentPositionCords] = useState<Coordinates>(
    { latitude: 0, longitude: 0 }
  );
  const [sliderValue, setSliderValue] = useState<number>(40);
  const [sliderValueDisplay, setSliderValueDisplay] = useState<number>(40);
  const markerCoords = useMemo(
    () => [coordinates?.longitude, coordinates?.latitude],
    [coordinates]
  );

  const searchParams = useSearchParams();

  const [savingMap, setSavingMap] = useState(false);

  // Search state
  const [searchValue, setSearchValue] = useState<string>("");

  // Hover State
  const [hoverEntryTimeout, setEntryHoverTimeout] =
    useState<NodeJS.Timeout | null>(null);

  // Tooltip state
  const [toolTipIsVisible, setToolTipIsVisible] = useState<boolean>(false);
  const [toolTipX, setToolTipX] = useState<number>(0);
  const [toolTipY, setToolTipY] = useState<number>(0);
  const [toolTipContent, setToolTipContent] = useState<string>("");

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

  const [circleCoordinates, setCircleCoordinates] = useState<number[][]>([]);

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
            setImagesLoaded((prev) => ({ ...prev, [config.id]: true }));
          }
          resolve();
        });
      });
    };

    try {
      await Promise.all(mapElements.map(loadImage));
    } catch (error) {
      console.error("Error loading images:", error);
    }
  };

  const handleMapLoad = (event: MapLoadEvent) => {
    const map = event.target;

    loadImages(map).catch((error) =>
      console.error("Error loading images:", error)
    );
  };

  const [amenitiesFilter, setAmenitiesFilter] = useState<string[]>(() =>
    mapElements.map((option) => option.value)
  );

  const { sessionToken } = useSession();

  const handleIconClick = (value: string) => {
    setAmenitiesFilter((prev) =>
      prev.includes(value)
        ? prev.filter((item) => item !== value)
        : [...prev, value]
    );
  };

  const [resetSelection, setResetSelection] = useState(true);
  const handleGlobalAmenitiesFilter = () => {
    setAmenitiesFilter(() =>
      resetSelection ? [] : mapElements.map((option) => option.value)
    );
    setResetSelection((prev) => !prev);
  };

  // Handle map click event
  const handleMapClick = (event: unknown) => {
    const mapClickEvent = event as MapClickEvent; // Type assertion
    if (mapClickEvent.coordinate) {
      const [longitude, latitude] = mapClickEvent.coordinate;
      setCoordinates({ latitude, longitude });
      setMarkerIsVisible(true);
    }
  };

  // Handle map hover event
  const handleHover = (event: unknown) => {
    const mapHoverEvent = event as MapHoverEvent;

    // Clear any existing timeout
    if (hoverEntryTimeout) {
      clearTimeout(hoverEntryTimeout);
    }

    // Hide the tooltip
    setToolTipIsVisible(false);

    // Check if the feature is undefined, if so, exit early
    if (!mapHoverEvent.coordinate) {
      return;
    }

    // Check if the map is too zoomed out, if so, exit early
    if (mapHoverEvent.viewport.zoom < 14) {
      return;
    }

    // Set a new timeout
    const newTimeout = setTimeout(() => {
      if (pointsGeoJson) {
        let closestPoint: Feature | undefined;
        let minDistance = Infinity;

        const hoverCoords = mapHoverEvent.coordinate as [number, number];

        Object.values(pointsGeoJson).forEach((geoJson) => {
          (geoJson as GeoJSON.FeatureCollection).features.forEach((feature) => {
            if (
              feature.geometry.type === "Point" &&
              feature.geometry.coordinates.length === 2
            ) {
              const pointCoords = feature.geometry.coordinates as [
                number,
                number
              ];

              if (isWithin20Meters(hoverCoords, pointCoords)) {
                const distance = haversineDistance(hoverCoords, pointCoords);

                if (distance < minDistance) {
                  minDistance = distance;
                  closestPoint = feature;
                }
              }
            }
          });
        });

        if (closestPoint) {
          const pointGeometry = closestPoint.geometry as GeoJSON.Point;
          setToolTipX(pointGeometry.coordinates[1]);
          setToolTipY(pointGeometry.coordinates[0]);

          if (closestPoint.properties) {
            fetchPointById(closestPoint.properties.Id as number).then(
              (data) => {
                setToolTipContent(data?.response?.content);
              }
            );
          }

          setToolTipIsVisible(true);
          setToolTipContent("");
        }
      }
    }, 400); // 1000 milliseconds = 1 second

    // Store the new timeout ID in state
    setEntryHoverTimeout(newTimeout);
  };

  // Clear the timeout when the component unmounts
  // This prevents a client-side memory leak
  useEffect(() => {
    return () => {
      if (hoverEntryTimeout) {
        clearTimeout(hoverEntryTimeout);
      }
    };
  }, [hoverEntryTimeout]);

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

  const { sessionUUID } = useSession();
  const { toast } = useToast();

  const handleSaveMap = async () => {
    try {
      // Validate coordinates
      if (!coordinates?.latitude || !coordinates?.longitude) {
        throw new Error("Location coordinates are missing");
      }
      // Validate other required data
      if (!amenitiesFilter?.length) {
        throw new Error("Please select at least one amenity type");
      }
      setSavingMap(true);
      console.log();
      const locationName = await getNameFromLocation()
      const response = await fetch(`/api/history?userid=${sessionUUID}`, {
        method: "POST",
        body: JSON.stringify({
          amenitytypes: amenitiesFilter.map(value => {
            return{
              count: (
                pointsGeoJson?.[
                value
                ] as GeoJSON.FeatureCollection
              )?.features?.length || 0,
              type: value
            }
          }),
          longlat: {
            type: "Point",
            coordinates: [
              coordinates.longitude,
              coordinates.latitude,
            ],
          },
          radius: sliderValue * 100,
          displayName: locationName
        }),
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${sessionToken}`,
        },
      });
      if (!response.ok) {
        const errorData = await response.json().catch(() => null);
        throw new Error(errorData?.message || "Failed to save map");
      }
      toast({
        title: "Map saved",
        description: "Your map has been saved successfully",
      });
    } catch (error) {
      console.error("Error saving map:", error);
      toast({
        title: "Error saving map",
        description: "Failed to save map",
        variant: "destructive",
      });
    } finally {
      setSavingMap(false);
    }
  };

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
        },
      }
    );

    const data = await response.json();
    const geoJson = convertToGeoJson(data?.response?.content);

    setPointsGeoJson(geoJson);
  };

  const getLocationFromText = async () => {
    const response = await fetch(`/api/search?q=${searchValue}&limit=1`, {
      method: "GET",
      credentials: "include",
      headers: {
        authorization: "Bearer " + sessionToken,
      },
    });

    const data = await response.json();

    const foundLong = parseFloat(data[0].lon);
    const foundLat = parseFloat(data[0].lat);

    if (foundLat || foundLong) {
      setCoordinates({ latitude: foundLat, longitude: foundLong });
      setMarkerIsVisible(true);
    }
  };

  const getNameFromLocation = async () => {
    const response = await fetch(`/api/location-lookup?lat=${coordinates.latitude}&lon=${coordinates.longitude}&format=jsonv2`, {
      method: "GET",
      credentials: "include",
      headers: {
        authorization: "Bearer " + sessionToken,
      },
    });

    const data = await response.json();

    if(!data.address){
      return null;
    }

    const road = data.address.road;
    const city = data.address.city || data.address.town || data.address.village || data.address.municipality;
    const country = data.address.country;
    
    if(!road && !city){
      if(country){
        return `Somewhere in ${country}`;
      } else {
        return null;
      }
    }

    return [road, city].filter(x => !!x).join(", ")
  }


  const fetchPointById = async (id: number) => {
    const response = await fetch(`/api/details?id=${id}`, {
      method: "GET",
      credentials: "include",
      headers: {
        authorization: "Bearer " + sessionToken,
      },
    });

    if (!response.ok) {
      throw new Error("Network response was not ok");
    }

    const data = await response.json();
    return data;
  };

  const { startOnborda, closeOnborda } = useOnborda();

  const version = packageJson.version;

  useEffect(() => {
    const radiusInMeters = sliderValueDisplay * 100; // Convert slider value to meters
    const radiusInDegrees = radiusInMeters / 111320; // Convert meters to degrees (approximation)
    const numPoints = 64; // Number of points to define the circle
    const angleStep = (2 * Math.PI) / numPoints;
    const localCoordinates = [];

    for (let i = 0; i < numPoints; i++) {
      const angle = i * angleStep;
      const x =
        coordinates.longitude +
        (radiusInDegrees * Math.cos(angle)) /
        Math.cos(coordinates.latitude * (Math.PI / 180));
      const y = coordinates.latitude + radiusInDegrees * Math.sin(angle);
      localCoordinates.push([x, y]);
    }

    // Close the circle
    localCoordinates.push(localCoordinates[0]);

    setCircleCoordinates(localCoordinates);
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
  });

  useEffect(() => {
    const paramLon = searchParams.get("marker_long")
    const paramLat = searchParams.get("marker_lat")
    const paramRad = searchParams.get("marker_rad")
    const paramTypes = searchParams.get("marker_types")

    if(!paramLon || !paramLat || !paramRad || !paramTypes){
      return;
    }
    const latitude = parseFloat(paramLat)
    const longitude = parseFloat(paramLon)
    const radius = parseInt(paramRad)
    const types = paramTypes.split(",")
    
    if(longitude && latitude && radius && types) {
      // Reason for this timeout: coordinates getting set to 0,0 after this fires during map load
      // I think this happens only in the dev environment (strict mode)
      // If the points are not loading after loading from history, re-enable this timeout
      // setTimeout(() => {
        setCoordinates({ latitude, longitude });
        setSliderValue(radius / 100);
        setSliderValueDisplay(radius / 100);
        setMarkerIsVisible(true);
        setAmenitiesFilter(() =>
          types
        );
      //},1000);
    }
  }, [searchParams]);

  useEffect(() => {
    closeOnborda();
    if (sessionToken && getCookiesAccepted() === false) {
      // Onborda seems to load before the map, so we need to wait a bit before starting the onboarding
      // This is awful, terrible code, I am so sorry you have to see this. I know it's bad. I tried
      // MANY different ways of doing this, over several hours. This being the only workable solution.
      // Feel free to take a crack at it yourself, I'd love to see this removed.
      // 1Solon - 19/11/2024
      setTimeout(() => {
        startOnborda("general-onboarding");
      }, 1000);
    }
  }, [closeOnborda, sessionToken, startOnborda]);

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
      <Toaster />
      {/* Onboarding help button */}
      <div
        className="absolute bottom-[5%] left-[1%] z-[999]"
        id="onboarding-step-3"
      >
        <div>
          <button
            onClick={() => startOnborda("general-onboarding")}
            className="mt-2 px-4 py-2 bg-white text-gray-800 rounded-full shadow-md"
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
          lg:h-screen
        "
        id="onboarding-step-2"
      >
        {mapBoxApiKey ? (
          <DeckGL
            effects={[lightingEffect]}
            initialViewState={INITIAL_VIEW_STATE}
            controller={true}
            onClick={handleMapClick}
            onHover={handleHover}
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
              {markerIsVisible && (
                <Marker
                  longitude={coordinates?.longitude}
                  latitude={coordinates?.latitude}
                  anchor="center"
                >
                  <div>
                    <FaLocationDot size={50} color="FFA15A" />
                  </div>
                </Marker>
              )}
              <Marker
                latitude={currentPositionCords?.latitude}
                longitude={currentPositionCords?.longitude}
              >
                <div>
                  <FaLocationDot size={35} color="blue" />
                </div>
              </Marker>
              {toolTipIsVisible && (
                <Popup
                  latitude={toolTipX}
                  longitude={toolTipY}
                  closeButton={false}
                  style={{ whiteSpace: "pre-wrap", padding: "8px" }}
                  maxWidth="350px"
                  anchor="bottom"
                >
                  <div className="popup-content">
                    <h3 className="popup-header">Amenity Details</h3>
                    {Object.entries(toolTipContent).map(([key, value]) => (
                      <div className="key-value-pair" key={key}>
                        <span className="key">{key}:</span>
                        <span className="value">{value}</span>
                      </div>
                    ))}
                  </div>
                </Popup>
              )}
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
      <div
        className="
        flex-none
        p-3
        bg-gray-50 
        overflow-y-auto
        xl:p-6
        lg:h-screen relative
        lg:p-6
        sm:p-4
      "
        id="onboarding-step-1"
      >
        <div className="space-y-3 sm:space-y-3 lg:space-y-3 max-w-lg mx-auto lg:max-w-none">
          {mapBoxApiKey ? (
            <>
              <div className="px-2 sm:px-3 lg:px-4">
                <div className="flex items-center space-x-4">
                  <Image
                    src="/images/BKlogo.svg"
                    alt="BK Logo"
                    width={64}
                    height={64}
                    className="inline-block"
                  />
                  <div>
                    <h1 className="text-lg sm:text-xl lg:text-2xl font-bold text-gray-900 tracking-tight">
                      Magpie Dashboard:{" "}
                      <span className="text-[#3e6e96]">v{version}</span>
                    </h1>
                    <span className="italic">Services at a glance</span>
                  </div>
                </div>
              </div>
              {/* Search Radius Card */}
              <div className="sticky top-0 bg-gray-50 px-2 sm:px-3 lg:px-4">
                <div className="w-full bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p4">
                  <div className="space-y-2 sm:space-y-3">
                    <div>
                      <label className="text-sm lg:text-base font-medium text-gray-700 mb-2 block">
                        Search Radius
                      </label>
                      <Slider
                        value={[sliderValueDisplay]}
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
                    <div className="text-sm lg:text-base font-medium text-gray-600 flex space-x-2">
                      {sliderValueDisplay * 100} meters
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button
                              onClick={() => {
                                setCoordinates({ latitude: 0, longitude: 0 });
                                setMarkerIsVisible(false);
                              }}
                              className="px-2 py-1 bg-gray-200 hover:bg-gray-300 transition-all duration-200 rounded ml-auto"
                            >
                              Clear
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Clear the currently selected point</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button
                              onClick={() => {
                                setSliderValueDisplay(1);
                                setSliderValue(1);
                              }}
                              className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded"
                            >
                              100m
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Set search radius to 100m</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button
                              onClick={() => {
                                setSliderValueDisplay(2);
                                setSliderValue(2);
                              }}
                              className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded"
                            >
                              200m
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Set search radius to 200m</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button
                              onClick={() => {
                                setSliderValueDisplay(5);
                                setSliderValue(5);
                              }}
                              className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded"
                            >
                              500m
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Set search radius to 500m</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button
                              onClick={() => {
                                setSliderValueDisplay(10);
                                setSliderValue(10);
                              }}
                              className="px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded"
                            >
                              1000m
                            </button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Set search radius to 1000m</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                    <div className="text-sm lg:text-base font-medium text-gray-600 flex space-x-2">
                      <Input
                        className="mx-auto"
                        onChange={(e) => setSearchValue(e.target.value)}
                        type="text" id="search" placeholder="Location Search"
                      />
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              className="w-15 mx-auto bg-neutral-700 transition-all duration-200 hover:scale-[1.02] hover:shadow-md active:scale-[0.98]"
                              onClick={getLocationFromText}
                            >
                              <Search size={16} />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Search for location</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              disabled={savingMap}
                              className="w-15 mx-auto bg-neutral-700 transition-all duration-200 hover:scale-[1.02] hover:shadow-md active:scale-[0.98]"
                              onClick={handleSaveMap}
                            >
                              {savingMap ? <Loader2 size={16} className="animate-spin" /> : <Save size={16} />}
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>Save the currently selected information</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                  </div>
                </div>
              </div>
            </>
          ) : null}
          {/* Combined Data and Filter Options Card */}
          <div className="px-2 sm:px-3 lg:px-4">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-3 sm:p-4">
              <Suspense
                fallback={<div className="animate-pulse">Loading...</div>}
              >
                {mapBoxApiKey ? (
                  <div>
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
                            <button
                              onClick={() => handleGlobalAmenitiesFilter()}
                            >
                              {resetSelection ? (
                                <Eye size={16} color="#3e6e96" />
                              ) : (
                                <EyeOff size={16} color="#3e6e96" />
                              )}
                            </button>
                          </th>
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {mapElements.map((option) => {
                          const imageConfig = mapElements.find(
                            (img) => img.value === option.value
                          );
                          return (
                            <tr
                              key={option.value}
                              className={`${!amenitiesFilter.includes(option.value)
                                ? "bg-gray-100"
                                : ""
                                }`}
                            >
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {imageConfig && (
                                  <Image
                                    src={imageConfig.path}
                                    alt={option.label}
                                    width={24}
                                    height={24}
                                    className={`w-6 h-6 ${!amenitiesFilter.includes(option.value)
                                      ? "filter grayscale"
                                      : ""
                                      }`}
                                  />
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                {option.label}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                {amenitiesFilter.includes(option.value) ? (
                                  (
                                    pointsGeoJson?.[
                                    option.value
                                    ] as GeoJSON.FeatureCollection
                                  )?.features?.length > 0 ? (
                                    <span className="font-bold">
                                      {(
                                        pointsGeoJson?.[
                                        option.value
                                        ] as GeoJSON.FeatureCollection
                                      )?.features?.length || 0}
                                    </span>
                                  ) : (
                                    (
                                      pointsGeoJson?.[
                                      option.value
                                      ] as GeoJSON.FeatureCollection
                                    )?.features?.length || 0
                                  )
                                ) : (
                                  "-"
                                )}
                              </td>
                              <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">
                                <button
                                  onClick={() => handleIconClick(option.value)}
                                >
                                  {amenitiesFilter.includes(option.value) ? (
                                    <Eye size={16} color="#3e6e96" />
                                  ) : (
                                    <EyeOff size={16} color="#3e6e96" />
                                  )}
                                </button>
                              </td>
                            </tr>
                          );
                        })}
                      </tbody>
                    </table>
                  </div>
                ) : null}
              </Suspense>
            </div>
          </div>
        </div>
      </div>
    </div >
  );
};

export default LocationAggregatorMap;
