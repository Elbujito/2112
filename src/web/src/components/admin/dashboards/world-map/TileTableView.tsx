import React, { useEffect, useState } from "react";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import useTileServiceStore from "services/tileService";
import GenericTableComponent from "components/table";
import { Tile } from "types/tiles";

interface TileTableViewProps {
    onSelectTile: (tile: string) => void;
}

export default function TileTableView({ onSelectTile }: TileTableViewProps) {
    const { tiles, loading, error, fetchTilesForLocation } = useTileServiceStore();

    const [search, setSearch] = useState<string>("");
    const [pageIndex, setPageIndex] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(20);
    const [paginatedTiles, setPaginatedTiles] = useState<Tile[]>([]);

    const defaultLocation = { latitude: 0, longitude: 0 };

    useEffect(() => {
        fetchTilesForLocation(defaultLocation);
    }, []);

    useEffect(() => {
        const start = pageIndex * pageSize;
        const end = start + pageSize;
        setPaginatedTiles(tiles.slice(start, end));
    }, [tiles, pageIndex, pageSize]);

    const handleSearchChange = (value: string) => {
        setSearch(value);
    };

    const handleSearchSubmit = () => {
        setPageIndex(0);
        fetchTilesForLocation(defaultLocation);
    };

    const handleOnPaginationChange = (index: number) => {
        console.log(index)
        setPageIndex(index);
    };

    const handleTileSelection = (tile: Tile) => {
        onSelectTile(tile.ID);
    };

    const columns = [
        {
            accessorKey: "ID",
            header: "Tile ID",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "CenterLat",
            header: "Center Lat",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "CenterLon",
            header: "Center Lon",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "ZoomLevel",
            header: "Zoom Level",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "NbFaces",
            header: "Nb Faces",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Radius",
            header: "Radius",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
    ];

    if (loading) {
        return (
            <Center className="grid h-full w-full place-items-center">
                <Spinner size="xl" color="blue.500" />
            </Center>
        );
    }

    if (error) {
        return (
            <Center className="grid h-full w-full place-items-center">
                <Box className="text-center">
                    <Text className="text-lg font-bold text-red-500">{error}</Text>
                </Box>
            </Center>
        );
    }

    return (
        <Box className="grid w-full gap-4 rounded-lg shadow-md">
            <GenericTableComponent
                getRowId={(row: Tile) => row.ID}
                columns={columns}
                data={paginatedTiles}
                totalItems={tiles.length}
                pageSize={pageSize}
                pageIndex={pageIndex}
                onPageChange={handleOnPaginationChange}
                searchValue={search}
                onSearchChange={handleSearchChange}
                onSearchSubmit={handleSearchSubmit}
                onRowClick={handleTileSelection}
            />
        </Box>
    );
}
