// ═══════════════════════════════════════
// AUTHENTICATION LOGIC
// ═══════════════════════════════════════

// Switch between login and register tabs
function switchTab(tab) {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const tabBtns = document.querySelectorAll('.tab-btn');

    tabBtns.forEach(btn => btn.classList.remove('active'));

    if (tab === 'login') {
        loginForm.classList.add('active');
        registerForm.classList.remove('active');
        tabBtns[0].classList.add('active');
    } else {
        loginForm.classList.remove('active');
        registerForm.classList.add('active');
        tabBtns[1].classList.add('active');
    }
}

// Show message to user
function showMessage(text, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = text;
    messageDiv.className = 'message ' + type;

    setTimeout(() => {
        messageDiv.className = 'message';
    }, 5000);
}

// ═══════════════════════════════════════
// LOGIN HANDLER
// ═══════════════════════════════════════
document.getElementById('loginForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const response = await fetch(getApiUrl(CONFIG.ENDPOINTS.LOGIN), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok && data.token) {
            // Save token and user info
            localStorage.setItem('token', data.token);
            localStorage.setItem('userId', data.user_id);
            localStorage.setItem('username', data.username);

            showMessage('Login successful! Redirecting...', 'success');

            // Redirect to chat
            setTimeout(() => {
                window.location.href = 'chat.html';
            }, 1000);
        } else {
            showMessage(data.error || 'Login failed', 'error');
        }
    } catch (error) {
        showMessage('Connection error. Is backend running?', 'error');
        console.error('Login error:', error);
    }
});

// ═══════════════════════════════════════
// REGISTER HANDLER
// ═══════════════════════════════════════
document.getElementById('registerForm').addEventListener('submit', async (e) => {
    e.preventDefault();

    const username = document.getElementById('registerUsername').value;
    const email = document.getElementById('registerEmail').value;
    const password = document.getElementById('registerPassword').value;

    try {
        const response = await fetch(getApiUrl(CONFIG.ENDPOINTS.REGISTER), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Registration successful! Please login.', 'success');
            
            // Switch to login tab
            setTimeout(() => {
                switchTab('login');
                document.getElementById('loginUsername').value = username;
            }, 1500);
        } else {
            showMessage(data.error || 'Registration failed', 'error');
        }
    } catch (error) {
        showMessage('Connection error. Is backend running?', 'error');
        console.error('Register error:', error);
    }
});

// Check if already logged in
window.addEventListener('load', () => {
    const token = localStorage.getItem('token');
    if (token) {
        window.location.href = 'chat.html';
    }
});
