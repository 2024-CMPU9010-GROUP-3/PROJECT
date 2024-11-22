import { DataTable } from "./data-table";
import { columns } from "./columns";
import { LocationItem } from "@/lib/interfaces/types";

const getData = () => {
  const locationData: LocationItem[] = [
    {
      id: 1,
      date: "2024-10-20",
      amenities: {
        parkingMeters: 1000,
        bikeStand: 1053,
        publicWiFi: 22,
        library: 16,
        multiStoreyCar: 16,
        drinkingWater: 1,
        publicToilet: 17,
        bikeStationSharing: 110,
        parking: 83788,
        accessibleParking: 231,
        publicBins: 2214,
        coachingPark: 2214,
      },
      location: "Phibsborough",
      coordinates: {
        longitude: 53.362854,
        latitude: 53.362854,
      },
    },
    {
      id: 2,
      date: "2024-10-20",
      amenities: {
        parkingMeters: 1000,
        bikeStand: 1053,
        publicWiFi: 22,
        library: 16,
        bikeStationSharing: 110,
        parking: 83788,
        accessibleParking: 11,
      },
      location: "Dublin 2",
      coordinates: {
        longitude: 6.26139,
        latitude: 53.339214,
      },
    },
    {
      id: 3,
      date: "2024-10-20",
      amenities: {
        parkingMeters: 1000,
        bikeStand: 1053,
        publicWiFi: 22,
        library: 16,
        bikeStationSharing: 110,
        parking: 83788,
      },
      location: "Dublin 1",
      coordinates: {
        longitude: 53.339214,
        latitude: 6.26139,
      },
    },
    {
      id: 4,
      date: "2024-10-20",
      amenities: {
        parkingMeters: 1000,
        bikeStand: 1053,
        publicWiFi: 22,
        library: 16,
        bikeStationSharing: 110,
      },
      location: "Dublin 7",
      coordinates: {
        longitude: 53.339214,
        latitude: 6.26139,
      },
    },
  ];
  return locationData;
};

const page = async () => {
  const data = await getData();
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-[100rem] mx-auto p-6 space-y-6">
        {/* Modern Header */}
        <div className="bg-white rounded-xl shadow-sm p-8">
          <h1 className="text-3xl font-semibold text-gray-900 mb-2">History</h1>
          <p className="text-gray-600">
            View and manage your past location searches and saved amenities.
          </p>
        </div>

        {/* Main Content Area */}
        <div className="bg-white rounded-xl shadow-sm p-8">
          <div className="flex flex-col space-y-6">
            <div>
              <h2 className="text-xl font-semibold text-gray-900">
                Welcome back!
              </h2>
              <p className="text-gray-600 mt-1">
                Here&apos;s a list of your saved locations for amenities!
              </p>
            </div>

            {/* Keep your existing DataTable component here */}
            {/* The rest of your content remains unchanged */}
            <div className="relative overflow-hidden rounded-lg border border-gray-200 w-full min-h-[calc(100vh-20rem)]">
              <div className="overflow-x-auto h-full">
                <DataTable columns={columns} data={data} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default page;
