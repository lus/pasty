export function format(milliseconds) {
    if (milliseconds < 0) {
        return "forever";
    }

    let parts = new Array();

    let days = Math.floor(milliseconds / 86400000);
    if (days > 0) {
        parts.push(`${days} ${days > 1 ? "days" : "day"}`);
        milliseconds -= days * 86400000;
    }

    let hours = Math.floor(milliseconds / 3600000);
    if (hours > 0) {
        parts.push(`${hours} ${hours > 1 ? "hours" : "hour"}`);
        milliseconds -= hours * 3600000;
    }

    let minutes = Math.floor(milliseconds / 60000);
    if (minutes > 0) {
        parts.push(`${minutes} ${minutes > 1 ? "minutes" : "minute"}`);
        milliseconds -= minutes * 60000;
    }

    let seconds = Math.ceil(milliseconds / 1000);
    if (seconds > 0) {
        parts.push(`${seconds} ${seconds > 1 ? "seconds" : "second"}`);
    }

    return parts.join(", ");
}