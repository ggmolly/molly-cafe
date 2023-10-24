import { Point } from "../../interfaces/point.interface";
import { ADrawable } from "./ADrawable";

export abstract class AClickable extends ADrawable {
    /**
     * Constructs a clickable object
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
    ) {
        super(sprite, context, pos, type, alpha, true);
    }

    /**
     * Called when the object is clicked on
     */
    public abstract onClick(e: MouseEvent): void;
}
