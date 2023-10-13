import { ADrawable } from "./objects/bases/ADrawable";

const avgSampleSize: number = 100;
let frameCount: number = 0;
let frameTimes: Array<number> = [];

/**
 * Main class for the rendering engine
 */
export class Sirius {
    private _init_functions: Array<CallableFunction>;
    private _ctx: CanvasRenderingContext2D;
    private _debugFrameTime: HTMLSpanElement | null;
    private _debugFrameRate: HTMLSpanElement | null;
    private _debugLoadedObjects: HTMLSpanElement | null;
    private _debugAvgFrameTime: HTMLSpanElement | null;
    private _debugAvgFrameRate: HTMLSpanElement | null;
    private _debugRenderedObjects: HTMLSpanElement | null;
    private _debugEnabled: boolean = false;
    /**
     * Initializes the rendering engine
     * @param functions Array of async functions that return a promise of an array of ADrawable objects
     */
    constructor(functions: Array<CallableFunction>, ctx: CanvasRenderingContext2D) {
        window.s_Objects = [];
        this._init_functions = functions;
        this._ctx = ctx;
        // Check the localStorage value
        this._debugFrameTime = document.getElementById("sirius-debug-ft");
        this._debugFrameRate = document.getElementById("sirius-debug-fps");
        this._debugLoadedObjects = document.getElementById("sirius-debug-objs");
        this._debugAvgFrameTime = document.getElementById("sirius-debug-avg-ft");
        this._debugAvgFrameRate = document.getElementById("sirius-debug-avg-fps");
        this._debugRenderedObjects = document.getElementById("sirius-debug-rendered-objs");
        this._debugEnabled = localStorage.getItem("sirius_debug") === "true";
        if (!this._debugEnabled) {
            // remove the DOM element to not take up unnecessary space
            document.getElementById("sirius-debug")?.remove();
        }
        this._init();
    }

    private _getAverageFrameTime(): number {
        let sum: number = 0;
        frameTimes.forEach(t => sum += t);
        return sum / frameTimes.length;
    }

    /**
     * Calls all the init functions and flattens the resulting arrays into a single array (window.s_Objects)
     */
    private async _init(): Promise<void> {
        for (let i = 0; i < this._init_functions.length; i++) {
            let f = this._init_functions[i];
            window.s_Objects = window.s_Objects.concat(await f(this._ctx));
            console.debug(`[sirius] Initialized ${window.s_Objects.length} objects`);
            console.log(`[sirius] ${i + 1}/${this._init_functions.length} init functions called`);
        }
    }

    private _updateDebug(frameTime: number) {
        if (
            !this._debugEnabled ||
            !this._debugFrameTime ||
            !this._debugFrameRate ||
            !this._debugLoadedObjects ||
            !this._debugAvgFrameTime ||
            !this._debugAvgFrameRate ||
            !this._debugRenderedObjects
        ) { return; }
        if (frameCount % 30 != 0) { return; }
        this._debugFrameTime.innerText = "Frame time: " + frameTime.toFixed(2) + "ms";
        this._debugFrameRate.innerText = "Frame rate: " + (1000 / frameTime).toFixed(2) + "fps";
        this._debugLoadedObjects.innerText = "Loaded objects: " + window.s_Objects.length;
        this._debugAvgFrameTime.innerText = "Average frame time: " + this._getAverageFrameTime().toFixed(2) + "ms";
        this._debugAvgFrameRate.innerText = "Average frame rate: " + (1000 / this._getAverageFrameTime()).toFixed(2) + "fps";
        this._debugRenderedObjects.innerText = "Rendered objects: " + window.s_Objects.filter(o => o.enabled).length;
    }

    /**
     * Runs the rendering engine, calling the tick() method of all the objects, and then the draw() method
     * Looping using requestAnimationFrame
     */
    run() {
        let start: number = performance.now();
        let tick = () => {
            this._ctx.clearRect(0, 0, this._ctx.canvas.width, this._ctx.canvas.height);
            let frameDelta: number = performance.now() - start;
            window.s_Objects.forEach(o => o._tick(frameDelta));
            frameCount = requestAnimationFrame(tick);
            this._updateDebug(frameDelta);
            if (frameTimes.length >= avgSampleSize) {
                frameTimes.shift();
            }
            frameTimes.push(frameDelta);
            start = performance.now();
        }
        tick();
    }
}