"use client";

import { LocationData } from "@/lib/interfaces/types";
import { ColumnDef } from "@tanstack/react-table";

export const columns: ColumnDef<LocationData>[] = [
  {
    accessorKey: "id",
    header: "ID",
  },
  {
    accessorKey: "datecreated",
    header: "Date Created",
    cell: ({ row }) =>
      new Date(row.getValue("datecreated")).toLocaleDateString(),
  },
  {
    header: "Amenity Types",
    accessorKey: "amenitytypes",
    cell: ({ row }) => (
      <div className="flex flex-wrap gap-2">
        {row.getValue<string[]>("amenitytypes").map((type) => (
          <span
            key={type}
            className="bg-gray-50 rounded-md p-2 text-sm border border-gray-500"
          >
            {type.replace("_", " ")}
          </span>
        ))}
      </div>
    ),
  },
  {
    header: "Coordinates",
    cell: ({ row }) => {
      const longlat = row.original.longlat;
      return `${longlat.coordinates[1]}, ${longlat.coordinates[0]}`;
    },
  },
  {
    header: "Radius (m)",
    accessorKey: "radius",
  },
  {
    header: "Actions",
    cell: ({ row }) => (
      <div className="flex gap-2">
        <button
          className="text-blue-600 hover:text-blue-800"
          onClick={() => {
            const coords = row.original.longlat.coordinates;
            window.open(`https://maps.google.com/?q=${coords[1]},${coords[0]}`);
          }}
        >
          View on Map
        </button>
      </div>
    ),
  },
];
