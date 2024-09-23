import dynamic from 'next/dynamic';


// Dynamically import the map to prevent SSR issues
const GoogleMap = dynamic(() => import('./GoogleMap'), {
  ssr: false
});

const Home = () => {
  return (
    <div>
      <h1>Google Map Integration</h1>
      <GoogleMap />
    </div>
  );
};

export default Home;