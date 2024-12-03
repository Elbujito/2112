import React, { useEffect, useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { FiSearch } from "react-icons/fi";
import { BiTargetLock } from "react-icons/bi";

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
  Action: string | null;
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
  const [selectedRowId, setSelectedRowId] = useState<string | null>(null);

  useEffect(() => {
    setCurrentData(tableData);
  }, [tableData, onPageChange, pageIndex]);

  const columns = [
    columnHelper.accessor("Action", {
      header: () => null,
      cell: (info) => {
        const row = info.row.original;
        return row.TleUpdatedAt ? (
          <button
            className="text-blue-500 hover:text-blue-700"
            onClick={() => onRowClick(row.NoradID)}
          >
            <BiTargetLock size={20} />
          </button>
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
    columnHelper.accessor("Name", {
      header: () => <p className="text-sm font-bold text-gray-600">NAME</p>,
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("NoradID", {
      header: () => <p className="text-sm font-bold text-gray-600">NORAD ID</p>,
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("Owner", {
      header: () => <p className="text-sm font-bold text-gray-600">OWNER</p>,
      cell: (info) => <p className="text-sm">{info.getValue()}</p>,
    }),
    columnHelper.accessor("LaunchDate", {
      header: () => (
        <p className="text-sm font-bold text-gray-600">LAUNCH DATE</p>
      ),
      cell: (info) => (
        <p className="text-sm">
          {info.getValue()
            ? new Date(info.getValue() as string).toLocaleDateString()
            : "N/A"}
        </p>
      ),
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
    <Card extra="h-full w-full pb-8 px-8">
      <div className="sticky top-0 z-10">
        <div className="flex items-center justify-between">
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
                if (e.key === "Enter") onSearchSubmit();
              }}
              placeholder="Search..."
              className="block h-full w-ful rounded-r-full bg-lightPrimary text-sm font-medium text-navy-700 outline-none placeholder:!text-gray-400 dark:bg-navy-900 dark:text-white dark:placeholder:!text-white"
            />
          </div>
          <p className="text-sm">Total Satellites: {totalItems}</p>
        </div>
      </div>

      <div className="mt-4 overflow-x-auto">
        <table className="w-full">
          <thead>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    colSpan={header.colSpan}
                    className="py-2 text-left text-xs font-bold text-gray-600"
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
          <tbody>
            {table.getRowModel().rows.map((row) => (
              <tr
                key={row.id}
                onClick={() => setSelectedRowId(row.id)}
                onDoubleClick={() => onRowClick(row.original.NoradID)}
                className={`cursor-pointer ${
                  selectedRowId === row.id ? "bg-brand-900 text-white dark:bg-brand-400 dark:text-white" : ""
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
        </table>
      </div>

      <div className="sticky flex justify-between">
        <button
          onClick={handlePreviousPage}
          disabled={pageIndex === 0}
          className="px-4 py-2 bg-gray-200 dark:bg-gray-700 rounded-lg"
        >
          <MdChevronLeft />
        </button>
        <p className="text-sm">
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
