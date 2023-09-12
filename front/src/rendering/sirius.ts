import { ADrawable } from "./objects/bases/ADrawable";

const avgSampleSize: number = 100;
let frameCount: number = 0;
let frameTimes: Array<number> = [];

/**
 * Main class for the rendering engine
 */
export class Sirius {
    private _init_functions: Array<CallableFunction>;
    private _objects: Array<ADrawable>;
    private _ctx: CanvasRenderingContext2D;
    private _debugSpan: HTMLElement | null;
    /**
     * Initializes the rendering engine
     * @param functions Array of async functions that return a promise of an array of ADrawable objects
     */
    constructor(functions: Array<CallableFunction>, ctx: CanvasRenderingContext2D) {
        this._init_functions = functions;
        this._objects = [];
        this._ctx = ctx;
        this._debugSpan = document.getElementById("sirius-debug");
        this._init();
    }

    private _getAverageFrameTime(): number {
        let sum: number = 0;
        frameTimes.forEach(t => sum += t);
        return sum / frameTimes.length;
    }

    /**
     * Calls all the init functions and flattens the resulting arrays into a single array (this._objects)
     */
    private _init() {
        this._init_functions.forEach(f => {
            f(this._ctx).then((objects: Array<ADrawable>) => {
                this._objects = this._objects.concat(objects);
                console.debug("[sirius] Initialized " + objects.length + " objects");
            });
        });
    }

    private _updateDebug(frameTime: number) {
        if (!this._debugSpan) { return; }
        if (frameCount % 30 != 0) { return; }
        this._debugSpan.innerText = "Frame time: " + frameTime.toFixed(2) + "ms (" + (1000 / frameTime).toFixed(2) + "fps) | Average frame time: " + this._getAverageFrameTime().toFixed(2) + "ms (" + (1000 / this._getAverageFrameTime()).toFixed(2) + "fps)";
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
            this._objects.forEach(o => o._tick(frameDelta));
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