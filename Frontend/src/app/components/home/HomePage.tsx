"use client";

import { useRouter } from 'next/navigation';
import dynamic from 'next/dynamic';
import { useEffect } from 'react';

// Dynamically import the Map component to prevent SSR issues with mapbox-gl
const MapBox = dynamic(() => import('../../components/map/mapBox').then(mod => mod.default), { ssr: false });

const HomePage = () => {
  const router = useRouter();



  return (
    <div>
      <h1>Hello Project!!!  test</h1>
      
      <h1>Dublin Map</h1>
      <MapBox />

    </div>
  );
};

export default HomePage;
