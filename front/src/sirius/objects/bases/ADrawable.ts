import { Point } from "../../interfaces/point.interface";

/**
 * Abstract class for drawable objects
 * 
 * This class is used to define the common properties of drawable objects such as:
 * 
 * - pos: the position of the object
 * - sprite: the sprite of the object
 * - context: the context of the canvas
 * - enabled: a boolean to know whether the object should be drawn or not
 * NOTE: the draw method will not be called if enabled is false
 */
export abstract class ADrawable {
    protected _sprite: HTMLImageElement;
    protected _context: CanvasRenderingContext2D;
    private _type: string = "ADrawable";
    private _clickable: boolean = false;
    public pos: Point;
    public enabled: boolean = true;
    public alpha: number = 1;

    /**
     * Constructs a drawable object
     * @param sprite Image of the object
     * @param context Context of the canvas
     * @param pos Position of the object
     */
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
        type: string,
        alpha: number = 1,
        clickable: boolean = false
    ) {
        this._sprite = sprite;
        this._context = context;
        this._type = type;
        this.pos = pos;
        console.assert(alpha >= 0 && alpha <= 1, "[sirius] Alpha must be between 0 and 1");
        this.alpha = alpha;
        this._clickable = clickable;
    }

    public preProcess(): void { return; }

    /**
     * Draws the object on the canvas, isn't called if enabled is false
     */
    public draw() {
        this._context.globalAlpha = this.alpha;
        this.preProcess();
        // bitwise hack to remove decimals
        this._context.drawImage(
            this.sprite,
            this.pos.x | 0,
            this.pos.y | 0,
            this.sprite.width,
            this.sprite.height
        );
        this._context.globalAlpha = 1;
        this._context.filter = "none";
    }

    /**
     * Ticks the object, called every frame, even if the object is disabled, calls tick
     * It is used to update the object's properties (if any) in order to not let the tick method do it
     * Internal method, not meant to be called by anything else than the rendering engine itself
     */
    public _tick(frameDelta: number) {
        this.tick(frameDelta);
        if (this.enabled) { this.draw(); }
    }

    /**
     * Tick method, called every frame, even if the object is disabled
     */
    public abstract tick(frameDelta: number): void;

    /**
     * Disables the object, it won't be drawn anymore (noop if already disabled)
     */
    public disable() {
        this.enabled = false;
    }

    /**
     * Enables the object, it will be drawn again (noop if already enabled)
     */
    public enable() {
        this.enabled = true;
    }

    /**
     * Updates the object's sprite
     */
    public updateSprite(sprite: HTMLImageElement) {
        this._sprite = sprite;
    }

    /**
     * Updates the object's context
     */
    public updateContext(context: CanvasRenderingContext2D) {
        this._context = context;
    }

    /**
     * Checks if a passed point is inside the object
     */
    public isPointInside(point: Point): boolean {
        return (
            point.x >= this.pos.x &&
            point.x <= this.pos.x + this.sprite.width &&
            point.y >= this.pos.y &&
            point.y <= this.pos.y + this.sprite.height
        );
    }

    /**
     * Returns the object's position
     */
    public get position(): Point {
        return this.pos;
    }

    /**
     * Returns the object's sprite
     */
    public get sprite(): HTMLImageElement {
        return this._sprite;
    }

    /**
     * Returns the object's context
    **/
    public get context(): CanvasRenderingContext2D {
        return this._context;
    }

    /**
     * Returns the object's type
     */
    public get type(): string {
        return this._type;
    }

    /**
     * Returns the object's clickable property
     */
    public get clickable(): boolean {
        return this._clickable;
    }
}