import React, { useState, useEffect } from "react";
import { Map } from "react-map-gl";
import "mapbox-gl/dist/mapbox-gl.css";
import * as Cesium from "cesium";
import "cesium/Build/Cesium/Widgets/widgets.css";
import Card from "components/card";

const MAPBOX_TOKEN =
  "pk.eyJ1Ijoic2ltbW1wbGxlIiwiYWQiOiJx9ssz4sR--YourTokenHere--"; // Set your Mapbox token

const MapCard = ({ orbitData, noradID }) => {
  const [darkmode, setDarkmode] = useState(
    document.body.classList.contains("dark")
      ? "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
      : "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
  );

  useEffect(() => {
    const observer = new MutationObserver((mutationsList) => {
      for (const mutation of mutationsList) {
        if (
          mutation.type === "attributes" &&
          mutation.attributeName === "class"
        ) {
          if (document.body.classList.contains("dark")) {
            setDarkmode(
              "mapbox://styles/simmmple/cl0qqjr3z000814pq7428ptk5"
            );
          } else {
            setDarkmode(
              "mapbox://styles/simmmple/ckwxecg1wapzp14s9qlus38p0"
            );
          }
        }
      }
    });
    observer.observe(document.body, { attributes: true });
    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    let viewer = null;

    if (orbitData && orbitData.length > 0) {
      viewer = new Cesium.Viewer("cesiumContainer", {
        terrainProvider: new Cesium.EllipsoidTerrainProvider(),
        timeline: true,
        animation: true,
      });

      const positionProperty = new Cesium.SampledPositionProperty();

      orbitData.forEach(({ latitude, longitude, altitude, time }) => {
        const position = Cesium.Cartesian3.fromDegrees(
          longitude,
          latitude,
          altitude * 1000
        );
        const julianTime = Cesium.JulianDate.fromIso8601(time);
        positionProperty.addSample(julianTime, position);
      });

      viewer.entities.add({
        position: positionProperty,
        point: { pixelSize: 10, color: Cesium.Color.RED },
        label: {
          text: `Satellite ${noradID || "Unknown"}`,
          font: "14pt sans-serif",
          fillColor: Cesium.Color.WHITE,
        },
        path: {
          material: Cesium.Color.RED.withAlpha(0.5),
          width: 2,
        },
      });

      return () => {
        if (viewer) {
          viewer.destroy();
        }
      };
    }
  }, [orbitData, noradID]);

  return (
    <Card extra={"relative w-full h-full bg-white px-3 py-[18px]"}>
      <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
        <div style={{ borderRadius: "20px", overflow: "hidden" }}>
          <Map
            initialViewState={{
              latitude: 37.692,
              longitude: -122.435,
              zoom: 13,
            }}
            style={{ width: "100%", minHeight: "300px" }}
            mapStyle={darkmode}
            mapboxAccessToken={MAPBOX_TOKEN}
          />
        </div>
        <div
          id="cesiumContainer"
          style={{ width: "100%", minHeight: "300px", borderRadius: "20px" }}
        />
      </div>
    </Card>
  );
};

export default MapCard;
