import React, { useEffect, useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { FiSearch } from "react-icons/fi";
import { BiTargetLock } from "react-icons/bi";
import { Tooltip } from "@chakra-ui/react"; // Assuming you're using Chakra UI for tooltips

import {
  createColumnHelper,
  useReactTable,
  getCoreRowModel,
  flexRender,
} from "@tanstack/react-table";

type Satellite = {
  Id: string;
  Name: string;
  NoradID: string;
  Owner: string;
  LaunchDate: string | null;
  Apogee: number | null;
  Perigee: number | null;
  TleUpdatedAt: string | null; // New field for TLE update time
  Action: string| null;
};

const columnHelper = createColumnHelper<Satellite>();

export default function SatelliteTableComponent(props: {
  tableData: Satellite[];
  totalItems: number;
  pageSize: number;
  pageIndex: number;
  onPageChange: (pageIndex: number) => void;
  onPageSizeChange: (pageSize: number) => void;
  searchValue: string;
  onSearchChange: (search: string) => void;
  onSearchSubmit: () => void;
  onRowClick: (noradID: string) => void; // Callback for emitting NoradID
}) {
  const {
    tableData,
    totalItems,
    pageSize,
    pageIndex,
    onPageChange,
    onPageSizeChange,
    searchValue,
    onSearchChange,
    onSearchSubmit,
    onRowClick, // Callback for emitting NoradID
  } = props;

  const [currentData, setCurrentData] = useState<Satellite[]>(tableData);

  // Effect to update current page data based on page index and size
  useEffect(() => {
    setCurrentData(tableData);
  }, [tableData, onPageChange, pageIndex]);

  const columns = [
    columnHelper.accessor("Name", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">NAME</p>
      ),
      cell: (info) => (
        <p className="text-navy-700 text-sm font-bold dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor("NoradID", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          NORAD ID
        </p>
      ),
      cell: (info) => (
        <p className="text-navy-700 text-sm font-bold dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor("Owner", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">OWNER</p>
      ),
      cell: (info) => (
        <p className="text-navy-700 text-sm font-bold dark:text-white">
          {info.getValue()}
        </p>
      ),
    }),
    columnHelper.accessor("LaunchDate", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          LAUNCH DATE
        </p>
      ),
      cell: (info) => (
        <p className="text-navy-700 text-sm font-bold dark:text-white">
          {info.getValue()
            ? new Date(info.getValue() as string).toLocaleDateString()
            : "N/A"}
        </p>
      ),
    }),
    columnHelper.accessor("TleUpdatedAt", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          TLE UPDATED AT
        </p>
      ),
      cell: (info) => (
        <p className="text-navy-700 text-sm font-bold dark:text-white">
          {info.getValue()
            ? new Date(info.getValue() as string).toLocaleString()
            : "N/A"}
        </p>
      ),
    }),
    columnHelper.accessor("Action", {
      header: () => (
        <p className="text-sm font-bold text-gray-600 dark:text-white">
          ACTION
        </p>
      ),
      cell: (info) => {
        const row = info.row.original;
        return row.TleUpdatedAt ? (
          <Tooltip
            label={`Last Updated: ${new Date(row.TleUpdatedAt).toLocaleString()}`}
            aria-label="TLE Update Tooltip"
          >
            <button
              className="text-blue-500 hover:text-blue-700"
              onClick={() => onRowClick(row.NoradID)}
            >
              <BiTargetLock size={20} />
            </button>
          </Tooltip>
        ) : (
          <button
            className="text-gray-400 cursor-not-allowed"
            disabled
            title="TLE data not available"
          >
            <BiTargetLock size={20} />
          </button>
        );
      },
    }),
  ];

  const table = useReactTable({
    data: currentData,
    columns,
    pageCount: Math.ceil(totalItems / pageSize),
    getCoreRowModel: getCoreRowModel(),
  });

  const handleNextPage = () => {
    if (pageIndex + 1 < table.getPageCount()) {
      onPageChange(pageIndex + 1);
    }
  };

  const handlePreviousPage = () => {
    if (pageIndex > 0) {
      onPageChange(pageIndex - 1);
    }
  };

  return (
    <Card extra={"h-full w-full pb-8 px-8"}>
      <div className="flex items-center justify-between">
        {/* Search Bar */}
        <div className="relative flex h-[61px] w-[355px] items-center rounded-full bg-white px-2 py-2 shadow-xl shadow-shadow-500 dark:bg-navy-800 dark:shadow-none">
          <button
            onClick={onSearchSubmit}
            className="flex h-full items-center justify-center rounded-l-full bg-lightPrimary text-navy-700 dark:bg-navy-900 dark:text-white px-4 text-xl"
          >
            <FiSearch className="h-4 w-4 text-gray-400 dark:text-white" />
          </button>
          <input
            type="text"
            value={searchValue}
            onChange={(e) => onSearchChange(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                onSearchSubmit();
              }
            }}
            placeholder="Search..."
            className="block h-full w-full rounded-r-full bg-lightPrimary text-sm font-medium text-navy-700 outline-none placeholder:!text-gray-400 dark:bg-navy-900 dark:text-white dark:placeholder:!text-white"
          />
        </div>
        <p className="text-sm text-gray-600 dark:text-white">
          Total Items: {totalItems}
        </p>
      </div>

      <div className="mt-8 overflow-x-scroll xl:overflow-x-hidden">
        <table className="w-full">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id} className="!border-px !border-gray-400">
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    colSpan={header.colSpan}
                    className="cursor-pointer border-b border-gray-200 pb-2 pr-4 text-start dark:border-white/30"
                  >
                    <div className="text-xs text-gray-600 dark:text-gray-200">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                    </div>
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {table.getRowModel().rows.map((row) => (
                <tr key={row.id}>
                {row.getVisibleCells().map((cell, index) => (
                    <td
                    key={`cell-${row.id}-${cell.column.id}-${index}`} // Unique key for each cell
                    className="min-w-[150px] border-white/0 py-3 pr-4"
                    >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                ))}
                </tr>
            ))}
            </tbody>
        </table>
      </div>

      {/* Pagination Controls */}
      <div className="mt-4 flex justify-between items-center">
        <button
          onClick={handlePreviousPage}
          disabled={pageIndex === 0}
          className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
        >
          <MdChevronLeft />
        </button>

        <p className="text-sm text-gray-600 dark:text-white">
          Page {pageIndex + 1} of {Math.ceil(totalItems / pageSize)}
        </p>

        <button
          onClick={handleNextPage}
          disabled={pageIndex + 1 >= Math.ceil(totalItems / pageSize)}
          className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
        >
          <MdChevronRight />
        </button>
      </div>
    </Card>
  );
}
