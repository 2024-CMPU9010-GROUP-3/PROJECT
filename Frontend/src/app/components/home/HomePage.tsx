import LocationAggregatorMap from "../map/MapboxMap";
import CookieBanner from '@/app/components/banner/CookieBanner';


const HomePage = () => {
  return (
    <div>
      <LocationAggregatorMap />
      <CookieBanner />
    </div>
  );
};

export default HomePage;
