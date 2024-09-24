"use client";

import { useRouter } from 'next/navigation';

const HomePage = () => {
  const router = useRouter();

  const handleButtonClick = () => {
    router.push('/googleMap');
  };

  return (
    <div>
      <h1>Hello Project!!!  test</h1>
      <button onClick={handleButtonClick}>Go to Google Map</button>
    </div>
  );
};

export default HomePage;
