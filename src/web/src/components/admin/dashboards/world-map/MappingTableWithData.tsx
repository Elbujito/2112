import React, { useState, useEffect } from "react";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import MappingTableComponent from "./MappingTable"; // Your table component
import axios from "axios";

interface TileMappingTableWithDataProps {
    onSelectTileID: (tileID: string) => void; // Callback for selected mapping ID
}

export type TileSatelliteMapping = {
    MappingID: string;
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

export default function TileMappingTableWithData({
    onSelectTileID,
}: TileMappingTableWithDataProps) {
    const [mappings, setMappings] = useState([]); // Tile mapping data
    const [loading, setLoading] = useState(true); // Loading state
    const [error, setError] = useState<string | null>(null); // Error state
    const [search, setSearch] = useState<string>(""); // Search term
    const [pageIndex, setPageIndex] = useState<number>(0); // Current page index
    const [pageSize, setPageSize] = useState<number>(20); // Rows per page
    const [totalItems, setTotalItems] = useState<number>(0); // Total items
    const [submit, setSubmit] = useState<boolean>(true); // Trigger fetch

    useEffect(() => {
        const fetchMappings = async () => {
            if (!submit) return;

            setLoading(true);
            setError(null);

            try {
                const response = await axios.get("http://localhost:8081/tiles/mappings", {
                    params: { page: pageIndex + 1, pageSize, search },
                    headers: { Accept: "application/json" },
                });

                setMappings(response.data.mappings || []);
                setTotalItems(response.data.totalRecords || 0);
            } catch (err) {
                console.error("Error fetching tile mappings:", err);
                setError("Failed to load tile mapping data.");
            } finally {
                setLoading(false);
                setSubmit(false);
            }
        };

        fetchMappings();
    }, [pageIndex, pageSize, submit]);

    const handleSearchChange = (value: string) => {
        setSearch(value);
    };

    const handleSearchSubmit = () => {
        setPageIndex(0);
        setSubmit(true);
    };

    const handleOnPaginationChange = (index: number) => {
        setPageIndex(index);
        setSubmit(true);
    };

    const handleMappingSelection = (mapping: TileSatelliteMapping) => {
        onSelectTileID(mapping.TileID);
    };

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
            <MappingTableComponent
                tableData={mappings}
                pageIndex={pageIndex}
                pageSize={pageSize}
                onPageChange={handleOnPaginationChange}
                searchValue={search}
                onSearchChange={handleSearchChange}
                onSearchSubmit={handleSearchSubmit}
                totalItems={totalItems}
                onRowClick={handleMappingSelection}
            />
        </Box>
    );
}
