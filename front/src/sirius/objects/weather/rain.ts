import { Rect } from "../../../interfaces/rect.interface";
import { getAssets, loadAssetsType } from "../../assets";
import { Velocity } from "../../interfaces/velocity.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";
import { AMovable } from "../bases/AMovable";

const N_RAINDROPS: number = 50;

class Raindrop extends AMovable {
    private _initialY: number;
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        initialVelocity: Velocity,
    ) {
        let pos = { x: Math.random() * context.canvas.width, y: 100 };
        super(sprite, context, pos, initialVelocity);
        pos.y = Math.random() * this.context.canvas.height;
        this._initialY = 0;
    }

    tick() {
        if (this.pos.y > this.context.canvas.height) {
            // random between -50 and 50
            let rnd = Math.random() * 100 - 50;
            this.pos.y = this._initialY + rnd;
            this.velocity.y = Math.random() * 2.5 + 2.5;
        }
        // Update velocity
        this.velocity.x = 0.1;
        if (this.pos.x > this.context.canvas.width - 50) {
            this.pos.x = -5;
        }
    }
}

export async function rainInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>> {
    let raindrops: Array<Raindrop> = new Array<Raindrop>();
    return loadAssetsType(AssetType.RAINDROP).then(() => {
        let raindropSprite = getAssets(AssetType.RAINDROP);
        for (let i = 0; i < N_RAINDROPS; i++) {
            raindrops.push(
                new Raindrop(
                    raindropSprite[0],
                    ctx,
                    { x: 0, y: 2.5 }
                )
            );
        }
        return raindrops;
    });
}
