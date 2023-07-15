/*
    Swindle is a multi-pass 2D renderer. It is designed to be used with multiple tilemaps which are called layers.
*/

class Tile {
    public readonly width: number;
    public readonly height: number;
    public readonly image?: HTMLImageElement = undefined;
    public skip: boolean = false;
    constructor(image?: HTMLImageElement, skip: boolean = false) {
        if (image) {
            this.width = image.width;
            this.height = image.height;
        } else {
            this.skip = true;
        }
        this.image = image;
    }
}

class Tilemap {
    private tiles: Array<Tile>; // Flatten'd Matrix
    public readonly width: number;
    public readonly height: number;
    public dirty: boolean = true;
    constructor(width: number, height: number, preFillImage?: HTMLImageElement) {
        this.width = width;
        this.height = height;
        this.tiles = new Array<Tile>(width * height);
        // fill with empty tiles
        if (!preFillImage) return;
        for (let i = 0; i < this.tiles.length; i++) {
            this.tiles[i] = new Tile(preFillImage);
        }
    }

    public getTile(x: number, y: number): Tile {
        if (x >= this.width || y >= this.height || x < 0 || y < 0) throw new Error("Tilemap index out of bounds.");
        return this.tiles[y * this.width + x];
    }

    public setTile(x: number, y: number, tile: Tile) {
        if (x >= this.width || y >= this.height || x < 0 || y < 0) throw new Error("Tilemap index out of bounds.");
        const oldState = this.tiles[y * this.width + x];
        if (oldState.image === tile.image) return;
        this.dirty = true;
        this.tiles[y * this.width + x] = tile;
    }

    public removeTile(x: number, y: number) {
        if (x >= this.width || y >= this.height || x < 0 || y < 0) throw new Error("Tilemap index out of bounds.");
        this.dirty = true;
        this.tiles[y * this.width + x] = null;
    }
}

/**
 * A SwindleClickEvent is fired when the user clicks on the canvas.
 * 
 * @param x The x coordinate of the click.
 * @param y The y coordinate of the click.
 * @param tile The tile that was clicked on.
 */
class SwindleClickEvent {
    public readonly x: number;
    public readonly y: number;
    public readonly tileX: number;
    public readonly tileY: number;
    public readonly tile: Tile;
    constructor(x: number, y: number, tileX: number, tileY: number, tile: Tile) {
        this.x = x;
        this.y = y;
        this.tileX = tileX;
        this.tileY = tileY;
        this.tile = tile;
    }
}

/**
 * A SwindleTileHoverEvent is fired when the user hovers over a tile.
 * 
 * @param x The x coordinate of the hover.
 * @param y The y coordinate of the hover.
 * @param tileX The x coordinate of the tile that was hovered over.
 * @param tileY The y coordinate of the tile that was hovered over.
 * @param tile The tile that was hovered over.
 */
class SwindleTileHoverEvent {
    public readonly x: number;
    public readonly y: number;
    public readonly tileX: number;
    public readonly tileY: number;
    public readonly tile: Tile;
    constructor(x: number, y: number, tileX: number, tileY: number, tile: Tile) {
        this.x = x;
        this.y = y;
        this.tileX = tileX;
        this.tileY = tileY;
        this.tile = tile;
    }
}

class Swindle {
    private spriteMap: Map<string, HTMLImageElement> = new Map<string, HTMLImageElement>();
    private canvas?: HTMLCanvasElement;
    private ctx?: CanvasRenderingContext2D;
    private tileWidth: number = 0;
    private tileHeight: number = 0;
    private fpsCounter?: HTMLElement;
    private lastFrameTime: number = 0;
    private frameCount: number = 0;
    public tilemaps: Array<Tilemap> = new Array<Tilemap>();
    public onSwindleClick: ((e: SwindleClickEvent) => void) | null = null;
    public onSwindleTileHover: ((e: SwindleTileHoverEvent) => void) | null = null;

    constructor() { }

    public async init(tileWidth: number, tileHeight: number, canvas: HTMLCanvasElement, fpsCounter?: HTMLElement) {
        this.canvas = canvas;
        this.ctx = canvas.getContext("2d");
        if (!this.ctx) throw new Error("Could not get 2D context.");
        this.ctx.imageSmoothingEnabled = false;
        this.tileWidth = tileWidth;
        this.tileHeight = tileHeight;
        // add click event
        canvas.addEventListener("click", (e) => {
            if (!this.ctx) throw new Error("Swindle not initialized.");
            let rect = this.canvas.getBoundingClientRect();
            let x = e.clientX - rect.left;
            let y = e.clientY - rect.top;
            let tileX = Math.floor(x / this.tileWidth);
            let tileY = Math.floor(y / this.tileHeight);
            let tile = this.tilemaps[this.tilemaps.length - 1].getTile(tileX, tileY);
            let event = new SwindleClickEvent(x, y, tileX, tileY, tile);
            if (this.onSwindleClick) this.onSwindleClick(event);
        });
        // add hover event
        canvas.addEventListener("mousemove", (e) => {
            if (!this.ctx) throw new Error("Swindle not initialized.");
            let rect = this.canvas.getBoundingClientRect();
            let x = e.clientX - rect.left;
            let y = e.clientY - rect.top;
            let tileX = Math.floor(x / this.tileWidth);
            let tileY = Math.floor(y / this.tileHeight);
            const tilemap = this.tilemaps[this.tilemaps.length - 1];
            if (tileX >= tilemap.width || tileY >= tilemap.height || tileX < 0 || tileY < 0) return;
            let tile = this.tilemaps[this.tilemaps.length - 1].getTile(tileX, tileY);
            let event = new SwindleTileHoverEvent(x, y, tileX, tileY, tile);
            if (this.onSwindleTileHover) this.onSwindleTileHover(event);
        });
        this.fpsCounter = fpsCounter;
    }

    public async loadSprite(name: string, path: string): Promise<HTMLImageElement> {
        if (this.spriteMap.has(name)) return this.spriteMap.get(name);
        let img = new Image();
        img.src = path;
        await img.decode();
        this.spriteMap.set(name, img);
        return img;
    }

    public async getSprite(name: string): Promise<HTMLImageElement> {
        if (!this.spriteMap.has(name)) throw new Error("Sprite not loaded: " + name);
        return this.spriteMap.get(name);
    }

    public async render(): Promise<void> {
        if (!this.ctx) throw new Error("Swindle not initialized.");
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        for (let layer of this.tilemaps) {
            if (!layer.dirty) continue;
            console.log("Rendering layer: " + layer)
            for (let y = 0; y < layer.height; y++) {
                for (let x = 0; x < layer.width; x++) {
                    let tile = layer.getTile(x, y);
                    if (!tile) continue;
                    if (tile.skip) continue;
                    this.ctx.drawImage(tile.image, x * tile.width, y * tile.height);
                }
            }
            layer.dirty = false;
        }
        this.frameCount++;
        let now = performance.now();
        let fps = 1000 / (now - this.lastFrameTime);
        this.lastFrameTime = now;
        if (this.fpsCounter && this.frameCount % 10 == 0) {
            this.fpsCounter.innerText = fps.toFixed(2) + " FPS";
        }
    }
}
