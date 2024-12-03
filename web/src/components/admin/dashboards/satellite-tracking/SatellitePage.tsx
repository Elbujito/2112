import React, { useState } from "react";
import SatelliteTableWithData from "./SatelliteTableWithData"; // Assume this is the table component
import MapSatelliteWithData from "./MapSatelliteWithData"; // Assume this is the map component

const SatellitePage: React.FC = () => {
  const [selectedNoradID, setSelectedNoradID] = useState<string | null>(null);

  // Handler to update NORAD ID
  const handleNoradIDChange = (noradID: string) => {
    setSelectedNoradID(noradID);
  };

  return (
    <div className="mt-3 grid h-full w-full grid-cols-1 lg:grid-cols-10 gap-5">
      {/* Left side: Satellite Table (60% width) */}
      <div
        className="col-span-1 lg:col-span-6"
        style={{
          minHeight: "60vh",
          maxHeight: "60vh",
        }}
      >
        <SatelliteTableWithData onSelectNoradID={handleNoradIDChange} />
      </div>

      {/* Right side: Map (40% width) */}
      <div
        className="col-span-1 lg:col-span-4"
        style={{
          minHeight: "60vh",
          maxHeight: "60vh",
        }}
      >
        <MapSatelliteWithData noradID={selectedNoradID} />
      </div>
    </div>
  );
};

export default SatellitePage;
