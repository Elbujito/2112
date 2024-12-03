import React, { useEffect, useState } from "react";
import axios from "axios";
import { Spinner, Box, Text, Center } from "@chakra-ui/react";
import Card from "components/card";
import MapSatellite from "./MapSatellite"; // Import the MapSatellite component

type OrbitDataItem = {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
};

interface MapSatelliteWithDataProps {
  noradID: string; // NORAD ID passed as a prop
}

const MapSatelliteWithData: React.FC<MapSatelliteWithDataProps> = ({
  noradID,
}) => {
  const [orbitData, setOrbitData] = useState<OrbitDataItem[]>([]);
  const [fetchingOrbit, setFetchingOrbit] = useState<boolean>(false);

  // Fetch orbit data for the provided NORAD ID
  useEffect(() => {
    const fetchOrbitData = async () => {
      if (!noradID) {
        console.warn("NORAD ID is required to fetch orbit data.");
        return;
      }

      setFetchingOrbit(true);
      try {
        const response = await axios.get("http://localhost:8081/satellites/orbit", {
          params: { noradID },
          headers: { Accept: "application/json" },
        });

        const orbitDataMapped = response.data.map((data: any) => ({
          latitude: data.Latitude,
          longitude: data.Longitude,
          altitude: data.Altitude,
          time: data.Time,
        }));
        setOrbitData(orbitDataMapped);
      } catch (error) {
        console.error("Error fetching orbit data:", error);
      } finally {
        setFetchingOrbit(false);
      }
    };

    fetchOrbitData();
  }, [noradID]);

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"}>
      <Box
        position="relative"
        w="100%"
        h="60vh"
        overflow="hidden"
        bg="gray.50"
        borderRadius="md"
      >
        {/* Use MapSatellite to Display Orbit Data */}
        <MapSatellite orbitData={orbitData} noradID={noradID} />

        {/* Loading Indicator */}
        {fetchingOrbit && (
          <Center
            position="absolute"
            top="0"
            left="0"
            right="0"
            bottom="0"
            bg="blackAlpha.700"
          >
            <Spinner thickness="4px" speed="0.65s" color="blue.500" size="xl" />
            <Text mt={4} color="white">
              Fetching Orbit Data...
            </Text>
          </Center>
        )}
      </Box>
    </Card>
  );
};

export default MapSatelliteWithData;
