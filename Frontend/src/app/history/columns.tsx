"use client";

import { LocationItem } from "@/lib/interfaces/types";
import { ColumnDef } from "@tanstack/react-table";

const formatKey = (key: string) => {
  return key
    .split(/(?=[A-Z])|_/)
    .map(
      (word: string) =>
        word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()
    )
    .join(" ");
};

export const columns: ColumnDef<LocationItem>[] = [
  {
    accessorKey: "datecreated",
    header: "Date",
    cell: ({ row }) => new Date(row.getValue("date")).toLocaleDateString(),
  },
  // {
  //   accessorKey: "location",
  //   header: "Location",
  // },
  {
    header: "Amenities",
    accessorKey: "amenities",
    cell: ({ row }) => {
      console.log("AMENITIES ROW>>>", row);

      const amenities = row.original.amenities;
      return (
        <div className="flex flex-wrap gap-2 text-sm">
          {Object.entries(amenities).map(([key, value]) => (
            <div
              key={key}
              className="flex items-center gap-2 bg-gray-50 rounded-md p-2 min-w-[80px] border border-gray-500"
            >
              <span className="font-medium">{formatKey(key)}:</span>
              <span className="text-gray-600">{value}</span>
            </div>
          ))}
        </div>
      );
    },
  },
  {
    header: "Coordinates",
    cell: ({ row }) => {
      const coords = row.original.coordinates;
      return `${coords.latitude}, ${coords.longitude}`;
    },
  },
  {
    header: "Actions",
    cell: ({ row }) => (
      <div className="flex gap-2">
        <button
          onClick={() =>
            window.open(
              `https://maps.google.com/?q=${row.original.coordinates.latitude},${row.original.coordinates.longitude}`
            )
          }
        >
          View on Map
        </button>
      </div>
    ),
  },
];
