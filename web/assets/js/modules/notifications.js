import * as Animation from "./animation.js";

const ELEMENT = document.getElementById("notifications");

// Shows a success notification
export function success(message) {
    create("success", message, 3000);
}

// Shows an error notification
export function error(message) {
    create("error", message, 3000);
}

// Creates a new custom notification
function create(type, message, duration) {
    const node = document.createElement("div");
    node.classList.add(type);
    Animation.animate(node, "animate__fadeInUp", "0.2s");
    node.innerHTML = message;

    ELEMENT.childNodes.forEach(child => Animation.animate(child, "animate__slideInUp", "0.2s"));
    ELEMENT.appendChild(node);
    setTimeout(() => Animation.animate(node, "animate__fadeOutUp", "0.2s", () => ELEMENT.removeChild(node)), duration);
} 