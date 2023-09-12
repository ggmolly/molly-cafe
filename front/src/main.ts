import { cloudInit } from "./rendering/objects/weather/clouds";
import { Sirius } from "./rendering/sirius";
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
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;

    updateTableRectangle();

    window.windSpeed = (document.getElementById("wind-speed")!! as HTMLInputElement).valueAsNumber;

    document.getElementById("wind-speed")!!.addEventListener("change", (event: Event) => {
        window.windSpeed = (event.target as HTMLInputElement).valueAsNumber;
        console.log("Wind speed changed to " + window.windSpeed);
    });

    let sirius = new Sirius([
        cloudInit,
    ], ctx).run();
});

// When the table is resized, resize the canvas
window.addEventListener("resize", (event: Event) => {
    let canvas: HTMLCanvasElement = document.getElementById("weather")!! as HTMLCanvasElement;
    canvas.width = window.innerWidth;
    updateTableRectangle();
});