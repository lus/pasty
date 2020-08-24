// apiBase defines the base URL of the API
const apiBase = location.protocol + "//" + location.host + "/api/v1";

// getAPIInformation returns the API information
export async function getAPIInformation() {
    return fetch(apiBase + "/info");
}

// getPaste retrieves a paste
export async function getPaste(id) {
    return fetch(apiBase + "/pastes/" + id);
}

// createPaste creates a new paste
export async function createPaste(content) {
    return await fetch(apiBase + "/pastes", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            content
        })
    });
}

// deletePaste deletes a paste
export async function deletePaste(id, deletionToken) {
    return await fetch(apiBase + "/pastes/" + id, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            deletionToken
        })
    });
}