import { useEffect, useState } from 'react';
import LocationAggregatorMap from "../map/MapboxMap";
import ProtectedRoute from '@/app/components/ProtectedRoute';
import { verifySession } from '@/lib/dal';

const HomePage = () => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);

  useEffect(() => {
    const checkSession = async () => {
      const session = await verifySession();
      setIsAuthenticated(!!session);
    };

    checkSession();
  }, []);

  return (
    <ProtectedRoute>
      <div>
        {isAuthenticated ? (
          <LocationAggregatorMap />
        ) : (
          <p>请登录以查看地图</p>
        )}
      </div>
    </ProtectedRoute>
  );
};

export default HomePage;
