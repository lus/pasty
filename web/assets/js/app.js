import * as API from "./modules/api.js";
import * as Notifications from "./modules/notifications.js";
import * as Spinner from "./modules/spinner.js";
import * as State from "./modules/state.js";

// Load the API information
let API_INFORMATION = {
    version: "error",
    modificationTokens: false,
    reports: false
};
const response = await API.getAPIInformation();
if (response.ok) {
    API_INFORMATION = await response.json();
} else {
    Notifications.error("Failed loading API information: <b>" + await response.text() + "</b>");
}

// Display the API version
document.getElementById("version").innerText = API_INFORMATION.version;

// Initialize the application state
Spinner.surround(State.initialize);
