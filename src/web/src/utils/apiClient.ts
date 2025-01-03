import axios from "axios";

const apiClient = axios.create({
    baseURL: process.env.API_BASE_URL || "http://localhost:8081", // Fallback for local development
    headers: { Accept: "application/json" },
});

apiClient.interceptors.request.use((config) => {
    return config;
});

apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        return Promise.reject(error);
    }
);

export default apiClient;
