import { getAssets, loadAssetsType } from "../../assets";
import { Point } from "../../interfaces/point.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";

const N_TREE: number = 1;

class Tree extends ADrawable {
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        pos: Point,
        alpha: number = 1,
    ) {
        super(sprite, context, pos, "tree", alpha);
        context.canvas.height = sprite.height;
        context.canvas.width = sprite.width;
    }

    public tick(frameDelta: number): void { return; }
}

export async function treeInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>> {
    return loadAssetsType(AssetType.TREE).then(() => {
        let treeSprites: Array<HTMLImageElement> = getAssets(AssetType.TREE);
        let trees: Array<ADrawable> = [];
        for (let i = 0; i < N_TREE; i++) {
            let treeSprite = treeSprites[Math.floor(Math.random() * treeSprites.length)];
            let tree = new Tree(treeSprite, ctx, { x: 0, y: 0 });
            trees.push(tree);
        }
        return trees;
    });
}