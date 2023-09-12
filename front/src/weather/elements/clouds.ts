import { Rect } from "../../interfaces/rect.interface";
import { getAssets, loadAssetsType } from "../assets";
import { AssetType } from "../types";

function behindTable(imageRect: Rect): boolean {
    if (window.tableRect == undefined) {
        return false;
    }
    if (
        imageRect.x > window.tableRect.x && // left side
        imageRect.x < window.tableRect.x + window.tableRect.width - imageRect.width && // right side
        imageRect.y > window.tableRect.y // top side
    ) {
        return true;
    }
    return false;
}

async function cloudRendering(canvas: HTMLCanvasElement, ctx: CanvasRenderingContext2D) {
    const speedX: number = 0.20;
    const maxYDelta: number = 0.5;
    let clouds: Array<HTMLImageElement> = getAssets(AssetType.CLOUD);
    let cloudPositions: Array<Array<number>> = new Array<Array<number>>();
    for (let i = 0; i < clouds.length; i++) {
        cloudPositions.push([Math.random() * canvas.width, Math.random() * canvas.height]);
    }
    setInterval(() => {
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        for (let i = 0; i < clouds.length; i++) {
            let image: HTMLImageElement = clouds[i];
            ctx.globalAlpha = i / clouds.length;
            if (behindTable({ x: cloudPositions[i][0], y: cloudPositions[i][1], width: image.width, height: image.height })) {
                continue;
            }
            ctx.drawImage(image, cloudPositions[i][0], cloudPositions[i][1]);
            ctx.globalAlpha = 1;
        }
        for (let i = 0; i < cloudPositions.length; i++) {
            cloudPositions[i][0] += speedX;
            if (cloudPositions[i][0] > canvas.width) {
                cloudPositions[i][0] = -clouds[i].width;
            }
            // Randomly decrease or increase cloud height
            let direction: boolean = Math.random() > 0.5;
            // If new height is out of bounds, reverse direction
            if (cloudPositions[i][1] < 0) {
                direction = false;
            } else if (cloudPositions[i][1] > canvas.height - clouds[i].height) {
                direction = true;
            }
            if (direction) { // Decrease
                cloudPositions[i][1] -= Math.random() * maxYDelta;
            } else { // Increase
                cloudPositions[i][1] += Math.random() * maxYDelta;
            }
        }
    }, 10);
}

/**
 * Main function of the clouds module, this will load assets (once) and start the rendering loop
 */
export function cloudMain(canvas: HTMLCanvasElement, ctx: CanvasRenderingContext2D) {
    // Load assets, then when they're loaded, start rendering
    loadAssetsType(AssetType.CLOUD).then(() => {
        cloudRendering(canvas, ctx);
    });
}