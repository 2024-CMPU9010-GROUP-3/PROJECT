import { Navbar } from "@/components/global/Navbar";
import HistoryClient from "./HistoryClient";

const page = () => {
  return (
    <>
      <Navbar />
      <div className="py-4">
        <HistoryClient />
      </div>
    </>
  );
};
export default page;
