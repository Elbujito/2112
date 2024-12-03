import React, { useState, useEffect } from "react";
import dynamic from "next/dynamic";
import axios from "axios";
import {
  CircularProgress,
  TextField,
  InputAdornment,
  IconButton,
  Box,
  Paper,
  Button,
} from "@mui/material";
import SearchIcon from "@mui/icons-material/Search";
import SatelliteTable from "./SatelliteTable";

type OrbitDataItem = {
  latitude: number;
  longitude: number;
  altitude: number;
  time: string;
};

export type Satellite = {
  id: string;
  name: string;
  noradID: string;
  owner: string;
  launchDate: string | null;
  apogee: number | null;
  perigee: number | null;
};

interface CesiumViewerProps {
  orbitData?: OrbitDataItem[];
  noradID?: string;
}

const CesiumViewer = dynamic<CesiumViewerProps>(
  () => import("../../shared/cesium/cesium"),
  {
    ssr: false,
    loading: () => <div style={{ height: "80vh", background: "#000" }} />,
  }
);

const SatelliteTrack: React.FC = () => {
  const [noradID, setNoradID] = useState<string>(""); // Current NORAD ID for orbit data
  const [orbitData, setOrbitData] = useState<OrbitDataItem[]>([]); // Orbit data for visualization
  const [satellites, setSatellites] = useState<Satellite[]>([]); // Satellite list
  const [loading, setLoading] = useState<boolean>(false); // Loading state for satellite list
  const [fetchingOrbit, setFetchingOrbit] = useState<boolean>(false); // Loading state for orbit data
  const [search, setSearch] = useState<string>(""); // Search term for filtering satellites
  const [page, setPage] = useState<number>(1); // Current page for pagination
  const [isTableOpen, setIsTableOpen] = useState<boolean>(true); // State for table visibility

  // Fetch satellite list on component mount and when search or page changes
  useEffect(() => {
    const fetchSatellites = async () => {
      setLoading(true);
      try {
        const response = await axios.get("http://localhost:8081/satellites/paginated", {
          params: { page, pageSize: 10, search }, // Include search in query parameters
          headers: { Accept: "application/json" }, // Explicitly set Accept header to fix 406 issue
        });
        setSatellites(response.data.satellites || []);
      } catch (error) {
        console.error("Error fetching satellites:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchSatellites();
  }, [search, page]);

  // Fetch orbit data for a specific NORAD ID
  const fetchOrbitData = async (noradID: string) => {
    setFetchingOrbit(true);
    try {
      const response = await axios.get("http://localhost:8081/satellites/orbit", {
        params: { noradID },
        headers: { Accept: "application/json" }, // Explicitly set Accept header to fix 406 issue
      });
      setOrbitData(
        response.data.payload.map((data: any) => ({
          latitude: data.Latitude,
          longitude: data.Longitude,
          altitude: data.Altitude,
          time: data.Time,
        }))
      );
      setNoradID(noradID);
    } catch (error) {
      console.error("Error fetching orbit data:", error);
    } finally {
      setFetchingOrbit(false);
    }
  };

  // Handle search input changes
  const handleSearchChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(event.target.value);
    setPage(1); // Reset to the first page when searching
  };

  return (
    <div className="p-6 relative">
      {/* Search Bar */}
      <div className="mb-6">
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Search satellites..."
          value={search}
          onChange={handleSearchChange}
          InputProps={{
            endAdornment: (
              <InputAdornment position="end">
                <IconButton>
                  <SearchIcon />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
      </div>

      {/* Cesium Viewer */}
      <div
        className="relative w-full"
        style={{
          height: isTableOpen ? "60vh" : "80vh", // Adjust viewer height based on table visibility
          transition: "height 0.3s ease",
        }}
      >
        <CesiumViewer orbitData={orbitData} noradID={noradID} />

        {fetchingOrbit && (
          <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
            <CircularProgress color="primary" />
          </div>
        )}
      </div>

      {/* Toggle Button */}
      <div className="text-center mt-4">
        <Button
          variant="contained"
          onClick={() => setIsTableOpen((prev) => !prev)}
        >
          {isTableOpen ? "Hide Table" : "Show Table"}
        </Button>
      </div>

      {/* Satellite Table (Expandable) */}
      <Box
        component={Paper}
        elevation={4}
        sx={{
          position: "relative",
          overflow: "hidden",
          height: isTableOpen ? "40vh" : "0vh",
          transition: "height 0.3s ease",
          display: isTableOpen ? "block" : "none",
        }}
      >
        {loading ? (
          <div className="flex justify-center p-6">
            <CircularProgress color="primary" />
          </div>
        ) : (
          <SatelliteTable satellites={satellites} onViewOrbit={fetchOrbitData} />
        )}
      </Box>
    </div>
  );
};

export default SatelliteTrack;
