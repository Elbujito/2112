import React, { useState } from "react";
import MapTileView from "./MapTileView";
import MappingTableView from "./MappingTableView";
import TileTableView from "./TileTableView"; // Import the TileTableView component
import { Box, Grid, GridItem, Divider, Flex } from "@chakra-ui/react";

const WorldMapPage: React.FC = () => {
    const [selectedTileID, setSelectedTileID] = useState<string | null>(null);

    const handleTileSelection = (tileID: string) => {
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
                alignItems="stretch"
            >
                <GridItem w="100%" h="100%" minHeight="50vh" display="flex">
                    <Box flex="1" h="100%">
                        <MapTileView selectedTileID={selectedTileID} />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%" minHeight="50vh" maxHeight="50vh"  display="flex">
                    <Flex w="100%" h="100%" gap={4}>
                        <Box flex="1" h="100%" overflow="hidden">
                            <MappingTableView onSelectTileID={handleTileSelection} />
                        </Box>

                        <Box flex="1" h="100%" overflow="hidden">
                            <TileTableView onSelectTile={handleTileSelection} />
                        </Box>
                    </Flex>
                </GridItem>
            </Grid>
        </Box>
    );
};

export default WorldMapPage;
