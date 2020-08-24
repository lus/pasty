// Import the used modules
import * as api from "./api.js";
import * as autoload from "./autoload.js";
import * as spinner from "./spinner.js";

// setupKeybinds initializes the keybinds for the buttons
export function setupKeybinds() {
    window.onkeydown = function(event) {
        // Return if the CTRL key was not pressed
        if (!event.ctrlKey) return;
    
        // Define the DOM element of the pressed button
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
    
        // Call the onClick function of the button
        if (element) {
            if (element.hasAttribute("disabled")) return;
            event.preventDefault();
            element.click();
        }
    }
}

// setupButtons configures the click listeners of the buttons
export function setupButtons() {
    // Define the behavior of the 'new' button
    document.getElementById("btn_new").addEventListener("click", function() {
        location.replace(location.protocol + "//" + location.host);
    });

    // Define the behavior of the 'save' button
    document.getElementById("btn_save").addEventListener("click", function() {
        spinner.surround(async function() {
            // Return if the text area is empty
            const input = document.getElementById("input");
            if (!input.value) return;

            // Create the paste
            const response = await api.createPaste(input.value);
            if (!response.ok) {
                alert("Error:\n\n" + data);
                return;
            }
            const data = await response.json();

            // Redirect the user to the paste page
            let address = location.protocol + "//" + location.host + "/" + data.id;
            if (data.suggestedSyntaxType) address += "." + data.suggestedSyntaxType;
            location.replace(address);

            // TODO: Find a solution to display the deletion token
        });
    });

    // Define the behavior of the 'delete' button
    document.getElementById("btn_delete").addEventListener("click", function() {
        spinner.surround(async function() {
            // Ask the user for the deletion token
            const deletionToken = window.prompt("Deletion Token:");
            if (!deletionToken) return;

            // Delete the paste
            const response = await api.deletePaste(autoload.PASTE_ID, deletionToken);
            const data = await response.text();
            if (!response.ok) {
                alert("Error:\n\n" + data);
                return;
            }

            // Redirect the user to the main page
            location.replace(location.protocol + "//" + location.host);
        });
    });

    // Define the behavior of the 'copy' button
    document.getElementById("btn_copy").addEventListener("click", function() {
        spinner.surround(async function() {
            // Ask for the clipboard permissions
            askClipboardPermissions();
            
            // Copy the code
            await navigator.clipboard.writeText(document.getElementById("code").innerText);
        });
    });
}

// askClipboardPermissions asks the user for the clipboard permissions
async function askClipboardPermissions() {
    try {
        const state = await navigator.permissions.query({
            name: "clipbaord-write"
        });
        return state === "granted";
    } catch (error) {
        return false;
    }
}