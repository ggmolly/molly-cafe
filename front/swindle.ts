/*
    Swindle is a multi-pass 2D renderer. It is designed to be used with multiple tilemaps which are called layers.
*/

class Tile {
    public readonly width: number;
    public readonly height: number;
    public readonly image?: HTMLImageElement = undefined;
    public skip: boolean = false;
    constructor (image?: HTMLImageElement, skip: boolean = false) {
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
    constructor(width: number, height: number, tiles: Array<Tile>) {
        this.width = width;
        this.height = height;
        this.tiles = tiles;
        if (width * height !== tiles.length) throw new Error("Tilemap size does not match tile array size.");
    }

    public getTile(x: number, y: number): Tile {
        if (x >= this.width || y >= this.height || x < 0 || y < 0) throw new Error("Tilemap index out of bounds.");
        return this.tiles[y * this.width + x];
    }

    public setTile(x: number, y: number, tile: Tile) {
        if (x >= this.width || y >= this.height || x < 0 || y < 0) throw new Error("Tilemap index out of bounds.");
        this.tiles[y * this.width + x] = tile;
    }
}

class Layer {
    public readonly width: number;
    public readonly height: number;
    public readonly tilemap: Tilemap;
    constructor(width: number, height: number, tilemap: Tilemap) {
        this.width = width;
        this.height = height;
        this.tilemap = tilemap;
        // Assert that all tiles are the same size
        for (let y = 0; y < this.height; y++) {
            let tile = this.tilemap.getTile(0, y);
            for (let x = 0; x < this.width; x++) {
                let currentTile = this.tilemap.getTile(x, y);
                if (currentTile.skip && (currentTile.width !== tile.width || currentTile.height !== tile.height)) throw new Error("All tiles must be the same size.");
            }
        }
    }
}

class Swindle {
    private spriteMap: Map<string, HTMLImageElement> = new Map<string, HTMLImageElement>();
    private canvas?: HTMLCanvasElement;
    private ctx?: CanvasRenderingContext2D;

    constructor () { }

    public async init(canvas: HTMLCanvasElement): Promise<void> {
        this.canvas = canvas;
        this.ctx = canvas.getContext("2d");
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

    public async render(layers: Array<Layer>): Promise<void> {
        if (!this.ctx) throw new Error("Swindle not initialized.");
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        for (let layer of layers) {
            for (let y = 0; y < layer.height; y++) {
                for (let x = 0; x < layer.width; x++) {
                    let tile = layer.tilemap.getTile(x, y);
                    if (tile.skip) continue;
                    this.ctx.drawImage(tile.image, x * tile.width, y * tile.height);
                }
            }
        }
    }
}