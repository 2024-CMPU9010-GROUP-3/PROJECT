"use client";

import { LocationData } from "@/lib/interfaces/types";
import { ColumnDef, Row } from "@tanstack/react-table";
import { Checkbox } from "@/components/ui/checkbox";
import React from "react";
import {Button} from "@/components/ui/button";
import {MapPinned} from "lucide-react";

export const getColumns = (handleShowOnMap : (row:Row<LocationData>) => void) : ColumnDef<LocationData>[] => {
  return [
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
      accessorKey: "datecreated",
      header: "Date Created",
      cell: ({ row }) =>
        new Date(row.getValue("datecreated")).toLocaleString(),
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
      header: "Show on Map",
      cell: ({ row }) => (
        <Button variant="outline" className="rounded-full w-10 h-10 p-0" onClick={() => {handleShowOnMap(row)}}>
          <MapPinned className="w-4 h-4"/>
        </Button>
      )
    }
  ]
}
