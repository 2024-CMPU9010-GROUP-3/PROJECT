"use client";

import { GoogleMap as GoogleMapComponent } from '@react-google-maps/api';
import { useRouter } from 'next/navigation';
import dynamic from 'next/dynamic';
import { useEffect } from 'react';

// Dynamically import dynamic map components to prevent SSR issues
const GoogleMap = dynamic(() => import('../../components/map/GoogleMap'), {
  ssr: false
});

const HomePage = () => {
  const router = useRouter();

  const handleButtonClick = () => {
    router.push('/googleMap');
  };

  return (
    <div>
      <h1>Hello Project!!!  test</h1>
      <h2>Google Map Integration</h2>
      <GoogleMap />
    </div>
  );
};

export default HomePage;
