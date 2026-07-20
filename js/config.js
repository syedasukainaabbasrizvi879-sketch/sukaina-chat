// ═══════════════════════════════════════
// API CONFIGURATION
// ═══════════════════════════════════════

const CONFIG = {
    // UPDATED: Deployed backend URLs
    API_URL: 'https://sukaina-chat.onrender.com',
    WS_URL: 'wss://sukaina-chat.onrender.com',
    
    ENDPOINTS: {
        REGISTER: '/api/v1/auth/register',
        LOGIN: '/api/v1/auth/login',
        MESSAGES: '/api/v1/messages',
        WEBSOCKET: '/ws'
    }
};

// Helper function to get full API URL
function getApiUrl(endpoint) {
    return CONFIG.API_URL + endpoint;
}

// Helper function to get WebSocket URL
function getWsUrl(token) {
    // HTTPS ke liye WSS use hota hai
    return CONFIG.WS_URL + CONFIG.ENDPOINTS.WEBSOCKET + '?token=' + token;
}
