import React, { useState } from "react";
import MapTileView from "./MapTileView";
import MappingTableView from "./MappingTableView";
import TileTableView from "./TileTableView";
import SatelliteTableView from "./SatelliteTableView";
import SearchBar from "components/search/SearchBar";
import { Box, Grid, GridItem } from "@chakra-ui/react";

const WorldMapPage: React.FC = () => {
    const [selectedTileID, setSelectedTileID] = useState<string | null>(null);
    const [searchQuery, setSearchQuery] = useState<string>(""); // For input field
    const [appliedSearchQuery, setAppliedSearchQuery] = useState<string>(""); // For filtering

    const handleTileSelection = (tileID: string) => {
        setSelectedTileID(tileID);
        console.log("Selected Tile ID:", tileID);
    };

    const handleSatelliteSelection = (noradID: string) => {
        console.log("Selected satellite ID:", noradID);
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
                        <MapTileView selectedTileID={selectedTileID} />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%" minHeight="50vh" maxHeight="50vh" display="flex">
                    <Box flex="1" h="100%" overflow="hidden">
                        <MappingTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectTileID={handleTileSelection}
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <TileTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectTile={handleTileSelection}
                        />
                    </Box>

                    <Box flex="1" h="100%" overflow="hidden">
                        <SatelliteTableView
                            searchQuery={appliedSearchQuery} // Use applied query
                            onSelectSatelliteID={handleSatelliteSelection}
                        />
                    </Box>
                </GridItem>
            </Grid>
        </Box>
    );
};

export default WorldMapPage;
