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
        <Box p={4} w="100%">
            <Grid
                templateRows="repeat(2, 1fr)" // Two vertical rows
                gap={4}
                w="100%"
                h="100%"
                templateColumns={{ base: "1fr", lg: "repeat(2, 1fr)" }} // Adjust for responsiveness
            >
                <GridItem w="100%">
                    <MapTileCardWithData />
                </GridItem>
                <GridItem w="100%">
                    <TileMappingTableWithData onSelectMappingID={handleMappingSelection} />
                </GridItem>
            </Grid>
        </Box>
    );
};

export default WorldMapPage;
