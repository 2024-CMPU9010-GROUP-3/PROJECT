"use client";

import { LocationData } from "@/lib/interfaces/types";
import { ColumnDef } from "@tanstack/react-table";
import { Button } from "@/components/ui/button";
import { CircleX } from "lucide-react";
import { Checkbox } from "@/components/ui/checkbox";

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
  // {
  //   header: "Radius (m)",
  //   accessorKey: "radius",
  // },
  {
    id: "actions",
    cell: ({ row }) => {
      let isDeleting = false;

      const handleDelete = async () => {
        if (isDeleting) return;
        isDeleting = true;

        try {
          const response = await fetch(`/api/history/${row.original.id}`, {
            method: "DELETE",
          });

          if (!response.ok) {
            throw new Error("Failed to delete");
          }
        } catch (error) {
          console.error("Delete failed:", error);
        } finally {
          isDeleting = false;
        }
      };

      return (
        <div className="flex gap-2">
          <Button
            className="bg-white text-red-600 hover:bg-red-200 hover:text-red-800 transition-all duration-500"
            size="sm"
            onClick={handleDelete}
            disabled={isDeleting}
          >
            <CircleX />
          </Button>
        </div>
      );
    },
  },
];
