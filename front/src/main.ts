import { cloudInit } from "./sirius/objects/weather/clouds";
import { rainInit } from "./sirius/objects/weather/rain";
import { Sirius } from "./sirius/sirius";
import { CafeSocket } from "./socket/socket";

window.tableRect = {
    x: 0,
    y: 0,
    width: 0,
    height: 0
}

function updateTableRectangle() {
    const table: HTMLElement = document.getElementsByTagName("tbody")[0];
    if (!table) { return; }
    let boundingRect: DOMRect = table.getBoundingClientRect();
    window.tableRect.x = boundingRect.x;
    window.tableRect.y = boundingRect.y + 6;
    window.tableRect.width = boundingRect.width;
    window.tableRect.height = boundingRect.height;
}

window.s_Weather = {
    windSpeed: 0,
    rainIntensity: 0,
    cloudiness: 0,
    temperature: 0,
    feelsLike: 0,
    humidity: 0,
    timeToSunrise: 0,
    timeToSunset: 0,
    currentTime: 0,
    currentCondition: "",
    onCloudinessChange: () => { }
}

// Set sirius_debug localStorage
if (!localStorage.getItem("sirius_debug")) {
    // Set default value to false
    localStorage.setItem("sirius_debug", "false");
}

// Make the disc spin, if he has the spinning class
document.addEventListener("readystatechange", (event: Event) => {
    if (document.readyState != "complete") {
        return;
    }
    let disc: HTMLElement | null = document.getElementById("strawberry-disc");
    let angle: number = 0;
    setInterval(() => {
        if (disc!!.classList.contains("spinning")) {
            angle += 0.25;
            disc!!.style.transform = "rotate(" + angle + "deg)";
            if (angle == 360) {
                angle = 0;
            }
        }
    }, 10);

    let socket: CafeSocket = new CafeSocket();
    let canvas: HTMLCanvasElement = document.getElementById("weather")!! as HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D = canvas.getContext("2d")!!;

    // Set the canvas size as the width of the page
    canvas.width = document.body.scrollWidth;
    canvas.height = document.body.scrollHeight * 2;

    updateTableRectangle();

    let sirius = new Sirius([
        cloudInit,
        rainInit,
    ], ctx).run();
});

// When the table is resized, resize the canvas
window.addEventListener("resize", (event: Event) => {
    let canvas: HTMLCanvasElement = document.getElementById("weather")!! as HTMLCanvasElement;
    canvas.width = window.innerWidth;
    updateTableRectangle();
});