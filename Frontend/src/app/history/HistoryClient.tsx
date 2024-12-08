"use client";

import { DataTable } from "./data-table";
import { getColumns } from "./columns";
import { useEffect, useState } from "react";
import { useSession } from "../context/SessionContext";
import { Button } from "@/components/ui/button";
import { LocationData } from "@/lib/interfaces/types";
import { ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { Row } from "@tanstack/react-table";

const HistoryClient = () => {
  const [history, setHistory] = useState<LocationData[]>([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const [rowSelection, setRowSelection] = useState({});
  const { sessionToken, sessionUUID } = useSession();
  const router = useRouter();
  const fetchHistory = async (sessionToken: string) => {
    if (sessionToken) {
      const response = await fetch(`/api/history?userid=${sessionUUID}`, {
        method: "GET",
        headers: {
          Authorization: `Bearer ${sessionToken}`,
        },
      });
      const data = await response.json();

      setHistory(data?.response?.content?.history);
    }
  };

  const handleDelete = async (deleteArray: number[]) => {
    setIsDeleting(true);
    try {
      const response = await fetch(`/api/history?userid=${sessionUUID}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${sessionToken}`,
        },
        body: JSON.stringify({
          idlist: deleteArray,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to delete");
      }

      // Reset selection and refetch data
      setRowSelection({});
      await fetchHistory(sessionToken);
    } catch (error) {
      console.error("Delete failed:", error);
    } finally {
      setIsDeleting(false);
    }
  };

  const handleShowOnMap = (row: Row<LocationData>) => {
    const longlat = row.original.longlat;
    router.push(
      `/home?marker_long=${longlat.coordinates[0]}&marker_lat=${
        longlat.coordinates[1]
      }&marker_rad=${
        row.original.radius
      }&marker_types=${row?.original?.amenitytypes?.map((entry) => entry.type)}`
    );
  };

  const columns = getColumns(handleShowOnMap);

  useEffect(() => {
    fetchHistory(sessionToken);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sessionToken]);

  return (
    <div className="max-w-[100rem] mx-auto p-6 space-y-6">
      <Button
        className="absolute top-5 left-5 z-[999] h-10 w-10 bg-white rounded-full p-2 hover:bg-neutral-100"
        onClick={() => {
          router.push("/home");
        }}
      >
        <ArrowLeft color="black" className="w-full h-full" />
      </Button>
      {/* Modern Header */}
      <div className="bg-white rounded-xl shadow-sm p-8">
        <h1 className="text-3xl font-semibold text-gray-900 mb-2">
          Saved Locations
        </h1>
        <p className="text-gray-600">
          Here&apos;s a list of your saved amenity locations!
        </p>
      </div>

      {/* Main Content Area */}
      <div className="flex flex-col space-y-6">
        <div className="flex w-full justify-end">
          <Button
            className="bg-white text-red-700 hover:bg-red-500 hover:text-white transition-colors duration-300"
            onClick={() => {
              // Get selected row IDs
              const selectedIds = Object.keys(rowSelection).map(
                (idx) => history[parseInt(idx)].id
              );
              handleDelete(selectedIds);
            }}
            disabled={Object.keys(rowSelection).length === 0 || isDeleting}
          >
            Delete
          </Button>
        </div>
        <DataTable
          columns={columns}
          data={history}
          setRowSelection={setRowSelection} // Add this prop
          rowSelection={rowSelection} // Add this prop
        />
      </div>
    </div>
  );
};
export default HistoryClient;
