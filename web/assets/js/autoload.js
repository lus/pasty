// Import the used modules
import * as api from "./api.js";
import * as buttons from "./buttons.js";
import * as spinner from "./spinner.js";
import * as notifications from "./notifications.js";

// Set up the buttons
buttons.setupButtons();
buttons.setupKeybinds();

// Load the API information
async function loadAPIInformation() {
    const response = await api.getAPIInformation();
    if (!response.ok) {
        const data = await response.text();
        notifications.error("Failed fetching the API information: <b>" + data + "</b>");
        return;
    }
    const data = await response.json();
    document.getElementById("version").innerText = data.version;
}
loadAPIInformation();

// Try to load a paste if one exists
export let PASTE_ID;
async function loadPaste() {
    if (location.pathname !== "/") {
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
        const code = document.getElementById("code");
        code.innerHTML = language
            ? hljs.highlight(language, data.content).value
            : hljs.highlightAuto(data.content).value;

        // Display the line numbers
        const lineNOs = document.getElementById("linenos");
        lineNOs.innerHTML = data.content.split(/\n/).map((_, index) => `<span>${index + 1}</span>`).join('');

        // Set the PASTE_ID variable
        PASTE_ID = pasteID;
    } else {
        const input = document.getElementById("input");
        input.classList.remove("hidden");
        input.focus();
        window.addEventListener("keydown", function (event) {
            if (event.keyCode != 9) return;
            event.preventDefault();

            insertTextAtCursor(input, "    ");
        });
    }
}
spinner.surround(loadPaste);

// 1:1 skid from https://stackoverflow.com/questions/7404366/how-do-i-insert-some-text-where-the-cursor-is
function insertTextAtCursor(el, text) {
    var val = el.value, endIndex, range, doc = el.ownerDocument;
    if (typeof el.selectionStart == "number"
        && typeof el.selectionEnd == "number") {
        endIndex = el.selectionEnd;
        el.value = val.slice(0, endIndex) + text + val.slice(endIndex);
        el.selectionStart = el.selectionEnd = endIndex + text.length;
    } else if (doc.selection != "undefined" && doc.selection.createRange) {
        el.focus();
        range = doc.selection.createRange();
        range.collapse(false);
        range.text = text;
        range.select();
    }
}