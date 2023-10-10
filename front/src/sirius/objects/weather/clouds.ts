import { Rect } from "../../../interfaces/rect.interface";
import { getAssets, loadAssetsType } from "../../assets";
import { Velocity } from "../../interfaces/velocity.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";
import { AMovable } from "../bases/AMovable";

const N_CLOUDS: number = 100;
const speedX: number = 0.10;
const yDelta: number = 0.003;
const maxCloudY = 200;
let constructedClouds = 0;

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

class Cloud extends AMovable {
    private _consecutiveSteps: number;
    private _lastDirection: number;
    private _bounciness: number;
    private _idleFrames: number;

    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        initialVelocity: Velocity
    ) {
        let pos = { x: Math.random() * context.canvas.width, y: Math.random() * maxCloudY };
        super(sprite, context, pos, initialVelocity);
        this._consecutiveSteps = 0; // number of consecutive steps in the same direction (used for bouncing)
        this._lastDirection = Math.random() > 0.5 ? 1 : -1; // random initial direction
        // bounciness is a random number between 100 and 200 (rounded)
        // it represents the number of frames the cloud will go up and down
        this._bounciness = Math.round(Math.random() * 100) + 100;
        // idle frames is a random number between 10 and 80 (rounded)
        // it represents the number of frames the cloud will stay idle after bouncing
        this._idleFrames = Math.round(Math.random() * 70) + 10;
        if ((constructedClouds / N_CLOUDS) < window.s_Weather.cloudiness) {
            this.enable();
        } else {
            this.disable();
        }
        constructedClouds++;
    }

    tick() {
        if (this.pos.x > this.context.canvas.width) {
            this.pos.x = -this.sprite.width;
        }
        this.enabled = !behindTable({ x: this.pos.x, y: this.pos.y, width: this.sprite.width, height: this.sprite.height }) && this.enabled;
        // Update velocity
        this.velocity.x = speedX + Math.log2(window.s_Weather.windSpeed / 100 + 1);
        if (this._consecutiveSteps == this._bounciness) {
            this._consecutiveSteps = this._idleFrames;
            this._lastDirection *= -1;
        }
        this._consecutiveSteps++;
        this.velocity.y = yDelta * this._lastDirection;
    }
}

export async function cloudInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>> {
    window.s_Weather.onCloudinessChange = onCloudinessChange;
    return loadAssetsType(AssetType.CLOUD).then(() => {
        let cloudSprites: Array<HTMLImageElement> = getAssets(AssetType.CLOUD);
        let clouds: Array<Cloud> = new Array<Cloud>();
        for (let i = 0; i < N_CLOUDS; i++) {
            clouds.push(
                new Cloud(
                    cloudSprites[Math.floor(Math.random() * cloudSprites.length)],
                    ctx,
                    { x: speedX, y: 0 }
                )
            );
        }
        return clouds;
    });
}

/**
 * This function will get called whenever window.s_Weather.cloudiness changes
 * It will be called with the new value of cloudiness, and the old value, allowing a smooth transition
 */
export function onCloudinessChange(newCloudiness: number) {
    let clouds = window.s_Objects.filter((obj) => obj.constructor.name === 'Cloud');
    let shownClouds: number = 0;
    // Show newCloudiness % of clouds
    for (const cloud of clouds) {
        if ((shownClouds / clouds.length * 100) < newCloudiness) {
            cloud.enable();
            shownClouds++;
        } else {
            cloud.disable();
        }
    }
}