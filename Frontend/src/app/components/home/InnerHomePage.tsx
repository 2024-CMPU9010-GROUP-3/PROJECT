import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from "../ProtectedRoute";
import Nav from "@/app/nav/page";
import CookieConsent from "../banner/CookieConsent";
import {Suspense} from "react";
import {Loader2} from "lucide-react";

const InnerHomePage = () => {
  return (
    <ProtectedRoute>
      <Suspense fallback={<Loader2/>}>
        <LocationAggregatorMap />
      </Suspense>
      <div className="nav-container">
        <div className="absolute top-5 left-5 z-[999]">
          <CookieConsent />
          <Nav />
        </div>
      </div>
    </ProtectedRoute>
  );
};

export default InnerHomePage;
