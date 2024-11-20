import React from 'react';
import { Source, Layer } from 'react-map-gl';
import { GeoJSON } from 'geojson';
import {
  parkingClusterStyles,
  parkingMeterLayers,
  bikeStandLayers,
  bikeSharingLayers,
  accessibleParkingLayers,
  publicBinLayers,
  publicWifiLayers,
  coachParkingLayers,
  libraryLayers,
  carParkLayers,
  waterFountainLayers,
  publicToiletLayers,
} from "./MapboxLayers";

interface MapSourcesProps {
  pointsGeoJson: Record<string, GeoJSON> | null;
  imagesLoaded: Record<string, boolean>;
  amenitiesFilter: string[];
}

const layerConfigs = [
  {
    id: "custom_parking",
    filterKey: "parking",
    dataKey: "parking",
    cluster: true,
    clusterMaxZoom: 14,
    clusterRadius: 50,
    layers: parkingClusterStyles,
  },
  {
    id: "custom_parking_meter",
    filterKey: "parking_meter",
    dataKey: "parking_meter",
    layers: parkingMeterLayers,
  },
  {
    id: "custom_bicycle",
    filterKey: "bike_stand",
    dataKey: "bike_stand",
    layers: bikeStandLayers,
  },
  {
    id: "custom_bicycle_share",
    filterKey: "bike_sharing_station",
    dataKey: "bike_sharing_station",
    layers: bikeSharingLayers,
  },
  {
    id: "custom_accessible_parking",
    filterKey: "accessible_parking",
    dataKey: "accessible_parking",
    layers: accessibleParkingLayers,
  },
  {
    id: "custom_public_bins",
    filterKey: "public_bins",
    dataKey: "public_bins",
    layers: publicBinLayers,
  },
  {
    id: "custom_public_wifi",
    filterKey: "public_wifi_access_point",
    dataKey: "public_wifi_access_point",
    layers: publicWifiLayers,
  },
  {
    id: "custom_bus",
    filterKey: "coach_parking",
    dataKey: "coach_parking",
    layers: coachParkingLayers,
  },
  {
    id: "custom_library",
    filterKey: "library",
    dataKey: "library",
    layers: libraryLayers,
  },
  {
    id: "custom_car_parks",
    filterKey: "multistorey_car_parking",
    dataKey: "multistorey_car_parking",
    layers: carParkLayers,
  },
  {
    id: "custom_water_fountain",
    filterKey: "drinking_water_fountain",
    dataKey: "drinking_water_fountain",
    layers: waterFountainLayers,
  },
  {
    id: "custom_toilet",
    filterKey: "public_toilet",
    dataKey: "public_toilet",
    layers: publicToiletLayers,
  },
];

const MapSources: React.FC<MapSourcesProps> = ({ pointsGeoJson, imagesLoaded, amenitiesFilter }) => {
  const sources = layerConfigs.map((config) => {
    if (imagesLoaded[config.id] && amenitiesFilter.includes(config.filterKey)) {
      return (
        <Source
          key={config.id}
          id={config.id}
          type="geojson"
          data={pointsGeoJson?.[config.dataKey]}
          {...(config.cluster !== undefined && { cluster: config.cluster })}
          {...(config.clusterMaxZoom !== undefined && { clusterMaxZoom: config.clusterMaxZoom })}
          {...(config.clusterRadius !== undefined && { clusterRadius: config.clusterRadius })}
        >
          <Layer {...config.layers?.close} />
          <Layer {...config.layers?.medium} />
          <Layer {...config.layers?.far} />
        </Source>
      );
    }
    return null;
  });

  return <>{sources}</>;
};

export default MapSources;