import LocationAggregatorMap from "../map/MapboxMap";

// Dynamically import dynamic map components to prevent SSR issues
// const GoogleMap = dynamic(() => import('../../components/map/GoogleMap'), {
//   ssr: false
// });

const HomePage = () => {
  // const router = useRouter();

  // const handleButtonClick = () => {
  //   router.push('/googleMap');
  // };

  return (
    <div>
      <h1>Hello Project!!! test</h1>
      <h2>Google Map Integration</h2>
      {/* <GoogleMap /> */}
      <LocationAggregatorMap />
    </div>
  );
};

export default HomePage;
