// Import the used modules
import * as api from "./api.js";
import * as buttons from "./buttons.js";
import * as spinner from "./spinner.js";
import * as notifications from "./notifications.js";

// Set up the buttons
buttons.setupButtons();
buttons.setupKeybinds();

// Define element handles
const versionElement = document.getElementById("version");
const lineNOsElement = document.getElementById("linenos");
const codeElement = document.getElementById("code");
const inputElement = document.getElementById("input");

// Load the API information
async function loadAPIInformation() {
    const response = await api.getAPIInformation();
    if (!response.ok) {
        const data = await response.text();
        notifications.error("Failed fetching the API information: <b>" + data + "</b>");
        return;
    }
    const data = await response.json();
    versionElement.innerText = data.version;
}
loadAPIInformation();

// Try to load a paste if one exists
export let PASTE_ID;
let CODE;
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
        CODE = (await response.json()).content;

        // Adjust the button states
        document.getElementById("btn_save").setAttribute("disabled", true);
        document.getElementById("btn_delete").removeAttribute("disabled");
        document.getElementById("btn_copy").removeAttribute("disabled");

        // Set the paste content to the DOM
        codeElement.innerHTML = language
            ? hljs.highlight(language, CODE).value
            : hljs.highlightAuto(CODE).value;

        // Display the line numbers
        renderLineNumbers();
        window.addEventListener("resize", renderLineNumbers);

        // Set the PASTE_ID variable
        PASTE_ID = pasteID;
    } else {
        inputElement.classList.remove("hidden");
        inputElement.focus();
        window.addEventListener("keydown", function (event) {
            if (event.keyCode != 9) return;
            event.preventDefault();

            insertTextAtCursor(inputElement, "    ");
        });
    }
}
spinner.surround(loadPaste);

function renderLineNumbers() {
    lineNOsElement.innerHTML = CODE.split(/\n/).map((line, index) => {
        let lineWidth = getTextWidth(line, "16px Source Code Pro");
        let linesSpace = Math.ceil(lineWidth / codeElement.offsetWidth);

        let result = `<span>${index+1}</span>`;
        if (linesSpace > 1) {
            result += "<span></span>".repeat(linesSpace - 1);
        }
        return result;
    }).join("");
}

// 1:1 skid from https://stackoverflow.com/questions/7404366/how-do-i-insert-some-text-where-the-cursor-is
function insertTextAtCursor(element, text) {
    let value = element.value, endIndex, range, doc = element.ownerDocument;
    if (typeof element.selectionStart == "number"
        && typeof element.selectionEnd == "number") {
        endIndex = element.selectionEnd;
        element.value = value.slice(0, endIndex) + text + value.slice(endIndex);
        element.selectionStart = element.selectionEnd = endIndex + text.length;
    } else if (doc.selection != "undefined" && doc.selection.createRange) {
        element.focus();
        range = doc.selection.createRange();
        range.collapse(false);
        range.text = text;
        range.select();
    }
}

// Also a kind of skid
function getTextWidth(text, font) {
    let canvas = getTextWidth.canvas || (getTextWidth.canvas = document.createElement("canvas"));
    let context = canvas.getContext("2d");
    context.font = font;
    return context.measureText(text).width;
}