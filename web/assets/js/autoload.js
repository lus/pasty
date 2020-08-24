// Import the used modules
import * as api from "./api.js";
import * as buttons from "./buttons.js";
import * as spinner from "./spinner.js";

// Set up the buttons
buttons.setupButtons();
buttons.setupKeybinds();

// Load the API information
async function loadAPIInformation() {
    const response = await api.getAPIInformation();
    if (response.ok) {
        const data = await response.json();
        document.getElementById("version").innerText = data.version;
    }
}
loadAPIInformation();

// Try to load a paste if one exists
export let PASTE_ID;
async function loadPaste() {
    if (location.pathname != "/") {
        // Define the paste ID and language
        const split = location.pathname.replace("/", "").split(".");
        const pasteID = split[0];
        const language = split[1];
    
        // Retrieve the paste from the API and redirect the user to the main page if it could not be found
        const response = await api.getPaste(pasteID);
        if (!response.ok) {
            location.replace(location.protocol + "//" + location.host);
            return;
        }
        const data = await response.json();
    
        // Adjust the button states
        document.getElementById("btn_save").setAttribute("disabled", true);
        document.getElementById("btn_delete").removeAttribute("disabled");
        document.getElementById("btn_copy").removeAttribute("disabled");
    
        // Set the paste content to the DOM
        document.getElementById("code").innerHTML = language
            ? hljs.highlight(language, data.content).value
            : hljs.highlightAuto(data.content).value;
        
        // Display the line numbers
        const lineNumbersElement = document.getElementById("linenos");
        data.content.split(/\n/).forEach(function(_currentValue, index) {
            lineNumbersElement.innerHTML += "<span>" + (index + 1) + "</span>";
        });
    
        // Set the PASTE_ID variable
        PASTE_ID = pasteID;
    } else {
        const input = document.getElementById("input");
        input.classList.remove("hidden");
        input.focus();
    }
}
spinner.surround(loadPaste);