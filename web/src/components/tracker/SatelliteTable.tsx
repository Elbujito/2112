import React from "react";
import {
  DataGrid,
  GridColDef,
} from "@mui/x-data-grid";
import { Button, Typography, useTheme } from "@mui/material";
import { Satellite } from "./Tracker";

interface SatelliteTableProps {
  satellites: Satellite[];
  onViewOrbit: (noradID: string) => void;
}

const SatelliteTable: React.FC<SatelliteTableProps> = ({ satellites, onViewOrbit }) => {
  const theme = useTheme(); // Access Material-UI theme

  // Define columns for the DataGrid
  const columns: GridColDef[] = [
    { field: "Name", headerName: "Name", flex: 1 },
    { field: "NoradID", headerName: "NORAD ID", flex: 1 },
    { field: "Owner", headerName: "Owner", flex: 1 },
    {
      field: "LaunchDate",
      headerName: "Launch Date",
      flex: 1,
      valueFormatter: (value?: string) => {
        return value? new Date(value).toLocaleDateString() : "N/A"}},
    { field: "Apogee", headerName: "Apogee (km)", flex: 1 },
    { field: "Perigee", headerName: "Perigee (km)", flex: 1 },
    {
      field: "Actions",
      headerName: "Actions",
      flex: 1,
      renderCell: (params) => (
        <Button
          variant="contained"
          size="small"
          onClick={() => onViewOrbit(params.row.noradID)}
        >
          View Orbit
        </Button>
      ),
      sortable: false,
      filterable: false,
    },
  ];

  console.log(satellites)

  return (
    <div style={{ height: 500, width: "100%" }}>
      <DataGrid
        rows={satellites.map((satellite) => ({ ...satellite, id: satellite.noradID }))}
        columns={columns}
        paginationModel={{ pageSize: 10, page: 0 }}
        pageSizeOptions={[5, 10, 20]}
        getRowId={(row) => row.ID}
        sx={{
          backgroundColor: theme.palette.background.paper,
          "& .MuiDataGrid-columnHeaders": {
            backgroundColor: theme.palette.action.hover,
          },
        }}
      />
    </div>
  );
};

export default SatelliteTable;
