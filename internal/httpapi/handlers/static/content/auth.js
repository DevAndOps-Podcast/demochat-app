import { initChat } from './app.js';

document.addEventListener('DOMContentLoaded', () => {
  const loginContainer = document.getElementById('login-container');
  const chatContainer = document.getElementById('chat-container');
  const loginButton = document.getElementById('login-button');
  const usernameInput = document.getElementById('username-input');
  const passwordInput = document.getElementById('password-input');
  const loginError = document.getElementById('login-error');

  const registerForm = document.getElementById('registerForm');

  // Check if a token exists (simple check for logged-in state)
  const accessToken = localStorage.getItem('accessToken');
  if (accessToken) {
    if (loginContainer) loginContainer.style.display = 'none';
    if (chatContainer) chatContainer.style.display = 'flex';
    initChat();
  } else {
    if (loginContainer) loginContainer.style.display = 'flex';
    if (chatContainer) chatContainer.style.display = 'none';
  }

  if (loginButton) {
    loginButton.addEventListener('click', async () => {
      const username = usernameInput.value;
      const password = passwordInput.value;

      try {
        const response = await fetch('/auth', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ username, password })
        });

        const data = await response.json();

        if (response.ok) {

          localStorage.setItem('accessToken', data.access_token);
          localStorage.setItem('refreshToken', data.refresh_token);
          if (loginContainer) loginContainer.style.display = 'none';
          if (chatContainer) chatContainer.style.display = 'flex';
          if (loginError) loginError.style.display = 'none';
          initChat();
        } else {
          if (loginError) {
            loginError.textContent = data.message || 'Authentication failed.';
            loginError.style.display = 'block';
          }
        }
      } catch (error) {
        if (loginError) {
          loginError.textContent = 'An error occurred during login.';
          loginError.style.display = 'block';
        }
        console.error('Login error:', error);
      }
    });
  }

  if (registerForm) {
    registerForm.addEventListener('submit', async (event) => {
      event.preventDefault();
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      try {
        const response = await fetch('/auth/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ username, password })
        });

        if (response.ok) {
          alert('Registration successful! You can now log in.');
          window.location.href = 'index.html';
        } else {
          const data = await response.json();
          alert(data.message || 'Registration failed.');
        }
      } catch (error) {
        alert('An error occurred during registration.');
        console.error('Registration error:', error);
      }
    });
  }
});
