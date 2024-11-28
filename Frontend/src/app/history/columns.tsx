"use client";

import { LocationData } from "@/lib/interfaces/types";
import { ColumnDef } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";
import React from "react";

export const columns: ColumnDef<LocationData>[] = [
  {
    id: "select",
    header: ({ table }) => (
      <div className="flex items-center justify-center h-8 w-8">
        <Checkbox
          checked={
            table.getIsAllPageRowsSelected() ||
            (table.getIsSomePageRowsSelected() && "indeterminate")
          }
          onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
          aria-label="Select all"
          className="p-2"
        />
      </div>
    ),
    cell: ({ row }) => (
      <div className="flex items-center justify-center h-8 w-8">
        <Checkbox
          checked={row.getIsSelected()}
          onCheckedChange={(value) => row.toggleSelected(!!value)}
          aria-label="Select row"
          className="p-2"
        />
      </div>
    ),
    enableSorting: false,
    enableHiding: false,
  },
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
];
