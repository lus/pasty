// Load the API information
loadAPIInfo();

// Set up the keybinds
setupKeybinds();

// Try to load a paste if one exists
let PASTE_ID = "";
function loadPaste() {
    let split = location.pathname.split(".");
    let pasteID = split[0];
    let language = split[1];
    getPaste(pasteID, function(success, data) {
        // Return if no paste was found
        if (!success) {
            location.replace(location.protocol + "//" + location.host);
            return;
        };

        // Enable and disable the corresponding buttons
        document.getElementById("btn_save").setAttribute("disabled", true);
        document.getElementById("btn_delete").removeAttribute("disabled");
        document.getElementById("btn_copy").removeAttribute("disabled");

        // Set the paste content to the DOM and display the line numbers
        document.getElementById("code").innerHTML = language
            ? hljs.highlight(language, data.content).value.replace("\n", "<br />")
            : hljs.highlightAuto(data.content).value.replace("\n", "<br />");
        for (i = 1; i <= data.content.split(/\n/).length; i++) {
            document.getElementById("linenos").innerHTML += "<span>" + i + "</span>";
        }

        // Set the PASTE_ID variable
        PASTE_ID = pasteID;
    });
}
if (location.pathname != "/") {
    loadPaste();
} else {
    const element = document.getElementById("input");
    element.classList.remove("hidden");
    element.focus();
}

// Define a function to copy text to the clipboard
function copyToClipboard(text) {
    const element = document.createElement("textarea");
    element.value = text;
    document.body.appendChild(element);
    element.select();
    document.execCommand("copy");
    document.body.removeChild(element);
}