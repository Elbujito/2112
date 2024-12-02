import React, { useEffect } from "react";
import * as Cesium from "cesium";
import "cesium/Build/Cesium/Widgets/widgets.css";
import Card from "components/card";

interface OrbitData {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
}

interface MapSatelliteProps {
  orbitData?: OrbitData[]; // Optional orbit data
  noradID?: string; // Optional NORAD ID
  userLocation?: { latitude: number; longitude: number }; // User location
}

const MapSatellite: React.FC<MapSatelliteProps> = ({
  orbitData = [],
  noradID = "Unknown",
  userLocation,
}) => {
  useEffect(() => {
    let viewer: Cesium.Viewer | null = new Cesium.Viewer("cesiumContainer", {
      terrainProvider: new Cesium.EllipsoidTerrainProvider(),
      timeline: true,
      animation: true,
      sceneMode: Cesium.SceneMode.SCENE3D, // Set default view to 3D
    });

    const plotOrbit = () => {
      if (!viewer) return;

      if (orbitData.length === 0) {
        console.warn("Orbit data not provided. Viewer will display Earth only.");
      } else {
        const positionProperty = new Cesium.SampledPositionProperty();

        orbitData.forEach(({ latitude, longitude, altitude, time }) => {
          try {
            const position = Cesium.Cartesian3.fromDegrees(
              longitude,
              latitude,
              altitude * 1000
            );
            const julianTime = Cesium.JulianDate.fromIso8601(time);
            positionProperty.addSample(julianTime, position);
          } catch (error) {
            console.error("Error processing orbit data:", error);
          }
        });

        const orbitPositions = orbitData.map(({ latitude, longitude, altitude }) =>
          Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000)
        );

        viewer.entities.add({
          name: `Orbit Path for Satellite ${noradID}`,
          polyline: {
            positions: orbitPositions,
            width: 1,
            material:Cesium.Color.PURPLE,
          },
        });

        viewer.entities.add({
          name: `Satellite ${noradID}`,
          position: positionProperty,
          point: {
            pixelSize: 12,
            color: Cesium.Color.PURPLE, // Purple color for the satellite
          },
          label: {
            text: `Satellite: ${noradID}`,
            font: "16pt Arial",
            fillColor: Cesium.Color.WHITE,
            style: Cesium.LabelStyle.FILL,
            outlineWidth: 2,
            outlineColor: Cesium.Color.BLACK,
            showBackground: true,
            backgroundColor: Cesium.Color.DARKSLATEGRAY.withAlpha(0.7),
            verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
            pixelOffset: new Cesium.Cartesian2(0, -20),
          },
          path: {
            material: Cesium.Color.PURPLE, 
            width: 2,
          },
        });

        try {
          const startTime = Cesium.JulianDate.fromIso8601(orbitData[0]?.time ?? "");
          const stopTime = Cesium.JulianDate.fromIso8601(
            orbitData[orbitData.length - 1]?.time ?? ""
          );

          if (startTime && stopTime) {
            viewer.clock.startTime = startTime.clone();
            viewer.clock.stopTime = stopTime.clone();
            viewer.clock.currentTime = startTime.clone();
            viewer.clock.clockRange = Cesium.ClockRange.LOOP_STOP;
            viewer.clock.multiplier = 1; // x1 playback speed
            viewer.clock.shouldAnimate = true; // Start animation by default
            viewer.timeline.zoomTo(startTime, stopTime);
          } else {
            console.warn("Invalid orbit data times. Ensure time format is ISO8601.");
          }
        } catch (error) {
          console.error("Error setting clock or camera:", error);
        }
      }

      // Add user location to the map, if available
      if (userLocation) {
        viewer.entities.add({
          name: "User Location",
          position: Cesium.Cartesian3.fromDegrees(
            userLocation.longitude,
            userLocation.latitude,
            0 // Altitude 0 for ground level
          ),
          point: {
            pixelSize: 10,
            color: Cesium.Color.BLUE, // Blue color for user
            outlineColor: Cesium.Color.WHITE,
            outlineWidth: 2,
          },
          label: {
            text: "Your Location",
            font: "14pt Arial",
            fillColor: Cesium.Color.WHITE,
            style: Cesium.LabelStyle.FILL,
            outlineWidth: 2,
            outlineColor: Cesium.Color.BLACK,
            showBackground: true,
            backgroundColor: Cesium.Color.DARKSLATEGRAY.withAlpha(0.7),
            verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
            pixelOffset: new Cesium.Cartesian2(0, -20),
          },
        });
      }

      // Set the camera view to focus on the entire Earth in 3D
      viewer.scene.camera.setView({
        destination: Cesium.Cartesian3.fromDegrees(0, 0, 20000000), // Height is in meters
      });
    };

    plotOrbit();

    return () => {
      if (viewer) {
        viewer.destroy();
        viewer = null;
      }
    };
  }, [orbitData, noradID, userLocation]);

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"} style={{ borderRadius: "20px" }}>
      <style>
        {`
          .cesium-viewer .cesium-widget-credits {
            display: none !important;
          }
        `}
      </style>
      <div id="cesiumContainer" style={{ width: "100%", height: "100%", borderRadius: "20px" }} />
    </Card>
  );
};

export default MapSatellite;
