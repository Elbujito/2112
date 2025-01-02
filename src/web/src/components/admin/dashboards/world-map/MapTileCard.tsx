import React, { useRef, useState, useEffect } from "react";
import { Box, Text } from "@chakra-ui/react";
import Map, { Source, Layer, NavigationControl, GeolocateControl, Popup } from "react-map-gl";
import { FeatureCollection, Polygon, Position } from "geojson";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "components/card";

const MAPBOX_TOKEN =
  "pk.eyJ1Ijoic2ltbW1wbGUiLCJhIjoiY2wxeG1hd24xMDEzYzNrbWs5emFkdm16ZiJ9.q9s0sSKQFFaT9fyrC-7--g"; // Replace with your Mapbox token

interface Tile {
  ID: string;
  Quadkey: string;
  ZoomLevel: number;
  CenterLat: number;
  CenterLon: number;
  NbFaces: number;
  Radius: number;
}

interface MapTileCardProps {
  tiles: Tile[];
  darkmode: string;
  onLocationChange: (location: { latitude: number; longitude: number }) => void; // Callback for location change
  selectedMappingID?: string;
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

const MapTileCard: React.FC<MapTileCardProps> = ({
  tiles,
  darkmode,
  onLocationChange,
  selectedMappingID,
}) => {
  const mapRef = useRef(null);
  const [hoveredTile, setHoveredTile] = useState<Tile | null>(null);

  useEffect(() => {
    if (selectedMappingID) {
      const selectedTile = tiles.find((tile) => tile.ID === selectedMappingID);
      if (selectedTile && mapRef.current) {
        mapRef.current.flyTo({
          center: [selectedTile.CenterLon, selectedTile.CenterLat],
          zoom: 8,
        });
      }
    }
  }, [selectedMappingID, tiles]);

  const handleTileHover = (event: any) => {
    const features = event.features;
    if (features && features.length > 0) {
      const hoveredFeature = features[0];
      const properties = hoveredFeature.properties;

      if (properties?.quadkey) {
        // Match the hovered feature to the tile in the tiles array
        const tile = tiles.find((t) => t.Quadkey === properties.quadkey);
        if (tile) {
          setHoveredTile(tile);
        }
      }
    } else {
      setHoveredTile(null); // Clear hover state if no feature is hovered
    }
  };

  const geoJsonSource: FeatureCollection = {
    type: "FeatureCollection",
    features: tiles.map((tile) => ({
      type: "Feature",
      geometry: generateSquare(tile.CenterLat, tile.CenterLon, tile.Radius),
      properties: {
        id: tile.ID,
        quadkey: tile.Quadkey,
        zoomLevel: tile.ZoomLevel,
        nbFaces: tile.NbFaces,
        radius: tile.Radius,
      },
    })),
  };

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"}>
      <Box position="relative" w="100%" h="100%" overflow="hidden" bg="gray.50" borderRadius="md">
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
          interactiveLayerIds={["tile-boundaries"]} // Make this layer interactable
          onMouseMove={handleTileHover} // Handle hover events
        >
          <Source id="tiles" type="geojson" data={geoJsonSource}>
            <Layer
              id="tile-boundaries"
              type="fill"
              paint={{
                "fill-color": [
                  "case",
                  ["==", ["get", "id"], selectedMappingID],
                  "#FF5733", // Highlight color for selected tile
                  "#888888", // Default fill color
                ],
                "fill-opacity": 0.4,
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
          <GeolocateControl
            position="top-right"
            onGeolocate={(position) => {
              const { latitude, longitude } = position.coords;
              onLocationChange({ latitude, longitude });
            }}
            trackUserLocation={true}
          />
          <NavigationControl position="top-right" />
          {hoveredTile && (
            <Popup
              longitude={hoveredTile.CenterLon} // Use center coordinates for the popup
              latitude={hoveredTile.CenterLat}
              closeButton={false}
              closeOnClick={false}
              anchor="bottom"
              style={{
                zIndex: 1000,
                backgroundColor: "white",
                padding: "5px",
                borderRadius: "5px",
              }}
            >
              <Text fontSize="sm" fontWeight="bold" color="black">
                Tile ID: {hoveredTile.ID}
              </Text>
              <Text fontSize="xs" color="gray">
                Quadkey: {hoveredTile.Quadkey}
              </Text>
              <Text fontSize="xs" color="gray">
                Zoom Level: {hoveredTile.ZoomLevel}
              </Text>
              <Text fontSize="xs" color="gray">
                NbFaces: {hoveredTile.NbFaces}
              </Text>
              <Text fontSize="xs" color="gray">
                Radius: {hoveredTile.Radius}
              </Text>
            </Popup>
          )}
        </Map>
      </Box>
    </Card>
  );
};

export default MapTileCard;
