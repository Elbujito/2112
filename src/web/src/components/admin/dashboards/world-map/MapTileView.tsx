import React, { useEffect } from "react";
import { Spinner, Center, Text } from "@chakra-ui/react";
import useTileServiceStore from "services/tileService"; // Import your tile service store
import MapTileCard from "./MapTileCard";

interface MapTileViewProps {
  selectedTileID?: string;
}

export default function MapTileView({
  selectedTileID,
}: MapTileViewProps) {
  const {
    tiles,
    loading,
    error,
    fetchTilesForLocation,
  } = useTileServiceStore();

  useEffect(() => {
    fetchTilesForLocation({ latitude: 0, longitude: 0 });
  }, [fetchTilesForLocation]);

  if (loading) {
    return (
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
          Fetching Tile Data...
        </Text>
      </Center>
    );
  }

  if (error) {
    return (
      <Center
        position="absolute"
        top="0"
        left="0"
        right="0"
        bottom="0"
        bg="red.700"
        color="white"
      >
        <Text>{error}</Text>
      </Center>
    );
  }

  return (
    <MapTileCard
      tiles={tiles}
      darkmode={
        document.body.classList.contains("dark")
          ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
          : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
      }
      onLocationChange={fetchTilesForLocation}
      selectedTileID={selectedTileID}
    />
  );
}
