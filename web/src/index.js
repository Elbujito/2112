import axios from "axios";
import './main.css';

var Cesium = require('cesium/Cesium');
require('./main.css');
require('cesium/Widgets/widgets.css'); // Import Cesium's CSS

let viewer; // Declare viewer globally for access after data fetch

// Fetch satellite orbit data from API
async function fetchSatelliteOrbit(noradID) {
  const SATELLITE_API_URL = "http://localhost:8081/satellites/orbit";

  try {
    const response = await axios.get(`${SATELLITE_API_URL}`, {
      params: { noradID },
      headers: { Accept: "application/json" },
    });

    if (response.status === 200 && Array.isArray(response.data.payload)) {
      return response.data.payload.map((data) => ({
        latitude: data.Latitude,
        longitude: data.Longitude,
        altitude: data.Altitude,
        time: data.Time,
      }));
    } else {
      console.error("Unexpected API response structure:", response.data);
      return [];
    }
  } catch (error) {
    console.error("Error fetching satellite orbit data:", error.message);
    return [];
  }
}

// Initialize Cesium viewer
function initializeViewer() {
  viewer = new Cesium.Viewer("cesiumContainer", {
    terrainProvider: new Cesium.EllipsoidTerrainProvider(),
    timeline: true,
    animation: true,
  });
}

// Plot orbit data on Cesium viewer
async function plotOrbitFromAPI(noradID, orbitData) {
  const positionProperty = new Cesium.SampledPositionProperty();

  orbitData.forEach(({ latitude, longitude, altitude, time }) => {
    const position = Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000);
    const julianTime = Cesium.JulianDate.fromIso8601(time);

    if (!julianTime) {
      console.error("Invalid time format:", time);
      return;
    }

    positionProperty.addSample(julianTime, position);
  });

  const orbitPositions = orbitData.map(({ latitude, longitude, altitude }) =>
    Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000)
  );

  viewer.entities.add({
    name: `Orbit Path for NORAD ID: ${noradID}`,
    polyline: {
      positions: orbitPositions,
      width: 1,
      material: Cesium.Color.YELLOW,
    },
  });

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
      show: true,
      leadTime: 60 * 60,
      trailTime: 60 * 60 * 24,
      resolution: 1,
      material: Cesium.Color.RED.withAlpha(0.5),
    },
  });

  const startTime = Cesium.JulianDate.fromIso8601(orbitData[0].time);
  const stopTime = Cesium.JulianDate.fromIso8601(orbitData[orbitData.length - 1].time);

  viewer.clock.startTime = startTime;
  viewer.clock.stopTime = stopTime;
  viewer.clock.currentTime = startTime;
  viewer.clock.clockRange = Cesium.ClockRange.LOOP_STOP;
  viewer.clock.multiplier = 60;
  viewer.timeline.zoomTo(startTime, stopTime);

  viewer.zoomTo(viewer.entities);
}

// Main application flow
async function main() {
  console.log("Fetching data...");
  const noradID = "25544";

  const orbitData = await fetchSatelliteOrbit(noradID);

  if (orbitData.length === 0) {
    console.error("No data available for NORAD ID:", noradID);
    return;
  }

  console.log("Initializing viewer...");
  initializeViewer();

  console.log("Plotting orbit data...");
  await plotOrbitFromAPI(noradID, orbitData);

  console.log("Done!");
}

main();
