import * as API from "./api.js";
import * as Notifications from "./notifications.js";
import * as Spinner from "./spinner.js";
import * as Animation from "./animation.js";
import * as Encryption from "./encryption.js";
import * as Duration from "./duration.js";

const CODE_ELEMENT = document.getElementById("code");
const LINE_NUMBERS_ELEMENT = document.getElementById("linenos");
const INPUT_ELEMENT = document.getElementById("input");

const LIFETIME_CONTAINER_ELEMENT = document.getElementById("lifetime_container");

const CHARACTER_AMOUNT_ELEMENT = document.getElementById("characters");
const LINES_AMOUNT_ELEMENT = document.getElementById("lines");

const BUTTONS_DEFAULT_ELEMENT = document.getElementById("buttons_default");
const BUTTON_NEW_ELEMENT = document.getElementById("btn_new");
const BUTTON_SAVE_ELEMENT = document.getElementById("btn_save");
const BUTTON_EDIT_ELEMENT = document.getElementById("btn_edit");
const BUTTON_DELETE_ELEMENT = document.getElementById("btn_delete");
const BUTTON_COPY_ELEMENT = document.getElementById("btn_copy");

const BUTTON_REPORT_ELEMENT = document.getElementById("btn_report");

const BUTTONS_EDIT_ELEMENT = document.getElementById("buttons_edit");
const BUTTON_EDIT_CANCEL_ELEMENT = document.getElementById("btn_edit_cancel");
const BUTTON_EDIT_APPLY_ELEMENT = document.getElementById("btn_edit_apply");

const BUTTON_TOGGLE_ENCRYPTION_ELEMENT = document.getElementById("btn_toggle_encryption");

let PASTE_ID;
let LANGUAGE;
let CODE;

let ENCRYPTION_KEY;
let ENCRYPTION_IV;

let EDIT_MODE = false;

let API_INFORMATION = {
    version: "error",
    pasteLifetime: -1,
    modificationTokens: false,
    reports: false
};

// Initializes the state system
export async function initialize() {
    loadAPIInformation();

    setupButtonFunctionality();
    setupKeybinds();

    // When embedded inside an iframe, add "embedded"
    // class to body element.
    if (window != window.parent) {
        document.body.classList += " embedded";
    }

    // Enable encryption if enabled from last session
    if (localStorage.getItem("encryption") === "true") {
        BUTTON_TOGGLE_ENCRYPTION_ELEMENT.classList.add("active");
    }

    if (location.pathname !== "/") {
        // Extract the paste data (ID and language)
        const split = location.pathname.replace("/", "").split(".");
        const pasteID = split[0];
        const language = split[1];

        // Try to retrieve the paste data from the API
        const response = await API.getPaste(pasteID);
        if (!response.ok) {
            Notifications.error("Could not load paste: <b>" + await response.text() + "</b>");
            setTimeout(() => location.replace(location.protocol + "//" + location.host), 3000);
            return;
        }

        // Set the persistent paste data
        PASTE_ID = pasteID;
        LANGUAGE = language;

        // Decode the response and decrypt the content if needed
        const json = await response.json();
        CODE = json.content;
        if (json.metadata.pf_encryption) {
            ENCRYPTION_KEY = location.hash.replace("#", "");
            while (ENCRYPTION_KEY.length == 0) {
                ENCRYPTION_KEY = prompt("Your decryption key:");
            }

            try {
                 CODE = await Encryption.decrypt(ENCRYPTION_KEY, json.metadata.pf_encryption.iv, CODE);
                 ENCRYPTION_IV = json.metadata.pf_encryption.iv;
            } catch (error) {
                console.log(error);
                Notifications.error("Could not decrypt paste; make sure the decryption key is correct.");
                setTimeout(() => location.replace(location.protocol + "//" + location.host), 3000);
                return;
            }
        }

        // Fill the code block with the just received data
        updateCode();
    } else {
        // Give the user the opportunity to paste his code
        INPUT_ELEMENT.classList.remove("hidden");
        INPUT_ELEMENT.focus();
        LIFETIME_CONTAINER_ELEMENT.classList.remove("hidden");
    }

    // Update the state of the buttons to match the current state
    updateButtonState();

    INPUT_ELEMENT.addEventListener("input", () => {
        updateLineNumbers(INPUT_ELEMENT.value);

        if (BUTTON_SAVE_ELEMENT.hasAttribute("disabled") && INPUT_ELEMENT.value.length > 0) {
            BUTTON_SAVE_ELEMENT.removeAttribute("disabled");
        }
        if (!BUTTON_SAVE_ELEMENT.hasAttribute("disabled") && INPUT_ELEMENT.value.length == 0) {
            BUTTON_SAVE_ELEMENT.setAttribute("disabled", true);
        }
    });
}

// Loads the API information
async function loadAPIInformation() {
    // try to retrieve the API information
    const response = await API.getAPIInformation();
    if (response.ok) {
        API_INFORMATION = await response.json();
    } else {
        Notifications.error("Failed loading API information: <b>" + await response.text() + "</b>");
    }

    // Display the API version
    document.getElementById("version").innerText = API_INFORMATION.version;

    // Display the paste lifetime
    document.getElementById("lifetime").innerText = Duration.format(API_INFORMATION.pasteLifetime);
}

// Sets the current persistent code to the code block, highlights it and updates the line numbers
function updateCode() {
    CODE_ELEMENT.innerHTML = LANGUAGE
        ? hljs.highlight(LANGUAGE, CODE).value
        : hljs.highlightAuto(CODE).value;
    updateLineNumbers(CODE);
}

function updateLineNumbers(content) {
    CHARACTER_AMOUNT_ELEMENT.innerText = content.length;
    LINES_AMOUNT_ELEMENT.innerText = content.split(/\n/).length;

    if (content == "") {
        LINE_NUMBERS_ELEMENT.innerHTML = "<span>></span>";
        return;
    }
    LINE_NUMBERS_ELEMENT.innerHTML = content.split(/\n/).map((_, index) => `<span>${index + 1}</span>`).join("");
}

// Updates the button state according to the current state
function updateButtonState() {
    if (PASTE_ID) {
        BUTTON_SAVE_ELEMENT.setAttribute("disabled", true);
        BUTTON_EDIT_ELEMENT.removeAttribute("disabled");
        BUTTON_DELETE_ELEMENT.removeAttribute("disabled");
        BUTTON_COPY_ELEMENT.removeAttribute("disabled");

        if (API_INFORMATION.reports) {
            BUTTON_REPORT_ELEMENT.classList.remove("hidden");
        }
    } else {
        BUTTON_EDIT_ELEMENT.setAttribute("disabled", true);
        BUTTON_DELETE_ELEMENT.setAttribute("disabled", true);
        BUTTON_COPY_ELEMENT.setAttribute("disabled", true);

        if (API_INFORMATION.reports) {
            BUTTON_REPORT_ELEMENT.classList.add("hidden");
        }
    }
}

// Toggles the edit mode
function toggleEditMode() {
    if (EDIT_MODE) {
        EDIT_MODE = false;
        INPUT_ELEMENT.classList.add("hidden");
        LIFETIME_CONTAINER_ELEMENT.classList.add("hidden");
        CODE_ELEMENT.classList.remove("hidden");
        updateLineNumbers(CODE);
        Animation.animate(BUTTONS_EDIT_ELEMENT, "animate__fadeOutDown", "0.3s", () => {
            BUTTONS_EDIT_ELEMENT.classList.add("hidden");
            BUTTONS_DEFAULT_ELEMENT.classList.remove("hidden");
            Animation.animate(BUTTONS_DEFAULT_ELEMENT, "animate__fadeInDown", "0.3s");
        });
    } else {
        EDIT_MODE = true;
        CODE_ELEMENT.classList.add("hidden");
        LIFETIME_CONTAINER_ELEMENT.classList.remove("hidden");
        INPUT_ELEMENT.classList.remove("hidden");
        INPUT_ELEMENT.value = CODE;
        INPUT_ELEMENT.focus();
        Animation.animate(BUTTONS_DEFAULT_ELEMENT, "animate__fadeOutUp", "0.3s", () => {
            BUTTONS_DEFAULT_ELEMENT.classList.add("hidden");
            BUTTONS_EDIT_ELEMENT.classList.remove("hidden");
            Animation.animate(BUTTONS_EDIT_ELEMENT, "animate__fadeInUp", "0.3s");
        });
    }
}

// Sets up the keybinds for the buttons
function setupKeybinds() {
    window.addEventListener("keydown", (event) => {
        // All keybinds in the default button set include the CTRL key
        if ((EDIT_MODE && !event.ctrlKey && !event.metaKey && event.code !== "Escape") || (!EDIT_MODE && !event.ctrlKey && !event.metaKey)) {
            return;
        }

        // Find the DOM element of the button to trigger
        let element;
        if (EDIT_MODE) {
            switch (event.code) {
                case "Escape": {
                    element = BUTTON_EDIT_CANCEL_ELEMENT;
                    break
                }
                case "KeyS": {
                    element = BUTTON_EDIT_APPLY_ELEMENT;
                    break;
                }
            }
        } else {
            switch (event.code) {
                case "KeyQ": {
                    element = BUTTON_NEW_ELEMENT;
                    break;
                }
                case "KeyS": {
                    element = BUTTON_SAVE_ELEMENT;
                    break;
                }
                case "KeyO": {
                    element = BUTTON_EDIT_ELEMENT;
                    break;
                }
                case "KeyX": {
                    element = BUTTON_DELETE_ELEMENT;
                    break;
                }
                case "KeyB": {
                    element = BUTTON_COPY_ELEMENT;
                    break;
                }
            }
        }

        // Trigger the found button
        if (element) {
            event.preventDefault();
            if (element.hasAttribute("disabled")) {
                return;
            }
            element.click();
        }
    });

    // Additionally fix the behaviour of the Tab key
    window.addEventListener("keydown", (event) => {
        if (event.code != "Tab") {
            return;
        }
        event.preventDefault();

        insertTextAtCursor(inputElement, "    ");
    });
}

// Sets up the different button functionalities
function setupButtonFunctionality() {
    BUTTON_NEW_ELEMENT.addEventListener("click", () => location.replace(location.protocol + "//" + location.host));

    BUTTON_SAVE_ELEMENT.addEventListener("click", () => {
        Spinner.surround(async () => {
            // Only proceed if the input is not empty
            if (!INPUT_ELEMENT.value) {
                return;
            }

            // Encrypt the paste if needed
            let value = INPUT_ELEMENT.value;
            let metadata;
            let key;
            if (BUTTON_TOGGLE_ENCRYPTION_ELEMENT.classList.contains("active")) {
                const encrypted = await Encryption.encrypt(await Encryption.generateEncryptionData(), value);
                value = encrypted.result;
                metadata = {
                    pf_encryption: {
                        alg: "AES-CBC",
                        iv: encrypted.iv
                    }
                };
                key = encrypted.key;
            }

            // Try to create the paste
            const response = await API.createPaste(value, metadata);
            if (!response.ok) {
                Notifications.error("Error while creating paste: <b>" + await response.text() + "</b>");
                return;
            }
            const data = await response.json();

            // Display the modification token if provided
            if (data.modificationToken) {
                prompt("The modification token for your paste is:", data.modificationToken);
            }

            // Redirect the user to his newly created paste
            location.replace(location.protocol + "//" + location.host + "/" + data.id + (key ? "#" + key : ""));
        });
    });
    
    BUTTON_EDIT_ELEMENT.addEventListener("click", toggleEditMode);

    BUTTON_DELETE_ELEMENT.addEventListener("click", () => {
        Spinner.surround(async () => {
            // Ask for the modification token
            const modificationToken = prompt("Modification token:");
            if (!modificationToken) {
                return;
            }

            // Try to delete the paste
            const response = await API.deletePaste(PASTE_ID, modificationToken);
            if (!response.ok) {
                Notifications.error("Error while deleting paste: <b>" + await response.text() + "</b>");
                return;
            }

            // Redirect the user to the start page
            location.replace(location.protocol + "//" + location.host);
        });
    });

    BUTTON_COPY_ELEMENT.addEventListener("click", async () => {
        if (!navigator.clipboard) {
            Notifications.error("Clipboard API not supported by your browser.");
            return;
        }

        await navigator.clipboard.writeText(CODE);
        Notifications.success("Successfully copied the code.");
    });

    BUTTON_EDIT_CANCEL_ELEMENT.addEventListener("click", toggleEditMode);

    BUTTON_EDIT_APPLY_ELEMENT.addEventListener("click", async () => {
        // Only proceed if the input is not empty
        if (!INPUT_ELEMENT.value) {
            return;
        }

        // Ask for the modification token
        const modificationToken = prompt("Modification token:");
        if (!modificationToken) {
            return;
        }

        // Re-encrypt the paste data if needed
        let value = INPUT_ELEMENT.value;
        if (ENCRYPTION_KEY && ENCRYPTION_IV) {
            const encrypted = await Encryption.encrypt(await Encryption.encryptionDataFromHex(ENCRYPTION_KEY, ENCRYPTION_IV), value);
            value = encrypted.result;
        }

        // Try to edit the paste
        const response = await API.editPaste(PASTE_ID, modificationToken, value);
        if (!response.ok) {
            Notifications.error("Error while editing paste: <b>" + await response.text() + "</b>");
            return;
        }

        // Update the code and leave the edit mode
        CODE = INPUT_ELEMENT.value;
        updateCode();
        toggleEditMode();
        Notifications.success("Successfully edited paste.");
    });

    BUTTON_TOGGLE_ENCRYPTION_ELEMENT.addEventListener("click", () => {
        const active = BUTTON_TOGGLE_ENCRYPTION_ELEMENT.classList.toggle("active");
        localStorage.setItem("encryption", active);
        Notifications.success((active ? "Enabled" : "Disabled") + " automatic paste encryption.");
    });

    BUTTON_REPORT_ELEMENT.addEventListener("click", async () => {
        // Ask the user for a reason
        const reason = prompt("Reason:");
        if (!reason) {
            return;
        }

        // Try to report the paste
        const response = await API.reportPaste(PASTE_ID, reason);
        if (!response.ok) {
            Notifications.error("Error while reporting paste: <b>" + await response.text() + "</b>");
            return;
        }

        // Show the response message
        const data = await response.json();
        if (!data.success) {
            Notifications.error("Error while reporting paste: <b>" + data.message + "</b>");
            return;
        }
        Notifications.success(data.message);
    });
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
