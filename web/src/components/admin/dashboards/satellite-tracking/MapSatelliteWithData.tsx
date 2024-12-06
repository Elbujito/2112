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
  userLocation?: { latitude: number; longitude: number }; // Optional user location
}

const MapSatelliteWithData: React.FC<MapSatelliteWithDataProps> = ({
  noradID,
  userLocation,
}) => {
  const [orbitData, setOrbitData] = useState<OrbitDataItem[]>([]);
  const [fetchingOrbit, setFetchingOrbit] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch orbit data for the provided NORAD ID
  useEffect(() => {
    const fetchOrbitData = async () => {
      if (!noradID) {
        return;
      }

      setFetchingOrbit(true);
      setError(null); // Reset error before fetching
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
        setError("Failed to fetch orbit data. Please try again later.");
      } finally {
        setFetchingOrbit(false);
      }
    };

    fetchOrbitData();
  }, [noradID]);

  return (
    <Card extra={"w-full h-full bg-white px-3 py-[18px]"}>
      <Box
        className="grid h-[60vh] grid-cols-1 grid-rows-1 rounded-md"
      >
        {/* MapSatellite Component */}
        <MapSatellite orbitData={orbitData} noradID={noradID} userLocation={userLocation} />

        {/* Loading Indicator */}
        {fetchingOrbit && (
          <Center className="absolute inset-0 bg-black/70 z-10 flex flex-col">
            <Spinner thickness="4px" speed="0.65s" color="blue.500" size="xl" />
            <Text mt={4} color="white">
              Fetching Orbit Data...
            </Text>
          </Center>
        )}

        {/* Error Message */}
        {error && (
          <Center className="absolute inset-0 bg-red-700 text-white z-10">
            <Text>{error}</Text>
          </Center>
        )}
      </Box>
    </Card>
  );
};

export default MapSatelliteWithData;
