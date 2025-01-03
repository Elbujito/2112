import { create } from "zustand";
import apiClient from "../utils/apiClient";
import { Tile, TileSatelliteMapping } from "types/tiles";

interface TileServiceState {
    tiles: Tile[];
    tileMappings: TileSatelliteMapping[];
    totalTileMappings: number;
    totalTiles: number;
    loading: boolean;
    error: string | null;
    fetchTileMappings: (pageIndex: number, pageSize: number, search: string) => Promise<void>;
    fetchTilesForLocation: (location: { latitude: number; longitude: number }) => Promise<void>;
}

const useTileServiceStore = create<TileServiceState>((set) => ({
    tiles: [],
    tileMappings: [],
    totalTileMappings: 0,
    totalTiles: 0,
    loading: false,
    error: null,

    fetchTileMappings: async (pageIndex: number, pageSize: number, search: string) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/tiles/mappings", {
                params: {
                    page: pageIndex + 1,
                    pageSize,
                    search,
                },
            });
            set({
                tileMappings: response.data?.mappings || [],
                totalTileMappings: response.data?.totalRecords || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching tile mappings:", err);
            set({
                error: "Failed to load tile mapping data.",
                loading: false,
            });
        }
    },

    fetchTilesForLocation: async ({ latitude, longitude }: { latitude: number; longitude: number }) => {
        set({ loading: true, error: null });

        try {
            const response = await apiClient.get("/tiles/all", {
                params: {
                    minLat: latitude - 1,
                    minLon: longitude - 1,
                    maxLat: latitude + 1,
                    maxLon: longitude + 1,
                },
            });

            set({
                tiles: response.data || [],
                totalTiles: response.data?.length || 0,
                loading: false,
            });
        } catch (err) {
            console.error("Error fetching tile data:", err);
            set({
                error: "Failed to fetch tile data. Please try again later.",
                loading: false,
            });
        }
    },
}));

export default useTileServiceStore;
