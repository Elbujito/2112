import React, { useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { FiSearch } from "react-icons/fi";

import {
  createColumnHelper,
  useReactTable,
  getCoreRowModel,
  flexRender,
} from "@tanstack/react-table";

type TileSatelliteMapping = {
  ID: string;
  NoradID: string;
  TileID: string;
  TileCenterLat: string;
  TileCenterLon: string;
  TileZoomLevel: number;
  IntersectionLongitude: number;
  IntersectionLatitude: number;
  Intersection: {
    Longitude: number;
    Latitude: number;
  };
};

const columnHelper = createColumnHelper<TileSatelliteMapping>();

// SearchBar Component
function SearchBar({
  searchValue,
  onSearchChange,
  onSearchSubmit,
}: {
  searchValue: string;
  onSearchChange: (value: string) => void;
  onSearchSubmit: () => void;
}) {
  return (
    <div className="relative flex h-[61px] w-full items-center rounded-full bg-white px-2 py-2 shadow-xl dark:bg-navy-800">
      <button
        onClick={onSearchSubmit}
        className="flex h-full items-center justify-center rounded-l-full bg-lightPrimary text-navy-700 dark:bg-navy-900 px-4"
      >
        <FiSearch className="h-4 w-4 text-gray-400 dark:text-white" />
      </button>
      <input
        type="text"
        value={searchValue}
        onChange={(e) => onSearchChange(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && onSearchSubmit()}
        placeholder="Search..."
        className="block h-full w-full rounded-r-full bg-lightPrimary text-sm font-medium text-navy-700 outline-none dark:bg-navy-900 dark:text-white"
      />
    </div>
  );
}

// TableHeader Component
function TableHeader({ headers }: { headers: any[] }) {
  return (
    <thead className="sticky top-0 bg-white dark:bg-navy-800">
      {headers.map((headerGroup) => (
        <tr key={headerGroup.id}>
          {headerGroup.headers.map((header) => (
            <th
              key={header.id}
              colSpan={header.colSpan}
              className="py-2 text-left text-xs font-bold text-gray-600 dark:text-white"
            >
              {flexRender(
                header.column.columnDef.header,
                header.getContext()
              )}
            </th>
          ))}
        </tr>
      ))}
    </thead>
  );
}

// TableBody Component
function TableBody({
  rows,
  onRowClick,
  selectedRowId,
}: {
  rows: any[];
  onRowClick: (id: string) => void;
  selectedRowId: string | null;
}) {
  return (
    <tbody>
      {rows.map((row) => (
        <tr
          key={row.MappingID}
          onClick={() => onRowClick(row.original.MappingID)}
          className={`cursor-pointer ${selectedRowId === row.original.MappingID
            ? "bg-blue-100 dark:bg-blue-900"
            : "hover:bg-gray-100 dark:hover:bg-gray-700"
            }`}
        >
          {row.getVisibleCells().map((cell) => (
            <td key={cell.id} className="py-2 px-4 text-sm">
              {flexRender(cell.column.columnDef.cell, cell.getContext())}
            </td>
          ))}
        </tr>
      ))}
    </tbody>
  );
}

// Pagination Component
function Pagination({
  pageIndex,
  pageCount,
  onPageChange,
}: {
  pageIndex: number;
  pageCount: number;
  onPageChange: (pageIndex: number) => void;
}) {
  const handleNextPage = () => {
    if (pageIndex + 1 < pageCount) onPageChange(pageIndex + 1);
  };

  const handlePreviousPage = () => {
    if (pageIndex > 0) onPageChange(pageIndex - 1);
  };

  return (
    <div className="flex justify-between items-center py-4">
      <button
        onClick={handlePreviousPage}
        disabled={pageIndex === 0}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
      >
        <MdChevronLeft />
      </button>
      <p className="text-sm">
        Page {pageIndex + 1} of {pageCount}
      </p>
      <button
        onClick={handleNextPage}
        disabled={pageIndex + 1 >= pageCount}
        className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
      >
        <MdChevronRight />
      </button>
    </div>
  );
}

// Main MappingTableComponent
export default function MappingTableComponent(props: {
  tableData: TileSatelliteMapping[];
  totalItems: number;
  pageSize: number;
  pageIndex: number;
  onPageChange: (pageIndex: number) => void;
  searchValue: string;
  onSearchChange: (search: string) => void;
  onSearchSubmit: () => void;
  onRowClick: (id: string) => void;
}) {
  const {
    tableData,
    totalItems,
    pageSize,
    pageIndex,
    onPageChange,
    searchValue,
    onSearchChange,
    onSearchSubmit,
    onRowClick,
  } = props;

  const [selectedRowId, setSelectedRowId] = useState<string | null>(null);

  const handleRowClick = (id: string) => {
    setSelectedRowId(id);
    onRowClick(id);
  };

  const columns = [
    columnHelper.accessor("NoradID", {
      header: "NORAD ID",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("TileID", {
      header: "Tile ID",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor((row) => row.TileCenterLat, {
      header: "T. Center Lat.",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor((row) => row.TileCenterLon, {
      header: "T. Center Lon.",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor((row) => row.TileZoomLevel, {
      header: "T. Zoom Level",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor((row) => row.Intersection.Longitude, {
      header: "Inter. Longitude",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor((row) => row.Intersection.Latitude, {
      header: "Inter. Latitude",
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
  ];

  const table = useReactTable({
    data: tableData,
    columns,
    pageCount: Math.ceil(totalItems / pageSize),
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <Card extra="h-full w-full pb-4 px-4">
      <SearchBar
        searchValue={searchValue}
        onSearchChange={onSearchChange}
        onSearchSubmit={onSearchSubmit}
      />
      <div className="mt-4 overflow-auto">
        <table className="w-full">
          <TableHeader headers={table.getHeaderGroups()} />
          <TableBody
            rows={table.getRowModel().rows}
            onRowClick={handleRowClick}
            selectedRowId={selectedRowId}
          />
        </table>
      </div>
      <Pagination
        pageIndex={pageIndex}
        pageCount={Math.ceil(totalItems / pageSize)}
        onPageChange={onPageChange}
      />
    </Card>
  );
}
