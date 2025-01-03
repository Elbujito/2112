import React, { useState } from "react";
import MapTileView from "./MapTileView";
import MappingTableView from "./MappingTableView";
import TileTableView from "./TileTableView";
import SatelliteTableView from "./SatelliteTableView";
import SearchBar from "components/search/SearchBar";
import { Box, Grid, GridItem } from "@chakra-ui/react";

const WorldMapPage: React.FC = () => {
    const [selectedTileIDs, setSelectedTileIDs] = useState<string[]>([]); // Array for multiple tile IDs
    const [searchQuery, setSearchQuery] = useState<string>(""); // For input field
    const [appliedSearchQuery, setAppliedSearchQuery] = useState<string>(""); // For filtering

    const handleTileSelection = (tileIDs: string[]) => {
        setSelectedTileIDs(tileIDs);
        console.log("Selected Tile IDs:", tileIDs);
    };

    const handleSatelliteSelection = (noradID: string, tileIDs: string[]) => {
        console.log("Selected Satellite NORAD ID:", noradID);
        handleTileSelection(tileIDs); // Update selected tiles based on satellite selection
    };

    const handleSearchChange = (value: string) => {
        setSearchQuery(value);
    };

    const handleSearchSubmit = () => {
        setAppliedSearchQuery(searchQuery.toLowerCase()); // Apply the search
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
                        <MapTileView selectedTileIDs={selectedTileIDs} />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%" minHeight="50vh" maxHeight="50vh" display="flex" gap={4}>
                    <Box flex="1" h="100%" overflow="hidden">
                        <MappingTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectTileID={(tileID) => handleTileSelection([tileID])} // Wrap single tile ID in array
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <TileTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectTile={(tileID) => handleTileSelection([tileID])} // Wrap single tile ID in array
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <SatelliteTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectSatelliteID={(noradID) => console.log("Satellite Selected:", noradID)}
                            onTilesSelected={(tileIDs) => handleTileSelection(tileIDs)} // Update tiles from satellite
                        />
                    </Box>
                </GridItem>
            </Grid>
        </Box>
    );
};

export default WorldMapPage;
