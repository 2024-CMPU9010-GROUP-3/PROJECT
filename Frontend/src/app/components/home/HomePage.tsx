import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from "../ProtectedRoute";

const HomePage = () => {
  return (
      <ProtectedRoute>
        <LocationAggregatorMap />
      </ProtectedRoute>
  );
};

export default HomePage;
