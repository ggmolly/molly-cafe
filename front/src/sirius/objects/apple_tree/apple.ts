import { getAssets, loadAssetsType } from "../../assets";
import { Point } from "../../interfaces/point.interface";
import { AssetType } from "../../types";
import { AClickable } from "../bases/AClickable";
import { ADrawable } from "../bases/ADrawable";

const positions: Array<Point> = [
    {x: 44, y: 43},
    {x: 116, y: 16},
    {x: 151, y: 18},
    {x: 92, y: 32},
    {x: 79, y: 69},
    {x: 41, y: 82},
    {x: 133, y: 64},
    {x: 195, y: 34},
    {x: 176, y: 74},
    {x: 227, y: 70},
    {x: 206, y: 63},
    {x: 183, y: 50},
    {x: 158, y: 44},
    {x: 134, y: 39},
    {x: 116, y: 55},
    {x: 78, y: 48},
    {x: 262, y: 73},
    {x: 24, y: 63},
    {x: 230, y: 53},
    {x: 20, y: 44},
    {x: 74, y: 26},
];

const N_APPLE: number = positions.length;
let constructedApples: number = 0;
let frame = 0;

class Apple extends AClickable {
    private _ripeness: number = 61; // 0 -> 255 (200+ being rotten)
    private _id: number = 0;
    private _additionnalFilter: string = "";
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
        alpha: number = 1,
    ) {
        super(sprite, context, pos, "apple", alpha);
        this._id = constructedApples++;
        this._ripeness = constructedApples * 20 % 255;
    }

    public tick(frameDelta: number): void {
        frame++;
        if (frame % 10 == 0) { this._ripeness++; }
    }

    public preProcess(): void {
        let min: number = 0;
        let max: number = 255;
        if (this._ripeness <= 60 && this._ripeness > 0) { // green
            min = 70;
            max = 130;
            this._additionnalFilter = "brightness(120%)";
        }
        if (this._ripeness > 60 && this._ripeness <= 120) { // yellow
            min = 30;
            max = 60;
            this._additionnalFilter = "brightness(180%)";
        }
        if (this._ripeness >= 120 && this._ripeness <= 140) { // orange
            min = 40;
            max = 10;
            this._additionnalFilter = "brightness(150%)";
        }
        if (this._ripeness > 140 && this._ripeness <= 200) { // red
            min = -30;
            max = 30;
            this._additionnalFilter = "brightness(120%) contrast(140%)";
        }
        if (this._ripeness > 200) { // rotten
            min = -20;
            max = 40;
            this._additionnalFilter = "brightness(80%)";
        }
        // Normalize an angle between min and max depending on the ripeness
        let angle: number = min + (max - min) * (this._ripeness > 255 ? 1 : this._ripeness / 255);
        this._context.filter = `hue-rotate(${angle}deg) ${this._additionnalFilter}`;
    }

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
            apples.push(new Apple(appleSprites[0], ctx, positions[i]));
        }
        return apples;
    });
}
