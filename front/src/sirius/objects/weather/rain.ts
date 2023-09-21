import { Rect } from "../../../interfaces/rect.interface";
import { getAssets, loadAssetsType } from "../../assets";
import { Velocity } from "../../interfaces/velocity.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";
import { AMovable } from "../bases/AMovable";

const N_RAINDROPS: number = 100;

class Raindrop extends AMovable {
    private _initialY: number;
    private _parentCloud: ADrawable; // rain falls from clouds, did you know?
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        initialVelocity: Velocity,
    ) {
        let clouds = window.s_Objects.filter((obj: ADrawable) => obj.constructor.name === "Cloud");
        const parentCloud = clouds[Math.floor(Math.random() * clouds.length)];
        let pos = {
            x: parentCloud.position.x + Math.random() * parentCloud.sprite.width,
            y: parentCloud.position.y + Math.random() * context.canvas.height,
        };
        super(sprite, context, pos, initialVelocity);
        this._parentCloud = parentCloud;
        this._initialY = 0;
    }

    tick() {
        if (this.pos.y > this.context.canvas.height) {
            this.pos.y = this._parentCloud.position.y + this._parentCloud.sprite.height;
            this.velocity.y = Math.random() * 2.5 + 2.5;
        }
        // Update velocity
        this.velocity.x = 0.1;
        if (this.pos.x > this.context.canvas.width - 50) {
            this.pos.x = this._parentCloud.position.x + Math.random() * this._parentCloud.sprite.width;
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
