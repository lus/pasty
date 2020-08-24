// element holds the spinners DOM element
const element = document.getElementById("spinner");

// show shows the spinner
export function show() {
    element.classList.remove("hidden");
}

// hide hides the spinner
export function hide() {
    element.classList.add("hidden");
}

// surround surrounds an action with a spinner
export async function surround(action) {
    show();
    await action();
    hide();
}