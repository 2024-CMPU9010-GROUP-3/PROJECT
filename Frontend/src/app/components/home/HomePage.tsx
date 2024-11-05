import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from "../ProtectedRoute";
import Nav from "@/app/nav/page";
import { CookieConsent } from "../banner/CookieConsent";


const HomePage = () => {
  return (
      <ProtectedRoute>
        <LocationAggregatorMap />
        <div className="nav-container">
            
  <CookieConsent />
          <Nav />
        </div>
      </ProtectedRoute>
  );
};

export default HomePage;
