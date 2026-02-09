// --- 1. Crypto Helpers (Mock implementation using Base64) ---
const PREFIX = "||ENC||";

function encryptText(text) {
    // In real life, use PGP or AES here
    return PREFIX + btoa(text);
}

function decryptText(encryptedText) {
    if (!encryptedText.startsWith(PREFIX)) return null;
    try {
        const raw = encryptedText.substring(PREFIX.length);
        return atob(raw);
    } catch (e) {
        return "Error decrypting";
    }
}

// --- 2. DOM Manipulation Logic ---

function addEncryptButton(composeWindow) {
    // Prevent adding duplicate buttons
    if (composeWindow.querySelector('.my-secure-send-btn')) return;

    // Locate the toolbar where the "Send" button lives
    // Gmail class names (like .gU.Up) are obfuscated and change.
    // We look for the row containing the 'Send' button.
    const sendButton = composeWindow.querySelector('div[role="button"][data-tooltip^="Send"]');

    if (sendButton) {
        const parentRow = sendButton.closest('tr') || sendButton.parentElement;

        const secureBtn = document.createElement('div');
        secureBtn.className = 'my-secure-send-btn T-I J-J5-Ji aoO v7 T-I-atl L3'; // Reuse Gmail button styles
        secureBtn.innerHTML = '🔒 Secure Send';
        secureBtn.style.marginRight = '10px';
        secureBtn.style.backgroundColor = '#d93025'; // Red color
        secureBtn.style.color = 'white';
        secureBtn.style.cursor = 'pointer';

        secureBtn.onclick = function() {
            // 1. Find the message body area
            const messageBody = composeWindow.querySelector('div[aria-label="Message Body"], div[aria-label="Message text"]');

            if (messageBody) {
                // 2. Encrypt content
                const originalText = messageBody.innerText;
                const encrypted = encryptText(originalText);

                // 3. Update content
                messageBody.innerText = encrypted;

                // 4. Trigger real send
                // Wait a split second for Angular/Gmail to recognize the text change
                setTimeout(() => {
                    sendButton.click();
                }, 200);
            }
        };

        // Insert before the send button container
        sendButton.parentElement.insertBefore(secureBtn, sendButton);
    }
}

function checkAndDecryptMessage(messageElement) {
    // Check if we already processed this specific message div
    if (messageElement.getAttribute('data-decrypted') === 'true') return;

    const content = messageElement.innerText;

    if (content.includes(PREFIX)) {
        // Simple logic: check if the whole body is the encrypted string
        // You might need regex to find the string inside standard email wrappers
        const decrypted = decryptText(content.trim());

        if (decrypted) {
            messageElement.innerHTML = `<div style="border: 2px solid green; padding: 10px; background: #e8f5e9;">
                <strong>🔓 Decrypted Message:</strong><br/>
                ${decrypted.replace(/\n/g, '<br>')}
            </div>`;
            messageElement.setAttribute('data-decrypted', 'true');
        }
    }
}

// --- 3. Mutation Observer (The Engine) ---

const observer = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
        for (const node of mutation.addedNodes) {
            if (node.nodeType !== 1) continue; // Element nodes only

            // A. DETECT COMPOSE WINDOW
            // Look for the "Compose" dialog or main window
            if (node.matches('div[role="dialog"]') || node.querySelector('div[role="dialog"]')) {
                // It might be a reply or a new compose
                const compose = node.matches('div[role="dialog"]') ? node : node.querySelector('div[role="dialog"]');
                if (compose) addEncryptButton(compose);
            }

            // Sometimes the button is added to an existing view (like hitting Reply)
            // So we also check generically for the toolbar presence
            const buttons = document.querySelectorAll('div[role="button"][data-tooltip^="Send"]');
            buttons.forEach(btn => {
                // Find the compose container for this specific button
                const composeContainer = btn.closest('div[role="dialog"]') || btn.closest('td.I5'); // I5 is often the reply area
                if (composeContainer) {
                    addEncryptButton(composeContainer);
                }
            });

            // B. DETECT OPENED MESSAGE
            // '.a3s' is a very common, stable-ish class for message body content in Gmail
            if (node.classList.contains('a3s') || node.querySelector('.a3s')) {
                const messages = node.classList.contains('a3s') ? [node] : node.querySelectorAll('.a3s');
                messages.forEach(checkAndDecryptMessage);
            }
        }
    }
});

// Start observing the entire body for changes
observer.observe(document.body, {
    childList: true,
    subtree: true
});
