import React, { useState } from "react";
import SatelliteTableWithData from "./SatelliteTableWithData"; // Assume this is the table component
import MapSatelliteWithData from "./MapSatelliteWithData"; // Assume this is the map component
import MapCard from "./MapCard";
import VisibilityTimeline from "./VisibilityTimeline"; // Assume this is the timeline component
import { CustomScrollbar } from "components/scrollbar/CustomScrollbar";

const SatellitePage: React.FC = () => {
  const [selectedNoradID, setSelectedNoradID] = useState<string | null>(null);
  const [userLocation, setUserLocation] = useState<{ latitude: number; longitude: number } | null>(
    null
  );

  // Handler to update NORAD ID
  const handleNoradIDChange = (noradID: string) => {
    setSelectedNoradID(noradID);
  };

  // Handler to update user location
  const handleLocationChange = (location: { latitude: number; longitude: number }) => {
    setUserLocation(location);
  };

  // Mock visibility data for the timeline
  const mockVisibilityData = [
    {
      satellite: "Hubble Space Telescope",
      noradID: "20580",
      day: "03",
      month: "12",
      weekday: "Wed",
      hours: "10:30 - 12:00",
      current: true,
    },
    {
      satellite: "ISS",
      noradID: "25544",
      day: "04",
      month: "12",
      weekday: "Thu",
      hours: "09:00 - 09:15",
    },
    {
      satellite: "Starlink-1234",
      noradID: "44238",
      day: "05",
      month: "12",
      weekday: "Fri",
      hours: "21:30 - 22:00",
    },
    {
      satellite: "Sentinel-2A",
      noradID: "40697",
      day: "06",
      month: "12",
      weekday: "Sat",
      hours: "15:00 - 15:30",
    },
    {
      satellite: "Landsat 9",
      noradID: "50294",
      day: "07",
      month: "12",
      weekday: "Sun",
      hours: "17:45 - 18:15",
    },
    {
      satellite: "TerraSAR-X",
      noradID: "31698",
      day: "08",
      month: "12",
      weekday: "Mon",
      hours: "14:20 - 14:50",
    },
    {
      satellite: "NOAA-20",
      noradID: "43013",
      day: "09",
      month: "12",
      weekday: "Tue",
      hours: "13:30 - 13:50",
    },
    {
      satellite: "Gaofen-4",
      noradID: "41019",
      day: "10",
      month: "12",
      weekday: "Wed",
      hours: "22:00 - 22:30",
    },
  ];

  return (
    <div className="mt-3 grid h-full w-full grid-rows-[70vh_auto] grid-cols-12 gap-5">
      {/* Left side: MapCard (spans 7 columns) */}
      <div className="row-span-1 col-span-12 lg:col-span-7">
        <div className="h-full">
          <MapCard onLocationChange={handleLocationChange} />
        </div>
      </div>

      {/* Right side: Satellite Table and Map (spans 5 columns) */}
      <div className="row-span-1 col-span-12 lg:col-span-5 grid grid-rows-2 gap-5">
        {/* Satellite Table */}
        <div className="row-span-1 overflow-auto">
          <CustomScrollbar style={{ height: "100%" }}>
            <SatelliteTableWithData onSelectNoradID={handleNoradIDChange} />
          </CustomScrollbar>
        </div>

        {/* Satellite View */}
        <div className="row-span-1">
          <MapSatelliteWithData
            noradID={selectedNoradID}
            userLocation={userLocation} // Pass user location as a prop
          />
        </div>
      </div>

      {/* Bottom row: Visibility Timeline (spans all 12 columns) */}
      <div className="row-span-1 col-span-12">
        <VisibilityTimeline data={mockVisibilityData} />
      </div>
    </div>
  );
};

export default SatellitePage;
