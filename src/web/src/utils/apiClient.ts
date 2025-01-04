import axios from "axios";

// Create Axios instance
const apiClient = axios.create({
    baseURL: process.env.API_BASE_URL || "http://localhost:8081", // Fallback for local development
    headers: { Accept: "application/json", "Content-Type": "application/json" },
    timeout: 30000, // Timeout in milliseconds
});

// Request interceptor
apiClient.interceptors.request.use(
    (config) => {
        // Attach Authorization token if available
        const token = localStorage.getItem("authToken");
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        console.error("Request setup error:", error);
        return Promise.reject(error);
    }
);

// Response interceptor
apiClient.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    console.warn("Unauthorized! Redirecting to login...");
                    break;
                case 406:
                    console.error("HTTP 406: Not Acceptable. Check headers or payload.");
                    break;
                case 500:
                    console.error("HTTP 500: Internal Server Error:", error.response.data.message || "No details provided.");
                    break;
                default:
                    console.warn(`HTTP ${error.response.status}: ${error.response.data.message || "No details provided."}`);
            }
        } else if (error.request) {
            console.error("Network error or no response received.");
        } else {
            console.error("Request setup error:", error.message);
        }
        return Promise.reject(error);
    }
);

// Debugging logs in development mode
if (process.env.NODE_ENV === "development") {
    apiClient.interceptors.request.use((config) => {
        console.log("Request Config:", config);
        return config;
    });
    apiClient.interceptors.response.use(
        (response) => {
            console.log("Response Data:", response.data);
            return response;
        },
        (error) => {
            console.error("Error Response:", error.response || error);
            return Promise.reject(error);
        }
    );
}

export default apiClient;
