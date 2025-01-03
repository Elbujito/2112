import React, { useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { FiSearch } from "react-icons/fi";
import { CustomScrollbar } from "components/scrollbar/CustomScrollbar";
import {
    useReactTable,
    getCoreRowModel,
    flexRender,
    ColumnDef,
} from "@tanstack/react-table";

function SearchBarWithPagination({
    searchValue,
    onSearchChange,
    onSearchSubmit,
    pageIndex,
    pageCount,
    onPageChange,
}: {
    searchValue: string;
    onSearchChange: (value: string) => void;
    onSearchSubmit: () => void;
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
        <div className="flex justify-between items-center mb-4">
            <div className="flex flex-1 max-w-lg">
                <div className="relative flex w-full items-center rounded-full bg-white px-2 py-2 shadow-xl dark:bg-navy-800">
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
            </div>
            <div className="flex items-center ml-4">
                <button
                    onClick={handlePreviousPage}
                    disabled={pageIndex === 0}
                    className="px-4 py-2 bg-gray-200 dark:bg-navy-800 rounded-lg mr-2"
                >
                    <MdChevronLeft />
                </button>
                <p className="text-sm">
                    Page {pageIndex + 1} of {pageCount}
                </p>
                <button
                    onClick={handleNextPage}
                    disabled={pageIndex + 1 >= pageCount}
                    className="px-4 py-2 bg-gray-200 dark:bg-navy-800 rounded-lg ml-2"
                >
                    <MdChevronRight />
                </button>
            </div>
        </div>
    );
}

function TableHeader({ headers }: { headers: any[] }) {
    return (
        <thead className="sticky top-0 bg-white dark:bg-navy-800 z-10">
            {headers.map((headerGroup) => (
                <tr key={headerGroup.id}>
                    {headerGroup.headers.map((header) => (
                        <th
                            key={header.id}
                            colSpan={header.colSpan}
                            className="py-2 text-left text-xs font-bold text-gray-600 dark:text-white"
                        >
                            {flexRender(header.column.columnDef.header, header.getContext())}
                        </th>
                    ))}
                </tr>
            ))}
        </thead>
    );
}

function TableBody<T>({
    rows,
    onRowClick,
    selectedRowId,
    getRowId,
}: {
    rows: any[];
    onRowClick?: (row: T) => void;
    selectedRowId?: string | null;
    getRowId: (row: T) => string;
}) {
    return (
        <tbody>
            {rows.map((row) => (
                <tr
                    key={getRowId(row.original)}
                    onClick={() => onRowClick?.(row.original)}
                    className={`cursor-pointer ${selectedRowId === getRowId(row.original)
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

export default function GenericTableComponent<T>({
    columns,
    data = [],
    totalItems,
    pageSize,
    pageIndex,
    onPageChange,
    searchValue,
    onSearchChange,
    onSearchSubmit,
    onRowClick,
    getRowId,
}: {
    columns: ColumnDef<T>[];
    data: T[];
    totalItems: number;
    pageSize: number;
    pageIndex: number;
    onPageChange: (pageIndex: number) => void;
    searchValue: string;
    onSearchChange: (search: string) => void;
    onSearchSubmit: () => void;
    onRowClick?: (row: T) => void;
    getRowId: (row: T) => string;
}) {
    const [selectedRowId, setSelectedRowId] = useState<string | null>(null);

    const handleRowClick = (row: T) => {
        const rowId = getRowId(row);
        setSelectedRowId(rowId);
        onRowClick?.(row);
    };

    const table = useReactTable({
        data,
        columns,
        pageCount: Math.ceil(totalItems / pageSize),
        getCoreRowModel: getCoreRowModel(),
    });

    const headerGroups = table?.getHeaderGroups() || [];
    const rows = table?.getRowModel()?.rows || [];

    return (
        <Card extra="h-full w-full pb-4 px-4">
            <SearchBarWithPagination
                searchValue={searchValue}
                onSearchChange={onSearchChange}
                onSearchSubmit={onSearchSubmit}
                pageIndex={pageIndex}
                pageCount={Math.ceil(totalItems / pageSize)}
                onPageChange={onPageChange}
            />
            <CustomScrollbar style={{ minHeight: "50vh", height: "100%" }}>
                <table className="w-full">
                    {headerGroups.length > 0 && <TableHeader headers={headerGroups} />}
                    {rows.length > 0 ? (
                        <TableBody
                            rows={rows}
                            onRowClick={handleRowClick}
                            selectedRowId={selectedRowId}
                            getRowId={getRowId}
                        />
                    ) : (
                        <tbody>
                            <tr>
                                <td colSpan={columns.length} className="text-center py-4">
                                    No data available
                                </td>
                            </tr>
                        </tbody>
                    )}
                </table>
            </CustomScrollbar>
        </Card>
    );
}
