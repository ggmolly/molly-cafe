import { Point } from "../../interfaces/point.interface";
import { Velocity } from "../../interfaces/velocity.interface";
import { ADrawable } from "./ADrawable";

export abstract class AMovable extends ADrawable {
    protected velocity: Velocity;

    /**
     * Constructs a movable object
     * @param sprite Image of the object
     * @param context Context of the canvas
     * @param pos Position of the object
     * @param initialVelocity Initial velocity of the object
     */
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
        type: string = "AMovable",
        initialVelocity: Velocity,
    ) {
        super(sprite, context, pos, type);
        this.velocity = initialVelocity;
    }

    public _tick(frameDelta: number) {
        if (frameDelta > 100) { return; }
        this.pos.x += this.velocity.x * frameDelta;
        this.pos.y += this.velocity.y * frameDelta;
        this.tick(frameDelta);
        if (this.enabled) { this.draw(); }
    }

    public setVelocity(velocity: Velocity) {
        this.velocity = velocity;
    }
}