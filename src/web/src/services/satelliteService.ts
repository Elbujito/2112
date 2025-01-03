import { create } from "zustand";
import apiClient from "../utils/apiClient";
import { SatelliteInfo } from "types/satellites";

interface SatelliteServiceState {
    satelliteInfo: SatelliteInfo[];
    totalSatelliteInfo: number;
    loading: boolean;
    error: string | null;
    fetchPaginatedSatelliteInfo: (pageIndex: number, pageSize: number, search: string) => Promise<void>;
}

const useSatelliteServiceStore = create<SatelliteServiceState>((set) => ({
    satelliteInfo: [],
    totalSatelliteInfo: 0,
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
}));

export default useSatelliteServiceStore;
