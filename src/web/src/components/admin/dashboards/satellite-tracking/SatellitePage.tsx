import React, { useState } from "react";
import SatelliteTableWithData from "./SatelliteTableWithData"; // Table Component
import MapSatelliteWithData from "./MapSatelliteWithData"; // Map Component
import MapCard from "./MapCard"; // Map with Location Selector
import { CustomScrollbar } from "components/scrollbar/CustomScrollbar"; // Custom Scrollbar Component
import VisibilityTimelineWithData from "./VisibilityTimelineWithData"; // Timeline Component

const SatellitePage: React.FC = () => {
  // State for selected NORAD ID
  const [selectedNoradID, setSelectedNoradID] = useState<string | null>(null);

  // State for user location
  const [userLocation, setUserLocation] = useState<{ latitude: number; longitude: number }>({
    latitude: 37.7749, // Default to San Francisco
    longitude: -122.4194,
  });

  // Handle NORAD ID selection from the Satellite Table
  const handleNoradIDChange = (noradID: string) => {
    setSelectedNoradID(noradID);
  };

  // Handle user location change from the MapCard
  const handleLocationChange = (location: { latitude: number; longitude: number }) => {
    setUserLocation(location);
  };

  return (
    <div className="mt-3 grid h-full w-full grid-rows-[70vh_auto] grid-cols-12 gap-5">
      {/* Left Side: MapCard (spans 7 columns) */}
      <div className="row-span-1 col-span-12 lg:col-span-7">
        <div className="h-full">
          <MapCard onLocationChange={handleLocationChange} />
        </div>
      </div>

      {/* Right Side: Satellite Table and Satellite Map (spans 5 columns) */}
      <div className="row-span-1 col-span-12 lg:col-span-5 grid grid-rows-2 gap-5">
        {/* Satellite Table */}
        <div className="row-span-1 overflow-auto">
          <CustomScrollbar style={{ height: "100%" }}>
            <SatelliteTableWithData onSelectNoradID={handleNoradIDChange} />
          </CustomScrollbar>
        </div>

        {/* Satellite Map */}
        <div className="row-span-1">
          <MapSatelliteWithData
            noradID={selectedNoradID}
            userLocation={userLocation}
          />
        </div>
      </div>

      {/* Bottom Row: Visibility Timeline (spans all 12 columns) */}
      <div className="row-span-1 col-span-12">
        <VisibilityTimelineWithData userLocation={userLocation} uid={"adrien-test"} />
      </div>
    </div>
  );
};

export default SatellitePage;
