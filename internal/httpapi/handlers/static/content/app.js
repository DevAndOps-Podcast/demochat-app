export function initChat() {
  const messagesDiv = document.getElementById("messages");
  const messageInput = document.getElementById("message-input");
  const sendButton = document.getElementById("send-button");
  const insightsContainer = document.getElementById("insights-container");

  async function fetchAndDisplayMessages() {
    try {
      const response = await fetch("/messages", {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${localStorage.getItem('accessToken')}`
        }
      });

      if (response.ok) {
        const messages = await response.json();
        messagesDiv.innerHTML = ''; // Clear existing messages
        messages.forEach(msg => {
          const messageElement = document.createElement("div");
          messageElement.classList.add("chat-message");
          const senderSpan = document.createElement("span");
          senderSpan.classList.add("message-sender");
          senderSpan.textContent = `${msg.username}: `;
          const contentSpan = document.createElement("span");
          contentSpan.classList.add("message-content");
          contentSpan.textContent = msg.message;
          messageElement.appendChild(senderSpan);
          messageElement.appendChild(contentSpan);
          messagesDiv.appendChild(messageElement);
        });
      } else {
        console.error("Failed to fetch messages");
        if (response.status === 401) {
          alert("Session expired. Please log in again.");
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          window.location.reload();
        }
      }
    } catch (error) {
      console.error("Error fetching messages:", error);
    }
  }

  async function fetchAndDisplayInsights() {
    try {
      const response = await fetch("/insights", {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${localStorage.getItem('accessToken')}`
        }
      });

      if (response.ok) {
        const insights = await response.json();
        insightsContainer.innerHTML = `
          <strong>Most Active User:</strong> ${insights.most_active_user} | 
          <strong>Total Messages:</strong> ${insights.total_messages} | 
          <strong>Average Message Rate:</strong> ${typeof insights.average_message_rate === 'number' ? insights.average_message_rate.toFixed(2) : '0.00'} messages/minute
        `;
      } else {
        console.error("Failed to fetch insights");
      }
    } catch (error) {
      console.error("Error fetching insights:", error);
    }
  }

  fetchAndDisplayMessages(); // Fetch and display messages on chat initialization
  fetchAndDisplayInsights(); // Fetch and display insights on chat initialization

  setInterval(fetchAndDisplayInsights, 5000); // Update insights every 5 seconds

  sendButton.addEventListener("click", async () => {
    const message = messageInput.value;
    if (message.trim() === "") {
      return;
    }

    try {
      const response = await fetch("/messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Authorization": `Bearer ${localStorage.getItem('accessToken')}`
        },
        body: JSON.stringify({ message }),
      });

      if (response.ok) {
        messageInput.value = "";
        // Instead of just appending, re-fetch to get the latest state including other users' messages
        fetchAndDisplayMessages();
      } else {
        console.error("Failed to send message");
        if (response.status === 401) {
          alert("Session expired. Please log in again.");
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
          window.location.reload();
        }
      }
    } catch (error) {
      console.error("Error sending message:", error);
    }
  });
}
