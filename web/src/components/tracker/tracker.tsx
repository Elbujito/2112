import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import axios from "axios";
import Skeleton from "@mui/material/Skeleton";
import CircularProgress from "@mui/material/CircularProgress";

type OrbitDataItem = {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string; // ISO 8601 date string
};

interface CesiumViewerProps {
  orbitData?: OrbitDataItem[];
  noradID?: string;
}

// Dynamic import for CesiumViewer with type annotations
const CesiumViewer = dynamic<CesiumViewerProps>(
  () => import("../../shared/cesium/cesium"),
  {
    ssr: false, // Ensure it only runs on the client
    loading: () => <Skeleton variant="rectangular" width="100%" height="80vh" />,
  }
);

const Tracker: React.FC = () => {
  const [noradID, setNoradID] = useState<string>("25544");
  const [currentNoradID, setCurrentNoradID] = useState<string>("25544");
  const [orbitData, setOrbitData] = useState<OrbitDataItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSatelliteOrbit = async () => {
      setLoading(true); // Start loading
      setError(null); // Clear previous errors

      try {
        const response = await axios.get("http://localhost:8081/satellites/orbit", {
          params: { noradID: currentNoradID },
          headers: { Accept: "application/json" },
        });

        if (response.status === 200 && Array.isArray(response.data.payload)) {
          setOrbitData(
            response.data.payload.map((data: any) => ({
              latitude: data.Latitude,
              longitude: data.Longitude,
              altitude: data.Altitude,
              time: data.Time,
            }))
          );
        } else {
          setError("Unexpected API response structure.");
        }
      } catch (err: any) {
        setError(`Error fetching data: ${err.message}`);
      } finally {
        setLoading(false); // End loading
      }
    };

    fetchSatelliteOrbit();
  }, [currentNoradID]);

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault();
    setCurrentNoradID(noradID);
  };

  return (
    <div className="p-6 relative">
      <div className="p-6 bg-gray-100">
        <form onSubmit={handleSubmit} className="mb-6">
          <label htmlFor="noradID" className="block text-gray-800 font-bold mb-2">
            Enter NORAD ID:
          </label>
          <div className="flex items-center">
            <input
              id="noradID"
              type="text"
              value={noradID}
              onChange={(e) => setNoradID(e.target.value)}
              className="flex-1 p-2 border border-gray-300 rounded text-gray-700 placeholder-gray-500 mr-4 focus:outline-none focus:ring-2 focus:ring-blue-400"
              placeholder="e.g., 25544"
            />
            <button
              type="submit"
              className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-400"
            >
              Fetch
            </button>
          </div>
        </form>
      </div>
      {error && <div className="text-red-500 mb-4">{error}</div>}

      <div
        className="relative w-full"
        style={{ height: "80vh" }}
      >
        {/* Cesium Viewer */}
        <CesiumViewer orbitData={orbitData} noradID={currentNoradID} />

        {/* Spinner Overlay (visible on hover) */}
        {loading && (
          <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-25 z-10 opacity-0 hover:opacity-100 transition-opacity">
            <CircularProgress color="primary" />
          </div>
        )}
      </div>
    </div>
  );
};

export default Tracker;
