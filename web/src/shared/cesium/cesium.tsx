import React, { useEffect } from "react";
import * as Cesium from "cesium";
import 'cesium/Build/Cesium/Widgets/widgets.css';

interface OrbitData {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
}

interface CesiumViewerProps {
  orbitData?: OrbitData[]; // Make orbitData optional
  noradID?: string;        // Make noradID optional
}

const CesiumViewer: React.FC<CesiumViewerProps> = ({ orbitData = [], noradID = "Unknown" }) => {
  useEffect(() => {
    let viewer: Cesium.Viewer | null = null;

    const initializeViewer = () => {
      viewer = new Cesium.Viewer("cesiumContainer", {
        terrainProvider: new Cesium.EllipsoidTerrainProvider(),
        timeline: true,
        animation: true,
      });

      // Adjust the camera to show Earth properly
      viewer.scene.camera.setView({
        destination: Cesium.Cartesian3.fromDegrees(0, 0, 20000000), // Overhead view
      });
    };

    const plotOrbit = () => {
      if (!viewer || orbitData.length === 0) {
        console.warn("Orbit data not provided. Viewer will display Earth only.");
        return;
      }

      const positionProperty = new Cesium.SampledPositionProperty();

      orbitData.forEach(({ latitude, longitude, altitude, time }) => {
        const position = Cesium.Cartesian3.fromDegrees(
          longitude,
          latitude,
          altitude * 1000
        );
        const julianTime = Cesium.JulianDate.fromIso8601(time);

        if (!julianTime) {
          console.error("Invalid time format:", time);
          return;
        }

        positionProperty.addSample(julianTime, position);
      });

      // Create polyline for orbit path
      const orbitPositions = orbitData.map(({ latitude, longitude, altitude }) =>
        Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000)
      );

      viewer.entities.add({
        name: `Orbit Path for NORAD ID: ${noradID}`,
        polyline: {
          positions: orbitPositions,
          width: 2,
          material: Cesium.Color.YELLOW,
        },
      });

      // Add satellite entity with path
      viewer.entities.add({
        name: `Satellite ${noradID}`,
        position: positionProperty,
        point: {
          pixelSize: 10,
          color: Cesium.Color.RED,
        },
        label: {
          text: `Satellite ${noradID}`,
          font: "14pt sans-serif",
          fillColor: Cesium.Color.WHITE,
          showBackground: true,
          backgroundColor: Cesium.Color.BLACK.withAlpha(0.7),
          verticalOrigin: Cesium.VerticalOrigin.BOTTOM,
          pixelOffset: new Cesium.Cartesian2(0, -15),
        },
        path: {
          material: Cesium.Color.RED.withAlpha(0.5),
          width: 2,
        },
      });

      // Configure the clock and timeline
      const startTime = Cesium.JulianDate.fromIso8601(orbitData[0]?.time ?? "");
      const stopTime = Cesium.JulianDate.fromIso8601(
        orbitData[orbitData.length - 1]?.time ?? ""
      );

      if (startTime && stopTime) {
        viewer.clock.startTime = startTime;
        viewer.clock.stopTime = stopTime;
        viewer.clock.currentTime = startTime;
        viewer.clock.clockRange = Cesium.ClockRange.LOOP_STOP;
        viewer.clock.multiplier = 60;
        viewer.timeline.zoomTo(startTime, stopTime);

        // Zoom to entities
        viewer.zoomTo(viewer.entities);
      }
    };

    // Initialize and plot orbit
    initializeViewer();
    plotOrbit();

    // Cleanup function to destroy the viewer
    return () => {
      if (viewer) {
        viewer.destroy();
        viewer = null;
      }
    };
  }, [orbitData, noradID]);

  return <div id="cesiumContainer" style={{ width: "100%", height: "100%" }} />;
};

export default CesiumViewer;
