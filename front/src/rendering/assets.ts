// This file will contains all the loaded images

import { AssetType } from "./types";

const typeCount: Record<AssetType, number> = {
    [AssetType.CLOUD]: 11,
}
let assetsMap: Record<AssetType, Array<HTMLImageElement>> = {
    [AssetType.CLOUD]: new Array<HTMLImageElement>(),
}

/**
 * Loads a single image, and returns a promise that resolves when the image is loaded
 * @param type Type of asset to load (cloud, etc...)
 * @param index Index of the asset to load (0 -> cloud_0.png, 1 -> cloud_1.png, etc...)
 * @returns Promise that resolves when the image is loaded
 */
function loadAsset(type: AssetType, index: number): Promise<void> {
    return new Promise((resolve, reject) => {
        let image: HTMLImageElement = new Image();
        image.src = "/assets/" + type + "_" + index + ".png";
        image.onload = () => {
            assetsMap[type].push(image);
            resolve();
        }
    });
}

/**
 * Loads all assets of a type
 * @param type Type of asset to load (cloud, etc...)
 * @returns Promise that resolves when all assets are loaded
 */
export function loadAssetsType(type: AssetType): Promise<void> {
    // If assets are already loaded
    if (assetsMap[type].length == typeCount[type]) {
        console.warn("Assets of type " + type + " are already loaded, not loading again.")
        return new Promise((resolve, reject) => {
            resolve();
        });
    }
    const max: number = typeCount[type];
    let promises: Array<Promise<void>> = new Array<Promise<void>>();
    for (let i = 0; i < max; i++) {
        promises.push(loadAsset(type, i));
    }
    return new Promise((resolve, reject) => {
        Promise.all(promises).then(() => {
            resolve();
        });
    });
}

export function getAssets(type: AssetType): Array<HTMLImageElement> {
    return assetsMap[type];
}