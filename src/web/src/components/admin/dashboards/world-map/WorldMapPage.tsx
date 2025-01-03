import React, { useState } from "react";
import MapTileView from "./MapTileView";
import MappingTableView from "./MappingTableView";
import TileTableView from "./TileTableView";
import SatelliteTableView from "./SatelliteTableView";
import SearchBar from "components/search/SearchBar";
import { Box, Grid, GridItem } from "@chakra-ui/react";
import { OrbitDataItem } from "types/satellites";

const WorldMapPage: React.FC = () => {
    const [selectedTileIDs, setSelectedTileIDs] = useState<string[]>([]);
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [appliedSearchQuery, setAppliedSearchQuery] = useState<string>("");
    const [satellitePositionData, setSatellitePositionData] = useState<Record<string, OrbitDataItem[]> | null>(null); // State to hold position data

    const handleTileSelection = (tileIDs: string[]) => {
        setSelectedTileIDs(tileIDs);
    };

    const handleSatelliteSelection = (noradID: string, positionData: Record<string, OrbitDataItem[]>) => {
        console.log("Selected Satellite NORAD ID:", noradID);
        setSatellitePositionData(positionData);
    };

    const handleSearchChange = (value: string) => {
        setSearchQuery(value);
    };

    const handleSearchSubmit = () => {
        setAppliedSearchQuery(searchQuery.toLowerCase());
        console.log("Search submitted with query:", searchQuery.toLowerCase());
    };

    return (
        <Box p={4} w="100%" h="100%">
            <Box mb={4}>
                <SearchBar
                    searchValue={searchQuery}
                    onSearchChange={handleSearchChange}
                    onSearchSubmit={handleSearchSubmit}
                />
            </Box>

            <Grid
                templateColumns={{ base: "1fr", lg: "repeat(2, 1fr)" }}
                gap={4}
                w="100%"
                h="100%"
                alignItems="stretch"
            >
                <GridItem w="100%" h="100%" minHeight="50vh" display="flex">
                    <Box flex="1" h="100%">
                        <MapTileView
                            selectedTileIDs={selectedTileIDs}
                            satellitePositionData={satellitePositionData} // Pass position data to the map
                        />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%" minHeight="50vh" maxHeight="50vh" display="flex" gap={4}>
                    <Box flex="1" h="100%" overflow="hidden">
                        <MappingTableView
                            searchQuery={appliedSearchQuery}
                            onSelectTileID={(tileID) => handleTileSelection([tileID])}
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <TileTableView
                            searchQuery={appliedSearchQuery}
                            onSelectTile={(tileID) => handleTileSelection([tileID])}
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <SatelliteTableView
                            searchQuery={appliedSearchQuery}
                            onSelectSatelliteID={(noradID) => console.log("Satellite Selected:", noradID)}
                            onTilesSelected={(tileIDs) => handleTileSelection(tileIDs)}
                            onPropagateSatellite={null}
                            onTargetSatellite={handleSatelliteSelection} // Pass the handler to SatelliteTableView
                        />
                    </Box>
                </GridItem>
            </Grid>
        </Box>
    );
};

export default WorldMapPage;
