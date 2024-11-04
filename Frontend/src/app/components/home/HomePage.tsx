import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from "../ProtectedRoute";
import Nav from "@/app/nav/page";


const HomePage = () => {
  return (
      <ProtectedRoute>
        <LocationAggregatorMap />
        <div className="nav-container">
          <Nav />
        </div>
      </ProtectedRoute>
  );
};

export default HomePage;
