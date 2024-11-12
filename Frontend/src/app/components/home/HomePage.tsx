import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from "../ProtectedRoute";
import Nav from "@/app/nav/page";


const HomePage = () => {
  return (
    <ProtectedRoute>
      <LocationAggregatorMap />
      <div className="nav-container">
        <div className="absolute top-5 left-5 z-[999]">
          <Nav />
        </div>
      </div>
    </ProtectedRoute>
  );
};

export default HomePage;
