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
      header: "Location",
      cell: ({ row }) => {
        const name = row.original.displayname;
        const longlat = row.original.longlat;
        if(name) {
          return (
            <div>
              <div className="font-bold">{name}</div>
              <div className="font-mono text-neutral-500 text-xs size-fit whitespace-nowrap">
                <div>lat: {longlat.coordinates[1].toFixed(6)}</div>
                <div>lng: {longlat.coordinates[0].toFixed(6)}</div>
              </div>
            </div>
          )
        }
        return (
        <div className="font-mono  whitespace-nowrap">
          <div>lat: {longlat.coordinates[1].toFixed(6)}</div>
          <div>lng: {longlat.coordinates[0].toFixed(6)}</div>
        </div>);
      },
    },
    {
      accessorKey: "radius",
      header: "Radius",
      cell: ({ row }) => {
        const radius : number = row.getValue("radius")
        return (
        <div className="w-fit whitespace-nowrap">
          {radius.toLocaleString()} m
        </div>)

      }
    },
    {
      header: "Amenity Types",
      accessorKey: "amenitytypes",
      cell: ({ row }) => (
        <div className="flex flex-wrap gap-2">
          {row.getValue<LocationData["amenitytypes"]>("amenitytypes")?
          row.getValue<LocationData["amenitytypes"]>("amenitytypes").map((type) => (
            <div key={type.type} className=" rounded-md py-1 px-2 text-s border border-gray-200 flex gap-1 flex-row flex-nowrap">
              <div className="font-bold">
                {type.type.split("_").map(word => word.charAt(0).toUpperCase() + word.substring(1)).join(' ') + ":"}
              </div>
              <div>
                {type.count}
              </div>
            </div>
          )):<div/>}
        </div>
      ),
    },
    {
      accessorKey: "datecreated",
      header: "Date Created",
      cell: ({ row }) =>
        new Date(row.getValue("datecreated")).toLocaleString(),
    },
    {
      header: "Show on Map",
      cell: ({ row }) => (
        <div className="flex items-center justify-center">
          <Button variant="outline" className="rounded-full w-10 h-10 p-0" onClick={() => {handleShowOnMap(row)}}>
            <MapPinned className="w-4 h-4"/>
          </Button>
        </div>
      )
    }
  ]
}
