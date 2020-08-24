// getAPIInformation returns the API information
export async function getAPIInformation() {
    return await fetch(location.protocol + "//" + location.host + "/api/v1/info");
}

// getPaste retrieves a paste
export async function getPaste(id) {
    return await fetch(location.protocol + "//" + location.host + "/api/v1/pastes/" + id);
}

// createPaste creates a new paste
export async function createPaste(content) {
    return await fetch(location.protocol + "//" + location.host + "/api/v1/pastes", {
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
    return await fetch(location.protocol + "//" + location.host + "/api/v1/pastes/" + id, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            deletionToken
        })
    });
}