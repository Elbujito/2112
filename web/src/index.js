import axios from "axios";

import './main.css';

var Cesium = require('cesium/Cesium');
require('./main.css');
require('cesium/Widgets/widgets.css'); // Import Cesium's CSS

// Cesium Viewer setup
const viewer = new Cesium.Viewer("cesiumContainer", {
  // terrainProvider: Cesium.createWorldTerrain(),
});


// API endpoint to fetch satellite data
const SATELLITE_API_URL = "http://localhost:8081/satellites/all";

// Function to fetch satellite data
async function fetchSatelliteData(noradID) {
  try {
    const response = await axios.get(`${SATELLITE_API_URL}`, {
      params: { noradID }, // Pass parameters properly
      headers: {
        Accept: 'application/json', // Ensure proper Accept header
      },
    });

    // Handle successful response
    if (response.status === 200) {
      if (response.data && response.data.payload) {
        return response.data.payload;
      } else {
        console.warn("Response received but payload is empty or missing.");
        return [];
      }
    } else {
      console.error(`Unexpected response status: ${response.status}`);
      return [];
    }
  } catch (error) {
    // Handle errors more thoroughly
    if (error.response) {
      // Server responded with a status code outside 2xx
      console.error(
        `Error fetching satellite data: ${error.response.status} - ${error.response.statusText}`
      );
      console.error("Server response:", error.response.data);
    } else if (error.request) {
      // Request was made but no response received
      console.error("No response received from server:", error.request);
    } else {
      // Other errors (e.g., invalid URL)
      console.error("Error setting up the request:", error.message);
    }
    return [];
  }
}

// Function to display satellite trajectory on Cesium Viewer
async function plotSatelliteTrajectory(noradID) {
  // Fetch data
  const data = await fetchSatelliteData(noradID);

  if (!data || data.length === 0) {
    console.error("No satellite data available.");
    return;
  }

  // Extract latitude, longitude, altitude, and time
  const positions = data.map((entry) => {
    const { latitude, longitude, altitude } = entry;
    const time = Cesium.JulianDate.fromIso8601(entry.time);
    return {
      position: Cesium.Cartesian3.fromDegrees(longitude, latitude, altitude * 1000),
      time,
    };
  });

  // Prepare position property
  const positionProperty = new Cesium.SampledPositionProperty();
  positions.forEach(({ position, time }) => {
    positionProperty.addSample(time, position);
  });

  // Add satellite trajectory to the viewer
  viewer.entities.add({
    name: `Satellite ${noradID}`,
    position: positionProperty,
    point: {
      pixelSize: 10,
      color: Cesium.Color.RED,
    },
    path: {
      show: true,
      leadTime: 0,
      trailTime: 60 * 60 * 24, // 24 hours
      width: 2,
      resolution: 1,
      material: new Cesium.ColorMaterialProperty(Cesium.Color.YELLOW),
    },
  });

  viewer.zoomTo(viewer.entities);
}

// Call the function for the International Space Station (NORAD ID 25544)
plotSatelliteTrajectory("25544");
