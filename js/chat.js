// ═══════════════════════════════════════
// CHAT LOGIC
// ═══════════════════════════════════════

let ws = null;
let currentRecipient = null;
const token = localStorage.getItem('token');
const userId = localStorage.getItem('userId');
const username = localStorage.getItem('username');

// Check authentication
if (!token) {
    window.location.href = 'index.html';
}

// Initialize UI
window.addEventListener('load', () => {
    document.getElementById('currentUsername').textContent = username;
    document.getElementById('userAvatar').textContent = username.charAt(0).toUpperCase();
    connectWebSocket();
});

// ═══════════════════════════════════════
// WEBSOCKET CONNECTION
// ═══════════════════════════════════════
function connectWebSocket() {
    ws = new WebSocket(getWsUrl(token));

    ws.onopen = () => {
        console.log('✅ WebSocket connected');
        updateConnectionStatus('Connected', true);
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        displayMessage(message.content, 'received', message.timestamp);
    };

    ws.onclose = () => {
        console.log('❌ WebSocket disconnected');
        updateConnectionStatus('Disconnected', false);
        
        // Reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        updateConnectionStatus('Error', false);
    };
}

// ═══════════════════════════════════════
// START CHAT WITH RECIPIENT
// ═══════════════════════════════════════
function startChat() {
    const recipientId = document.getElementById('recipientInput').value.trim();
    
    if (!recipientId) {
        alert('Please enter a recipient User ID');
        return;
    }

    currentRecipient = recipientId;
    
    // Update UI
    document.getElementById('chatWith').textContent = 'Chatting with: ' + recipientId.substring(0, 8) + '...';
    document.getElementById('chatAvatar').textContent = recipientId.charAt(0).toUpperCase();
    
    // Enable input
    document.getElementById('messageInput').disabled = false;
    document.getElementById('sendBtn').disabled = false;
    
    // Clear welcome message
    document.getElementById('messages').innerHTML = '';
    
    // Load previous messages
    loadMessages();
}

// ═══════════════════════════════════════
// SEND MESSAGE
// ═══════════════════════════════════════
function sendMessage() {
    const input = document.getElementById('messageInput');
    const content = input.value.trim();

    if (!content || !currentRecipient) {
        return;
    }

    if (ws && ws.readyState === WebSocket.OPEN) {
        const message = {
            type: 'chat_message',
            recipient_id: currentRecipient,
            content: content
        };

        ws.send(JSON.stringify(message));
        displayMessage(content, 'sent', Date.now() / 1000);
        input.value = '';
    } else {
        alert('Not connected. Please wait...');
    }
}

// ═══════════════════════════════════════
// DISPLAY MESSAGE
// ═══════════════════════════════════════
function displayMessage(content, type, timestamp) {
    const messagesDiv = document.getElementById('messages');
    
    const msgDiv = document.createElement('div');
    msgDiv.className = 'msg ' + type;
    
    const time = new Date(timestamp * 1000).toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit'
    });
    
    msgDiv.innerHTML = `
        <div class="msg-content">${escapeHtml(content)}</div>
        <div class="msg-time">${time}</div>
    `;
    
    messagesDiv.appendChild(msgDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// ═══════════════════════════════════════
// LOAD PREVIOUS MESSAGES
// ═══════════════════════════════════════
async function loadMessages() {
    try {
        const response = await fetch(getApiUrl(CONFIG.ENDPOINTS.MESSAGES), {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        const data = await response.json();

        if (data.messages) {
            data.messages.reverse().forEach(msg => {
                const type = msg.sender_id === userId ? 'sent' : 'received';
                const timestamp = new Date(msg.created_at).getTime() / 1000;
                displayMessage(msg.content, type, timestamp);
            });
        }
    } catch (error) {
        console.error('Failed to load messages:', error);
    }
}

// ═══════════════════════════════════════
// UPDATE CONNECTION STATUS
// ═══════════════════════════════════════
function updateConnectionStatus(status, connected) {
    const statusEl = document.getElementById('connectionStatus');
    statusEl.textContent = status;
    statusEl.className = connected ? 'connected' : '';
}

// ═══════════════════════════════════════
// LOGOUT
// ═══════════════════════════════════════
function logout() {
    if (ws) {
        ws.close();
    }
    localStorage.clear();
    window.location.href = 'index.html';
}

// ═══════════════════════════════════════
// UTILITY: Escape HTML (prevent XSS)
// ═══════════════════════════════════════
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ═══════════════════════════════════════
// ENTER KEY TO SEND
// ═══════════════════════════════════════
document.getElementById('messageInput').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});
