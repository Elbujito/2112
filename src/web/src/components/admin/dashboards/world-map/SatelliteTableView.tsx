import React, { useState, useEffect, useRef } from "react";
import { Spinner, Box, Text, Center, AlertDialog, AlertDialogOverlay, AlertDialogContent, AlertDialogHeader, AlertDialogBody, AlertDialogFooter, Button } from "@chakra-ui/react";
import useSatelliteServiceStore from "services/satelliteService"; // Satellite store
import useTileServiceStore from "services/tileService"; // Tile store
import GenericTableComponent from "components/table";
import { OrbitDataItem, SatelliteInfo } from "types/satellites";
import { BiTargetLock } from "react-icons/bi";

interface SatelliteTableViewProps {
    onSelectSatelliteID: (satelliteID: string) => void;
    searchQuery: string;
    onTilesSelected: (tileIDs: string[]) => void;
    onTargetSatellite: (noradID: string, positionData: Record<string, OrbitDataItem[]>) => void; // Callback for targeting satellite with position data
    onPropagateSatellite: (noradID: string) => void; // Callback for targeting satellite
}

export default function SatelliteTableView({
    onSelectSatelliteID,
    searchQuery,
    onTilesSelected,
    onTargetSatellite,
    onPropagateSatellite,
}: SatelliteTableViewProps) {
    const { satelliteInfo, totalSatelliteInfo, loading, error, orbitData, fetchPaginatedSatelliteInfo, fetchSatellitePositions } = useSatelliteServiceStore();
    const { tileMappings } = useTileServiceStore();

    const [pageIndex, setPageIndex] = useState<number>(0);
    const [pageSize, setPageSize] = useState<number>(20);
    const [isErrorDialogOpen, setIsErrorDialogOpen] = useState(false); // Manage error dialog visibility
    const [errorMessage, setErrorMessage] = useState<string>("");

    const cancelRef = useRef<HTMLButtonElement>(null); // Reference for the cancel button

    useEffect(() => {
        fetchPaginatedSatelliteInfo(pageIndex, pageSize, searchQuery);
    }, [pageIndex, pageSize, searchQuery, fetchPaginatedSatelliteInfo]);

    const handleOnPaginationChange = (index: number) => {
        setPageIndex(index);
    };

    const handleSatelliteSelection = (satellite: SatelliteInfo) => {
        const noradId = satellite.Satellite.NoradID;

        const matchingTiles = tileMappings.filter((mapping) => mapping.NoradID === noradId);
        const matchingTileIDs = matchingTiles.map((tile) => tile.TileID);

        onSelectSatelliteID(noradId);
        onTilesSelected(matchingTileIDs);
    };

    const handleTargetSatellite = async (noradID: string) => {
        onTargetSatellite(noradID, { [noradID]: [] });
        const startTime = new Date(Date.now()).toISOString();
        const endTime = new Date(Date.now() + 60 * 60 * 1000 * 24).toISOString(); // 24 hours ahead

        try {
            const positionData = await fetchSatellitePositions(noradID, startTime, endTime);
            onTargetSatellite(noradID, orbitData); // Pass the position data to the parent callback
        } catch (error) {
            console.error("Error fetching satellite positions:", error);
            setErrorMessage(error.message || "Failed to fetch satellite positions.");
            setIsErrorDialogOpen(true); // Show error dialog
            onTargetSatellite(noradID, { [noradID]: [] });
        }
    };

    const handlePropagateSatellite = (noradID: string) => {
        // Invoke the callback to propagate satellite
        onPropagateSatellite(noradID);
    };

    const columns = [
        {
            accessorKey: "",
            header: "Actions",
            cell: (info: any) => <p className="text-sm">{info.getValue()}</p>,
        },
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
    ];

    if (loading) {
        return (
            <Center className="grid h-full w-full place-items-center">
                <Spinner size="xl" color="blue.500" />
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
                onRowClick={handleSatelliteSelection}
                actions={(row: SatelliteInfo) => [
                    {
                        label: "Target",
                        onClick: () => handleTargetSatellite(row.Satellite.NoradID),
                        icon: <BiTargetLock />,
                    },
                    {
                        label: "Propagate",
                        onClick: () => handlePropagateSatellite(row.Satellite.NoradID),
                    },
                ]}
            />

            {/* Error Dialog (Popup) */}
            <AlertDialog
                isOpen={isErrorDialogOpen}
                onClose={() => setIsErrorDialogOpen(false)}
                leastDestructiveRef={cancelRef} // Add the cancelRef here
            >
                <AlertDialogOverlay>
                    <AlertDialogContent>
                        <AlertDialogHeader>Error</AlertDialogHeader>
                        <AlertDialogBody>
                            {errorMessage || "An unexpected error occurred."}
                        </AlertDialogBody>
                        <AlertDialogFooter>
                            <Button ref={cancelRef} onClick={() => setIsErrorDialogOpen(false)}>
                                Close
                            </Button>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialogOverlay>
            </AlertDialog>
        </Box>
    );
}
