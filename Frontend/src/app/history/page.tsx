// import ProtectedRoute from "../components/ProtectedRoute";
import HistoryClient from "./HistoryClient";

const Page = async () => {
  return (
    <>
      <div className="min-h-screen bg-gray-50">
        <HistoryClient />
      </div>
    </>
  );
};

export default Page;
