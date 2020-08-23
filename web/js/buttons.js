// setupKeybinds initializes the keybinds for the buttons
function setupKeybinds() {
    window.onkeydown = function(event) {
        if (!event.ctrlKey) return;

        let element = null;
        switch (event.keyCode) {
            case 81: {
                element = document.getElementById("btn_new");
                break;
            }
            case 83: {
                element = document.getElementById("btn_save");
                break;
            }
            case 88: {
                element = document.getElementById("btn_delete");
                break;
            }
            case 67: {
                element = document.getElementById("btn_copy");
                break;
            }
        }

        if (element) {
            if (element.hasAttribute("disabled")) return;
            event.preventDefault();
            element.onclick();
        }
    }
}

// Define the behavior of the 'new' button
document.getElementById("btn_new").onclick = function() {
    location.replace(location.protocol + "//" + location.host);
}

// Define the behavior of the 'save' button
document.getElementById("btn_save").onclick = function() {
    // Return if the text area is empty
    if (!document.getElementById("input").value) return;

    // Create the paste
    createPaste(document.getElementById("input").value, function(success, data) {
        // Notify the user about an error if one occurs
        if (!success) {
            alert("Error:\n\n" + data);
            return;
        }

        // Redirect the user to the paste page
        let address = location.protocol + "//" + location.host + "/" + data.id;
        if (data.suggestedSyntaxType) address += "." + data.suggestedSyntaxType;
        copyToClipboard(data.deletionToken);
        location.replace(address);
    });
}

// Define the behavior of the 'delete' button
document.getElementById("btn_delete").onclick = function() {
    // Ask the user for the deletion token
    let deletionToken = window.prompt("Deletion Token:");
    if (!deletionToken) return;

    // Delete the paste
    deletePaste(PASTE_ID, deletionToken, function(success, data) {
        // Notify the user about an error if one occurs
        if (!success) {
            alert("Error:\n\n" + data);
            return;
        }

        // Redirect the user to the default page
        location.replace(location.protocol + "//" + location.host);
    });
}

// Define the behavior of the 'copy' button
document.getElementById("btn_copy").onclick = function() {
    copyToClipboard(document.getElementById("code").innerText);
}