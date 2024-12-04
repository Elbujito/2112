import React, { useRef } from "react";
import { Box } from "@chakra-ui/react";
import Map, { Source, Layer, NavigationControl, GeolocateControl } from "react-map-gl";
import { FeatureCollection, Polygon } from "geojson";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "components/card";

const MAPBOX_TOKEN =
  "pk.eyJ1Ijoic2ltbW1wbGUiLCJhIjoiY2wxeG1hd24xMDEzYzNrbWs5emFkdm16ZiJ9.q9s0sSKQFFaT9fyrC-7--g"; // Replace with your Mapbox token

interface Tile {
  quadkey: string;
  zoomLevel: number;
  centerLat: number;
  centerLon: number;
  spatialIndex?: string;
  nbFaces: number;
  radius: number;
  boundariesJSON?: string;
}

interface MapTileCardProps {
  tiles: Tile[];
  darkmode: string;
  onLocationChange: (location: { latitude: number; longitude: number }) => void; // Callback for location change
}

// Helper function to generate a circle's GeoJSON geometry
const generateSquare = (lat: number, lon: number, size: number): Polygon => {
    const earthRadius = 6371000; // Earth's radius in meters
    const meterPerDegreeLat = 111320; // Approximation: 1 degree latitude ≈ 111.32 km
    const meterPerDegreeLon = Math.cos((lat * Math.PI) / 180) * (2 * Math.PI * earthRadius / 360);
  
    // Convert the size in meters to degrees
    const deltaLat = size / 2 / meterPerDegreeLat;
    const deltaLon = size / 2 / meterPerDegreeLon;

    console.log("deltaLat", deltaLat)
    console.log("deltaLon", deltaLon)
    console.log("meterPerDegreeLat", meterPerDegreeLat)
    console.log("meterPerDegreeLon", meterPerDegreeLon)

  
    // Create square coordinates
    const coordinates: [number, number][] = [
      [lon - deltaLon, lat - deltaLat], // Bottom-left
      [lon + deltaLon, lat - deltaLat], // Bottom-right
      [lon + deltaLon, lat + deltaLat], // Top-right
      [lon - deltaLon, lat + deltaLat], // Top-left
      [lon - deltaLon, lat - deltaLat], // Close the polygon (back to bottom-left)
    ];

    console.log("coordinates", coordinates)
  
    return {
      type: "Polygon",
      coordinates: [coordinates],
    };
  };

const MapTileCard: React.FC<MapTileCardProps> = ({ tiles, darkmode, onLocationChange }) => {
  const mapRef = useRef(null);

  const handleGeolocate = (position: GeolocationPosition) => {
    const { latitude, longitude } = position.coords;
    onLocationChange({ latitude, longitude }); // Pass location to parent
  };
  const geoJsonSource: FeatureCollection = {
    type: "FeatureCollection",
    features: tiles.map((tile) => ({
      type: "Feature",
      geometry: generateSquare(tile.centerLat, tile.centerLon, tile.radius), // Generate circle geometry
      properties: {
        quadkey: tile.quadkey,
        zoomLevel: tile.zoomLevel,
        nbFaces: tile.nbFaces,
        radius: tile.radius,
      },
    })),
  };

  console.log(geoJsonSource)
  console.log(tiles)

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
        <Map
          ref={mapRef}
          initialViewState={{
            latitude: 49.6117, // Default latitude
            longitude: 6.1319, // Default longitude
            zoom: 5, // Adjust initial zoom level
          }}
          style={{
            borderRadius: "20px",
            width: "100%",
            height: "100%",
          }}
          mapStyle={darkmode}
          mapboxAccessToken={MAPBOX_TOKEN}
        >
          <Source id="tiles" type="geojson" data={geoJsonSource}>
            <Layer
              id="tile-boundaries"
              type="fill"
              paint={{
                "fill-color": "#888888", // Fill color for the tile
                "fill-opacity": 0.5,
              }}
            />
            <Layer
              id="tile-borders"
              type="line"
              paint={{
                "line-color": "#000000", // Border color
                "line-width": 2,
              }}
            />
          </Source>
          {/* Geolocation and Navigation Controls */}
          <GeolocateControl
            position="top-right"
            onGeolocate={handleGeolocate} // Pass user's location to parent component
            trackUserLocation={true}
          />
          <NavigationControl position="top-right" />
        </Map>
      </Box>
    </Card>
  );
};

export default MapTileCard;
