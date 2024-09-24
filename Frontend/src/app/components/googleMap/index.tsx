import dynamic from 'next/dynamic';



const GoogleMap = dynamic(() => import('./GoogleMap'), {
  ssr: false
});

const GoogleMapPage = () => {
  return (
    <div>
      <h1>Google Map Integration</h1>
      <GoogleMap />
    </div>
  );
};

export default GoogleMapPage;