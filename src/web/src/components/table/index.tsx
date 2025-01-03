import React, { useState } from "react";
import Card from "components/card";
import { MdChevronRight, MdChevronLeft, MdMoreVert } from "react-icons/md";
import {
    useReactTable,
    getCoreRowModel,
    getSortedRowModel,
    flexRender,
    ColumnDef,
} from "@tanstack/react-table";
import { Menu, MenuItem, MenuButton } from "@chakra-ui/react";
import { CustomScrollbar } from "components/scrollbar/CustomScrollbar";

function TableContainer<T>({
    columns,
    data,
    pageIndex,
    pageSize,
    totalItems,
    onPageChange,
    getRowId,
    actions,
    onRowClick,
}: {
    columns: ColumnDef<T>[];
    data: T[];
    pageIndex: number;
    pageSize: number;
    totalItems: number;
    onPageChange: (pageIndex: number) => void;
    getRowId: (row: T) => string;
    actions?: (row: T) => { label: string; onClick: () => void }[] | null;
    onRowClick?: (row: T) => void;
}) {
    const [selectedRowId, setSelectedRowId] = useState<string | null>(null);

    const table = useReactTable({
        data,
        columns,
        pageCount: Math.ceil(totalItems / pageSize),
        getCoreRowModel: getCoreRowModel(),
        getSortedRowModel: getSortedRowModel(),
        columnResizeMode: "onChange",
    });

    const headerGroups = table.getHeaderGroups() || [];
    const rows = table.getRowModel().rows || [];

    return (
        <Card extra="h-full w-full pb-2 px-2">
            <CustomScrollbar style={{ height: "55vh", overflowX: "auto" }}>
                <table
                    className="w-full border-collapse table-fixed"
                    style={{ width: table.getTotalSize() }}
                >
                    {/* Table Headers */}
                    <thead className="sticky top-0 bg-navy-700 z-10 text-white">
                        {headerGroups.map((headerGroup) => (
                            <tr key={headerGroup.id}>
                                {headerGroup.headers.map((header) => (
                                    <th
                                        key={header.id}
                                        className="py-2 px-3 text-left text-xs font-bold relative group"
                                        style={{
                                            width: header.getSize(),
                                            minWidth: header.column.columnDef.minSize,
                                            maxWidth: header.column.columnDef.maxSize,
                                        }}
                                    >
                                        {flexRender(
                                            header.column.columnDef.header,
                                            header.getContext()
                                        )}
                                    </th>
                                ))}
                                {actions && (
                                    <th className="py-2 px-3 text-left text-xs font-bold">
                                        Actions
                                    </th>
                                )}
                            </tr>
                        ))}
                    </thead>

                    {/* Table Rows */}
                    <tbody>
                        {rows.length > 0 ? (
                            rows.map((row) => (
                                <tr
                                    key={getRowId(row.original)}
                                    onClick={() => {
                                        setSelectedRowId(getRowId(row.original));
                                        onRowClick?.(row.original);
                                    }}
                                    className={`cursor-pointer ${selectedRowId === getRowId(row.original)
                                            ? "bg-blue-600 text-white"
                                            : "hover:bg-gray-200"
                                        }`}
                                >
                                    {row.getVisibleCells().map((cell) => (
                                        <td
                                            key={cell.id}
                                            className="py-2 px-3 text-xs truncate"
                                            title={String(cell.getValue())}
                                        >
                                            {flexRender(
                                                cell.column.columnDef.cell,
                                                cell.getContext()
                                            )}
                                        </td>
                                    ))}
                                    {actions && (
                                        <td className="py-2 px-3 text-xs">
                                            <Menu>
                                                <MenuButton>
                                                    <MdMoreVert className="cursor-pointer" />
                                                </MenuButton>
                                                {(actions(row.original) || []).map(
                                                    (action, index) => (
                                                        <MenuItem
                                                            key={index}
                                                            onClick={action.onClick}
                                                        >
                                                            {action.label}
                                                        </MenuItem>
                                                    )
                                                )}
                                            </Menu>
                                        </td>
                                    )}
                                </tr>
                            ))
                        ) : (
                            <tr>
                                <td
                                    colSpan={columns.length + 1}
                                    className="text-center py-4 text-xs"
                                >
                                    No data available
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </CustomScrollbar>

            {/* Sticky Pagination */}
            <div className="sticky bottom-0 bg-navy-700 z-10 text-white p-2 flex justify-between items-center">
                <p className="text-xs">Total Results: {totalItems}</p>
                <div className="flex items-center">
                    <button
                        onClick={() => onPageChange(Math.max(0, pageIndex - 1))}
                        disabled={pageIndex === 0}
                        className="px-2 py-1 text-xs bg-blue-500 hover:bg-blue-600 text-white rounded-lg mr-2 disabled:bg-gray-400 disabled:text-gray-200"
                    >
                        <MdChevronLeft />
                    </button>
                    <p className="text-xs">
                        Page {pageIndex + 1} of {Math.ceil(totalItems / pageSize)}
                    </p>
                    <button
                        onClick={() => onPageChange(Math.min(pageIndex + 1, Math.ceil(totalItems / pageSize) - 1))}
                        disabled={pageIndex + 1 >= Math.ceil(totalItems / pageSize)}
                        className="px-2 py-1 text-xs bg-blue-500 hover:bg-blue-600 text-white rounded-lg ml-2 disabled:bg-gray-400 disabled:text-gray-200"
                    >
                        <MdChevronRight />
                    </button>
                </div>
            </div>
        </Card>
    );
}

export default TableContainer;
