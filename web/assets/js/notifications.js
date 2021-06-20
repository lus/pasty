// element holds the notification containers DOM element
const element = document.getElementById("notifications");

// error shows an error notifications
export function error(message) {
    create("error", message, 3000);
}

// success shows a success notifications
export function success(message) {
    create("success", message, 3000);
}

// create creates a new notification
function create(type, message, duration) {
    const node = document.createElement("div");
    node.classList.add(type);
    node.innerHTML = message;

    element.appendChild(node);
    setTimeout(() => element.removeChild(node), duration);
}