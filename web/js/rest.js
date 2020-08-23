// loadAPIInfo loads and displays the API information
function loadAPIInfo() {
    fetch(location.protocol + "//" + location.host + "/api/v1/info")
        .then(response => response.json())
        .then(data => document.getElementById("version").innerText = data.version);
}

// getPaste retrieves a paste
function getPaste(id, callback) {
    fetch(location.protocol + "//" + location.host + "/api/v1/pastes/" + id)
        .then(response => {
            if (response.status != 200) {
                response.text().then(data => callback(false, data));
                return;
            }
            response.json().then(data => callback(true, data));
        });
}

// createPaste creates a new paste
function createPaste(content, callback) {
    fetch(location.protocol + "//" + location.host + "/api/v1/pastes", {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            content: content
        })
    }).then(response => {
        if (response.status != 200) {
            response.text().then(data => callback(false, data));
            return;
        }
        response.json().then(data => callback(true, data));
    });
}

// deletePaste deletes a paste
function deletePaste(id, deletionToken, callback) {
    fetch(location.protocol + "//" + location.host + "/api/v1/pastes/" + id, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            deletionToken: deletionToken
        })
    }).then(response => {
        if (response.status != 200) {
            response.text().then(data => callback(false, data));
            return;
        }
        response.text().then(data => callback(true, data));
    });
}