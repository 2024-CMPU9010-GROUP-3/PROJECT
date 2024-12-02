// Haversine formula to calculate the distance between two points in meters
export const haversineDistance = (coords1: [number, number], coords2: [number, number]) => {
    const [lon1, lat1] = coords1;
    const [lon2, lat2] = coords2;

    const R = 6371000; // Radius of the Earth in meters
    const dLat = ((lat2 - lat1) * Math.PI) / 180;
    const dLon = ((lon2 - lon1) * Math.PI) / 180;

    const a =
      Math.sin(dLat / 2) * Math.sin(dLat / 2) +
      Math.cos((lat1 * Math.PI) / 180) * Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) * Math.sin(dLon / 2);

    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c; // Distance in meters
  };

  // Function to check if the mapHoverEvent is within 100 meters of a point
  export const isWithin20Meters = (hoverCoords: [number, number], pointCoords: [number, number]) => {
    const distance = haversineDistance(hoverCoords, pointCoords);
    return distance < 20;
  };