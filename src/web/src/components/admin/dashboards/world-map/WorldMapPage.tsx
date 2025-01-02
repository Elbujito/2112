import React, { useState } from "react";
import MapTileCardWithData from "./MapTileCardWithData"; // Map component with data
import TileMappingTableWithData from "./MappingTableWithData"; // Table component with data
import { Box, Grid, GridItem } from "@chakra-ui/react";

const WorldMapPage: React.FC = () => {
    const [selectedMappingID, setSelectedMappingID] = useState<string | null>(null);

    const handleMappingSelection = (mappingID: string) => {
        setSelectedMappingID(mappingID);
        console.log("Selected Mapping ID:", mappingID);
    };

    return (
        <Box p={4} w="100%" h="100%">
            <Grid
                templateColumns={{ base: "1fr", lg: "repeat(2, 1fr)" }} // 1 column for small screens, 2 columns for large screens
                gap={4}
                w="100%"
                h="100%"
            >
                <GridItem w="100%" h="100%" minHeight="50vh">
                    <Box w="100%" h="100%">
                        <MapTileCardWithData selectedMappingID={selectedMappingID} />
                    </Box>
                </GridItem>

                <GridItem w="100%" h="100%">
                    <Box w="100%" h="100%">
                        <TileMappingTableWithData onSelectMappingID={handleMappingSelection} />
                    </Box>
                </GridItem>
            </Grid>
        </Box>

    );
};

export default WorldMapPage;
