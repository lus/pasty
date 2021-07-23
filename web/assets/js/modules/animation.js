// Properly animates an element
export function animate(element, animation, duration, after) {
    element.style.setProperty("--animate-duration", duration);
    element.classList.add("animate__animated", animation);
    element.addEventListener("animationend", () => {
        element.style.removeProperty("--animate-duration");
        element.classList.remove("animate__animated", animation);
        if (after) {
            after();
        }
    }, {once: true});
}