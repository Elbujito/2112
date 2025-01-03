import { create } from "zustand";
import apiClient from "../utils/apiClient";
import { SatelliteInfo, OrbitDataItem } from "types/satellites";
import { gql } from "@apollo/client";
import AppoloClient from "lib/ApolloClient";

interface SatelliteServiceState {
    satelliteInfo: SatelliteInfo[];
    totalSatelliteInfo: number;
    orbitData: Record<string, OrbitDataItem[]>; // Orbit data keyed by NORAD ID
    loading: boolean;
    error: string | null;
    fetchPaginatedSatelliteInfo: (pageIndex: number, pageSize: number, search: string) => Promise<void>;
    fetchSatellitePositions: (noradID: string, startTime: string, endTime: string) => Promise<void>; // Method to fetch positions (void return type)
}

const GET_SATELLITE_POSITIONS = gql`
  query GetSatellitePositions($id: ID!, $startTime: String!, $endTime: String!) {
    satellitePositionsInRange(id: $id, startTime: $startTime, endTime: $endTime) {
      latitude
      longitude
      altitude
      timestamp
    }
  }
`;

const useSatelliteServiceStore = create<SatelliteServiceState>((set) => ({
    satelliteInfo: [],
    totalSatelliteInfo: 0,
    orbitData: {},
    loading: false,
    error: null,

    fetchPaginatedSatelliteInfo: async (pageIndex: number, pageSize: number, search: string) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/satellites/paginated/tles", {
                params: {
                    page: pageIndex + 1,
                    pageSize,
                    search,
                },
            });

            set({
                satelliteInfo: response.data?.satellites || [],
                totalSatelliteInfo: response.data?.totalRecords || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching satellite info:", err);
            set({
                error: "Failed to load satellite information.",
                loading: false,
            });
        }
    },

    fetchSatellitePositions: async (noradID: string, startTime: string, endTime: string): Promise<void> => {
        set({ loading: true, error: null });

        try {
            // Use Apollo Client's `query` method to fetch data
            const { data } = await AppoloClient.query({
                query: GET_SATELLITE_POSITIONS,
                variables: { id: noradID, startTime, endTime },
            });

            let orbitData: OrbitDataItem[] = [];
            if (data && Array.isArray(data.satellitePositionsInRange)) {
                orbitData = data.satellitePositionsInRange.map((item: any) => ({
                    latitude: item.latitude,
                    longitude: item.longitude,
                    altitude: item.altitude,
                    time: item.timestamp,
                }));
            } else {
                console.warn("Unexpected data structure:", data?.satellitePositionsInRange);
            }

            if (orbitData.length > 0) {
                set({ orbitData: { [noradID]: orbitData }, loading: false });
            } else {
                set({
                    error: "No satellite position data found.",
                    loading: false,
                });
            }
        } catch (err) {
            console.error("Error fetching satellite positions:", err);
            set({
                error: "Failed to load satellite positions.",
                loading: false,
            });
        }
    },
}));

export default useSatelliteServiceStore;
