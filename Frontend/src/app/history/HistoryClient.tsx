"use client";

import { DataTable } from "./data-table";
import { columns } from "./columns";
import { useEffect, useState } from "react";
import { useSession } from "../context/SessionContext";
import { Button } from "@/components/ui/button";
import { LocationData } from "@/lib/interfaces/types";

const HistoryClient = () => {
  const [history, setHistory] = useState<LocationData[]>([]);
  const [isDeleting, setIsDeleting] = useState(false);
  const [rowSelection, setRowSelection] = useState({});
  const { sessionToken, sessionUUID } = useSession();
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

  useEffect(() => {
    fetchHistory(sessionToken);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sessionToken]);

  useEffect(() => {
    console.log("ROW SELECTION", rowSelection);
  }, [rowSelection]);

  return (
    // Adjust main container padding for smaller screens
    <div className="max-w-[100rem] mx-auto p-3 sm:p-4 md:p-6 space-y-4 sm:space-y-6">
      {/* Responsive Header */}
      <div className="bg-white rounded-lg sm:rounded-xl shadow-sm p-4 sm:p-6 md:p-8">
        <h1 className="text-2xl sm:text-3xl font-semibold text-gray-900 mb-1 sm:mb-2">
          History
        </h1>
        <p className="text-sm sm:text-base text-gray-600">
          View and manage your past location searches and saved amenities.
        </p>
      </div>

      {/* Responsive Main Content Area */}
      <div className="bg-white rounded-lg sm:rounded-xl shadow-sm p-4 sm:p-6 md:p-8">
        <div className="flex flex-col space-y-4 sm:space-y-6">
          <div>
            <h2 className="text-lg sm:text-xl font-semibold text-gray-900">
              Welcome back!
            </h2>
            <p className="text-sm sm:text-base text-gray-600 mt-1">
              Here&apos;s a list of your saved locations for amenities!
            </p>
          </div>

          {/* Keep your existing DataTable component here */}
          {/* The rest of your content remains unchanged */}
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
          <div className="relative overflow-hidden rounded-lg border border-gray-200 w-full min-h-[calc(100vh-20rem)]">
            <DataTable
              columns={columns}
              data={history ?? []}
              setRowSelection={setRowSelection} // Add this prop
              rowSelection={rowSelection} // Add this prop
            />
          </div>
        </div>
      </div>
    </div>
  );
};
export default HistoryClient;
