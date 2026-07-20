// ═══════════════════════════════════════
// API CONFIGURATION
// ═══════════════════════════════════════

const CONFIG = {
    // CHANGE THIS to your deployed backend URL
    API_URL: 'http://localhost:8080',
    WS_URL: 'ws://localhost:8080',
    
    // For production (Railway):
    // API_URL: 'https://sukaina-chat.up.railway.app',
    // WS_URL: 'wss://sukaina-chat.up.railway.app',
    
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
    return CONFIG.WS_URL + CONFIG.ENDPOINTS.WEBSOCKET + '?token=' + token;
}
