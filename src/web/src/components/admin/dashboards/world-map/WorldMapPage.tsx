import React, { useState } from "react";
import MapTileView from "./MapTileView"
import MappingTableView from "./MappingTableView";
import { Box, Grid, GridItem } from "@chakra-ui/react";

const WorldMapPage: React.FC = () => {
    const [selectedTileID, setSelectedTileID] = useState<string | null>(null);

    const handleMappingSelection = (tileID: string) => {
        setSelectedTileID(tileID);
        console.log("Selected Tile ID:", tileID);
    };

    return (
        <Box p={4} w="100%" h="100%">
            <Grid
                templateColumns={{ base: "1fr", lg: "repeat(2, 1fr)" }}
                gap={4}
                w="100%"
                h="100%"
            >
                <GridItem w="100%" h="100%" minHeight="50vh">
                    <MapTileView selectedTileID={selectedTileID} />
                </GridItem>

                <GridItem w="100%" h="100%">
                    <MappingTableView onSelectTileID={handleMappingSelection} />
                </GridItem>
            </Grid>
        </Box>

    );
};

export default WorldMapPage;
