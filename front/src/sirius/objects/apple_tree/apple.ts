import { getAssets, loadAssetsType } from "../../assets";
import { Point } from "../../interfaces/point.interface";
import { AssetType } from "../../types";
import { AClickable } from "../bases/AClickable";
import { ADrawable } from "../bases/ADrawable";

const N_APPLE: number = 1;

class Apple extends AClickable {
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
        alpha: number = 1,
    ) {
        super(sprite, context, pos, "apple", alpha);
    }

    public tick(frameDelta: number): void { return; }

    public onClick(): void {
        alert("You clicked an apple!");
    }
}

export async function appleInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>>
{
    return loadAssetsType(AssetType.APPLE).then(() => {
        let appleSprites: Array<HTMLImageElement> = getAssets(AssetType.APPLE);
        let apples: Array<ADrawable> = [];
        for (let i = 0; i < N_APPLE; i++) {
            // Transform the apple to a random color
            apples.push(new Apple(appleSprites[0], ctx, { x: 50, y: 20 }));
        }
        return apples;
    });
}
