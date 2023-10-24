import { getAssets, loadAssetsType } from "../../assets";
import { Point } from "../../interfaces/point.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";
import { AMovable } from "../bases/AMovable";

const ASCEND_DESCEND_DELAY = 3600; // 1 hour

/**
 * Returns where the sun should be drawn according to the time of the day, sunrises and sunsets.
 * If x or y is negative, the sun should not be drawn.
 */
function getSunPosition(): Point {
    // Check if the sun has risen
    if (window.s_Weather.currentTime > window.s_Weather.timeToSunrise) {
        
    }
    // If the sun has set
    if (window.s_Weather.currentTime > window.s_Weather.timeToSunset) {
        
    }
    // If the sun is rising
    if (window.s_Weather.currentTime + ASCEND_DESCEND_DELAY > window.s_Weather.timeToSunrise) {

    }
    // If the sun is setting
    if (window.s_Weather.currentTime + ASCEND_DESCEND_DELAY > window.s_Weather.timeToSunset) {

    }
}

class Sun extends AMovable {
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
    ) {
        super(sprite, context, pos, "tree", { x: 0, y: 0 });
    }

    public tick(frameDelta: number): void { return; }
}

export async function sunInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>> {
    return loadAssetsType(AssetType.SUN).then(() => {
        let sunSprites: Array<HTMLImageElement> = getAssets(AssetType.SUN);
        let suns: Array<ADrawable> = [];
        let sunSprite = sunSprites[Math.floor(Math.random() * sunSprites.length)];
        let sun = new Sun(sunSprite, ctx, { x: 0, y: 0 });
        suns.push(sun);
        return suns;
    });
}