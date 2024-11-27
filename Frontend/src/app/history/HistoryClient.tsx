"use client";

import { DataTable } from "./data-table";
import { columns } from "./columns";
import { useEffect, useState } from "react";
import { useSession } from "../context/SessionContext";

const HistoryClient = () => {
  const [history, setHistory] = useState([]);
  const { sessionToken, sessionUUID } = useSession();
  const fetchHistory = async (sessionToken: string) => {
    if (sessionToken) {
      const response = await fetch(`/api/history?userid=${sessionUUID}`, {
        headers: {
          Authorization: `Bearer ${sessionToken}`,
        },
      });
      const data = await response.json();

      setHistory(data?.response?.content?.history);
    }
  };

  useEffect(() => {
    fetchHistory(sessionToken);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [sessionToken]);

  return (
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
            <DataTable columns={columns} data={history ?? []} />
          </div>
        </div>
      </div>
    </div>
  );
};
export default HistoryClient;
