import React, { useEffect, useState } from "react";
import axios from "axios";
import { Spinner, Center, Text } from "@chakra-ui/react";
import MapTileCard from "./MapTileCard";

interface Tile {
  Quadkey: string;
  ZoomLevel: number;
  CenterLat: number;
  CenterLon: number;
  SpatialIndex?: string;
  NbFaces: number;
  Radius: number;
  BoundariesJSON?: string;
  ID: string;
}

interface MapTileCardWithDataProps {
  selectedTileID?: string;
}

export default function MapTileCardWithData({
  selectedTileID,
}: MapTileCardWithDataProps) {
  const [tiles, setTiles] = useState<Tile[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [darkmode, setDarkmode] = useState(
    document.body.classList.contains("dark")
      ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
      : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
  );

  useEffect(() => {
    const observer = new MutationObserver(() => {
      setDarkmode(
        document.body.classList.contains("dark")
          ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
          : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
      );
    });
    observer.observe(document.body, { attributes: true });
    return () => observer.disconnect();
  }, []);

  const fetchTilesForLocation = async ({ latitude, longitude }: { latitude: number; longitude: number }) => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get("http://localhost:8081/tiles/all", { //!user region for use specific location
        headers: { Accept: "application/json" },
        params: {
          minLat: latitude - 1,
          minLon: longitude - 1,
          maxLat: latitude + 1,
          maxLon: longitude + 1,
        },
      });
      setTiles(response.data);
    } catch (err) {
      console.error("Error fetching tile data:", err);
      setError("Failed to fetch tile data. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTilesForLocation({ latitude: 0, longitude: 0 });
  }, []);

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
      darkmode={darkmode}
      onLocationChange={fetchTilesForLocation}
      selectedTileID={selectedTileID}
    />
  );
};

