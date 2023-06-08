import * as Animation from "./animation.js";

const ELEMENT = document.getElementById("spinner-container");

// SHows the spinner
export function show() {
    ELEMENT.classList.remove("hidden");
    Animation.animate(ELEMENT, "animate__zoomIn", "0.2s");
}

// Hides the spinner
export function hide() {
    Animation.animate(ELEMENT, "animate__zoomOut", "0.2s", () => ELEMENT.classList.add("hidden"));
}

// Surrounds an async action with a spinner
export async function surround(innerFunction) {
    show();
    await innerFunction();
    hide();
}