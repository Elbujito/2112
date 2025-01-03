import React, { useEffect, useState } from "react";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import useSatelliteServiceStore from "services/satelliteService"; // New store for SatelliteInfo
import GenericTableComponent from "components/table";
import { SatelliteInfo } from "types/satellites";

interface SatelliteTableViewProps {
    onSelectSatelliteID: (satelliteID: string) => void;
    searchQuery: string;
}

export default function SatelliteTableView({
    onSelectSatelliteID,
    searchQuery,
}: SatelliteTableViewProps) {
    const { satelliteInfo, totalSatelliteInfo, loading, error, fetchPaginatedSatelliteInfo } = useSatelliteServiceStore();

    const [pageIndex, setPageIndex] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(20);

    useEffect(() => {
        fetchPaginatedSatelliteInfo(pageIndex, pageSize, searchQuery);
    }, [pageIndex, pageSize, searchQuery]);


    const handleSearchSubmit = () => {
        setPageIndex(0);
        fetchPaginatedSatelliteInfo(0, pageSize, searchQuery);
    };

    const handleOnPaginationChange = (index: number) => {
        setPageIndex(index);
    };

    const handleSatelliteSelection = (satellite: SatelliteInfo) => {
        onSelectSatelliteID(satellite.Satellite.NoradID);
    };

    const columns = [
        {
            accessorKey: "Satellite.NoradID",
            header: "NORAD ID",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.Name",
            header: "Name",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.Owner",
            header: "Owner",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
        {
            accessorKey: "Satellite.LaunchDate",
            header: "Launch Date",
            cell: (info: any) =>
                info.getValue() ? (
                    <p className="text-sm">{new Date(info.getValue()).toLocaleDateString()}</p>
                ) : (
                    <p className="text-sm">N/A</p>
                ),
        },
        {
            accessorKey: "Satellite.Apogee",
            header: "Apogee (km)",
            cell: (info: any) => <p className="text-sm">{info.getValue() ?? "N/A"}</p>,
        },
        {
            accessorKey: "Satellite.Perigee",
            header: "Perigee (km)",
            cell: (info: any) => <p className="text-sm">{info.getValue() ?? "N/A"}</p>,
        },
        {
            accessorKey: "TLEs[0].Epoch",
            header: "Latest TLE Update",
            cell: (info: any) =>
                info.getValue() ? (
                    <p className="text-sm">{info.getValue()}</p>
                ) : (
                    <p className="text-sm">N/A</p>
                ),
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
                getRowId={(row: SatelliteInfo) => row.Satellite.NoradID}
                columns={columns}
                data={satelliteInfo}
                totalItems={totalSatelliteInfo}
                pageSize={pageSize}
                pageIndex={pageIndex}
                onPageChange={handleOnPaginationChange}
            />
        </Box>
    );
}
