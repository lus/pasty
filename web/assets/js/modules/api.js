const API_BASE_URL = location.protocol + "//" + location.host + "/api/v2";

export async function getAPIInformation() {
    return fetch(API_BASE_URL + "/info");
}

export async function getPaste(pasteID) {
    return fetch(API_BASE_URL + "/pastes/" + pasteID);
}

export async function createPaste(content, metadata) {
    return fetch(API_BASE_URL + "/pastes", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            content,
            metadata
        })
    });
}

export async function editPaste(pasteID, modificationToken, content, metadata) {
    return fetch(API_BASE_URL + "/pastes/" + pasteID, {
        method: "PATCH",
        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + modificationToken,
        },
        body: JSON.stringify({
            content,
            metadata
        })
    });
}

export async function deletePaste(pasteID, modificationToken) {
    return fetch(API_BASE_URL + "/pastes/" + pasteID, {
        method: "DELETE",
        headers: {
            "Authorization": "Bearer " + modificationToken,
        }
    });
}

export async function reportPaste(pasteID, reason) {
    return fetch(API_BASE_URL + "/pastes/" + pasteID + "/report", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            reason
        })
    });
}
