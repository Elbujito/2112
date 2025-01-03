import React, { useRef, useState, useEffect, useCallback } from "react";
import { Box } from "@chakra-ui/react";
import Map, { Source, Layer, NavigationControl, GeolocateControl, Marker, MapRef } from "react-map-gl";
import { FeatureCollection, Point, Polygon, Position } from "geojson";
import "mapbox-gl/dist/mapbox-gl.css";
import Card from "components/card";
import { Tile } from "types/tiles";
import { OrbitDataItem } from "types/satellites";
import supercluster from "supercluster";

const MAPBOX_TOKEN = process.env.MAPBOX_TOKEN;

interface MapTileCardProps {
  tiles: Tile[];
  darkmode: string;
  onLocationChange: (location: { latitude: number; longitude: number }) => void;
  selectedTileIDs?: string[]; // Updated to accept an array of selected tile IDs
  satellitePositionData?: Record<string, OrbitDataItem[]>; // Added satellite position data
}

const generateSquare = (lat: number, lon: number, size: number): Polygon => {
  const earthRadius = 6378137; // Earth's radius in meters for Mercator projection
  const halfSize = size;

  const mercatorX = (lon * Math.PI * earthRadius) / 180;
  const mercatorY =
    earthRadius * Math.log(Math.tan(Math.PI / 4 + (lat * Math.PI) / 360));

  const cornersMercator = [
    [mercatorX - halfSize, mercatorY - halfSize],
    [mercatorX + halfSize, mercatorY - halfSize],
    [mercatorX + halfSize, mercatorY + halfSize],
    [mercatorX - halfSize, mercatorY + halfSize],
    [mercatorX - halfSize, mercatorY - halfSize],
  ];

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
  selectedTileIDs = [], // Default to an empty array
  satellitePositionData,
}) => {
  const mapRef = useRef<MapRef | null>(null);
  const [hoveredTile, setHoveredTile] = useState<Tile | null>(null);
  const [clusters, setClusters] = useState<any[]>([]);
  const [clusterZoom, setClusterZoom] = useState<number>(5);

  useEffect(() => {
    if (selectedTileIDs.length > 0) {
      const firstSelectedTile = tiles.find((tile) => tile.ID === selectedTileIDs[0]);
      if (firstSelectedTile && mapRef.current) {
        mapRef.current.flyTo({
          center: [firstSelectedTile.CenterLon, firstSelectedTile.CenterLat],
          zoom: 3,
        });
      }
    }
  }, [selectedTileIDs, tiles]);

  // Prepare data for supercluster
  const getClusteredData = useCallback(() => {
    const points = [];

    // Flatten satellite position data into a point array
    satellitePositionData &&
      Object.keys(satellitePositionData).forEach((satelliteID) => {
        const positions = satellitePositionData[satelliteID];
        positions.forEach((position) => {
          points.push({
            type: "Feature",
            properties: { satelliteID },
            geometry: {
              type: "Point",
              coordinates: [position.longitude, position.latitude],
            },
          });
        });
      });

    const index = new supercluster({
      radius: 40,
      maxZoom: 16,
    });

    index.load(points);

    const bounds = mapRef.current?.getMap()?.getBounds().toArray();

    if (bounds) {
      const clusters = index.getClusters(bounds[0].concat(bounds[1]), clusterZoom);
      setClusters(clusters);
    }
  }, [satellitePositionData, clusterZoom]);

  // Update clusters when zoom level or satellite positions change
  useEffect(() => {
    getClusteredData();
  }, [satellitePositionData, clusterZoom, getClusteredData]);

  const handleZoomChange = () => {
    const currentZoom = mapRef.current?.getMap()?.getZoom();
    if (currentZoom) {
      setClusterZoom(Math.floor(currentZoom));
    }
  };

  const handleTileHover = (event: any) => {
    const features = event.features;
    if (features && features.length > 0) {
      const hoveredFeature = features[0];
      const properties = hoveredFeature.properties;

      if (properties?.id) {
        const tile = tiles.find((t) => t.ID === properties.id);
        if (tile) {
          setHoveredTile(tile);
        }
      }
    } else {
      setHoveredTile(null);
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
            latitude: 49.6117,
            longitude: 6.1319,
            zoom: 5,
          }}
          style={{
            borderRadius: "20px",
            width: "100%",
            height: "100%",
          }}
          mapStyle={darkmode}
          mapboxAccessToken={MAPBOX_TOKEN}
          interactiveLayerIds={["tile-boundaries"]}
          onMouseMove={handleTileHover}
          onZoomEnd={handleZoomChange} // Track zoom level changes
        >
          <Source id="tiles" type="geojson" data={geoJsonSource}>
            <Layer
              id="tile-boundaries"
              type="fill"
              paint={{
                "fill-color": [
                  "case",
                  ["in", ["get", "id"], ["literal", selectedTileIDs]],
                  "#0D47A1", // Highlight color for selected tiles
                  "#888888", // Default color
                ],
                "fill-opacity": 0.4,
              }}
            />
            <Layer
              id="tile-borders"
              type="line"
              paint={{
                "line-color": "#000000",
                "line-width": 2,
              }}
            />
          </Source>

          {/* Render clustered satellite position markers */}
          {clusters &&
            clusters.map((cluster: any, index: number) => {
              const [longitude, latitude] = cluster.geometry.coordinates;
              const count = cluster.properties?.point_count_abbreviated;

              return (
                <Marker key={index} latitude={latitude} longitude={longitude}>
                  <div
                    style={{
                      backgroundColor: "red",
                      padding: "10px",
                      borderRadius: "50%",
                      color: "white",
                      fontSize: "12px",
                      textAlign: "center",
                    }}
                  >
                    {count}
                  </div>
                </Marker>
              );
            })}

          <GeolocateControl
            position="top-right"
            onGeolocate={(position) => {
              const { latitude, longitude } = position.coords;
              onLocationChange({ latitude, longitude });
            }}
            trackUserLocation={true}
          />
          <NavigationControl position="top-right" />
        </Map>
      </Box>
    </Card>
  );
};

export default MapTileCard;
