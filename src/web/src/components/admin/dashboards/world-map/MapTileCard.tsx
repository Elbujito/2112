import React, { useRef } from "react";
import { Box } from "@chakra-ui/react";
import Map, { Source, Layer, NavigationControl, GeolocateControl } from "react-map-gl";
import { FeatureCollection, Polygon, Position } from "geojson";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "components/card";

const MAPBOX_TOKEN =
  "pk.eyJ1Ijoic2ltbW1wbGUiLCJhIjoiY2wxeG1hd24xMDEzYzNrbWs5emFkdm16ZiJ9.q9s0sSKQFFaT9fyrC-7--g"; // Replace with your Mapbox token

interface Tile {
  Quadkey: string;
  ZoomLevel: number;
  CenterLat: number;
  CenterLon: number;
  SpatialIndex?: string;
  NbFaces: number;
  Radius: number;
  BoundariesJSON?: string;
}

interface MapTileCardProps {
  tiles: Tile[];
  darkmode: string;
  onLocationChange: (location: { latitude: number; longitude: number }) => void; // Callback for location change
}

const generateSquare = (lat: number, lon: number, size: number): Polygon => {
    const earthRadius = 6378137; // Earth's radius in meters for Mercator projection

    const halfSize = size;

    // Convert latitude and longitude to Mercator x and y
    const mercatorX = (lon * Math.PI * earthRadius) / 180;
    const mercatorY =
      earthRadius *
      Math.log(Math.tan(Math.PI / 4 + (lat * Math.PI) / 360));

    // Create square in Mercator coordinates
    const cornersMercator = [
      [mercatorX - halfSize, mercatorY - halfSize], // Bottom-left
      [mercatorX + halfSize, mercatorY - halfSize], // Bottom-right
      [mercatorX + halfSize, mercatorY + halfSize], // Top-right
      [mercatorX - halfSize, mercatorY + halfSize], // Top-left
      [mercatorX - halfSize, mercatorY - halfSize], // Close the polygon
    ];

    // Convert Mercator coordinates back to latitude and longitude
    const cornersLatLon = cornersMercator.map(([x, y]) => {
      const lonDeg = (x * 180) / (Math.PI * earthRadius);
      const latRad = (2 * Math.atan(Math.exp(y / earthRadius))) - Math.PI / 2;
      const latDeg = (latRad * 180) / Math.PI;
      return [lonDeg, latDeg] as [number, number];
    });

    // Create GeoJSON Polygon
    const coordinates: Position[][] = [cornersLatLon];

    return {
      type: "Polygon",
      coordinates,
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
      geometry: generateSquare(tile.CenterLat, tile.CenterLon, tile.Radius), // Generate circle geometry
      properties: {
        quadkey: tile.Quadkey,
        zoomLevel: tile.ZoomLevel,
        nbFaces: tile.NbFaces,
        radius: tile.Radius,
      },
    })),
  };

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
                "fill-opacity": 0.2,
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
