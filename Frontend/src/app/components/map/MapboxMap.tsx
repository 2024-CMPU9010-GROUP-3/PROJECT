"use client";

import React, { useState } from "react";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import Map, { Marker } from "react-map-gl";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { HexagonLayer } from "@deck.gl/aggregation-layers";
import { FaLocationDot } from "react-icons/fa6";
import DeckGL from "@deck.gl/react";
import "mapbox-gl/dist/mapbox-gl.css";
import { MapClickEvent } from "@/lib/interfaces/types";
import { Coordinates } from "@/lib/interfaces/types";

import { lightingEffect, INITIAL_VIEW_STATE } from "@/lib/mapconfig";
import { FloatingDock } from "@/components/ui/floating-dock";
import { IconHome } from "@tabler/icons-react";
import { Slider } from "@/components/ui/slider";
import { cn } from "@/lib/utils";
import { Input } from "@/components/ui/input";


type SliderProps = React.ComponentProps<typeof Slider>;

const LocationAggregatorMap = ({ className, ...props }: SliderProps) => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [coordinates, setCoordinates] = useState<Coordinates>({
    latitude: 0,
    longitude: 0,
  });

  const [sliderValue, setSliderValue] = useState<number>(0);

  const links = [
    {
      title: "Home",
      icon: (
        <IconHome className="h-full w-full text-neutral-500 dark:text-neutral-300" />
      ),
      href: "#",
    },

    {
      title: "Home",
      icon: (
        <IconHome className="h-full w-full text-neutral-500 dark:text-neutral-300" />
      ),
      href: "#",
    },
    {
      title: "Home",
      icon: (
        <IconHome className="h-full w-full text-neutral-500 dark:text-neutral-300" />
      ),
      href: "#",
    },
  ];

  // Handle map click event
  const handleMapClick = (event: unknown) => {
    const mapClickEvent = event as MapClickEvent; // Type assertion
    if (mapClickEvent.coordinate) {
      const [longitude, latitude] = mapClickEvent.coordinate;
      setCoordinates({ latitude, longitude });
      console.log(`Map clicked at longitude: ${longitude}, latitude: ${latitude}`);
    }
  };

  return (
    <div>
      <DeckGL
        effects={[lightingEffect]}
        initialViewState={INITIAL_VIEW_STATE}
        controller={true}
        onClick={handleMapClick}
      >
        <Map
          mapboxAccessToken={process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN}
          mapStyle="mapbox://styles/mapbox/streets-v8"
          antialias={true}
        >
          <Marker
            longitude={coordinates?.longitude}
            latitude={coordinates?.latitude}
            anchor="bottom"
          >
            <FaLocationDot size={50} color="FFA15A" />
          </Marker>
        </Map>
      </DeckGL>
      <div className="absolute bg-transparent top-20 right-32 scale-125 transition-all">
        <div className="2xl:min-w-[200px] 2xl:max-w-[300px] bg-white rounded-xl ">
          <div className="px-2 py-4 space-y-2">
            <div className="space-y-1">
              <label>Distance</label>
              <Slider
                onValueChange={(value) => setSliderValue(value[0])}
                defaultValue={[0]}
                max={100}
                step={1}
                className={cn("w-[60%]", className)}
                {...props}
              />
            </div>
            <div className="">
              <span className="p-1">{sliderValue}</span>
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
      <div className="absolute bg-transparent bottom-10 right-20 scale-105 z-20">
        <FloatingDock items={links} desktopClassName="bg-transparent " />
      </div>
    </div>
  );
};

export default LocationAggregatorMap;
