import React, { useState } from "react";
import MapTileCardWithData from "./MapTileCardWithData"
import TileMappingTableWithData from "./MappingTableWithData";
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
                    <Box w="100%" h="100%">
                        <MapTileCardWithData selectedTileID={selectedTileID} />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%">
                    <Box w="100%" h="100%">
                        <TileMappingTableWithData onSelectTileID={handleMappingSelection} />
                    </Box>
                </GridItem>
            </Grid>
        </Box>

    );
};

export default WorldMapPage;
