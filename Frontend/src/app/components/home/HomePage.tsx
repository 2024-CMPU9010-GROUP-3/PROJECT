"use client";
import dynamic from "next/dynamic";

// Dynamically import dynamic map components to prevent SSR issues
const GoogleMap = dynamic(() => import("../../components/map/GoogleMap"), {
  ssr: false,
});

const HomePage = () => {
  return (
    <div>
      <h1>Hello Project!!! test</h1>
      <h2>Google Map Integration</h2>
      <GoogleMap />
    </div>
  );
};

export default HomePage;
