class Packet {
    id: number;
    data: Uint8Array;
    constructor(id: number, data: Uint8Array) {
        this.id = id;
        this.data = data;
    }
}

class Socket {
    socket: WebSocket;
    eventTicks: number = 0;
    lastX: number = 0;
    lastY: number = 0;
    constructor() {
        this.socket = new WebSocket("ws://localhost:3000/ws");
        this.socket.binaryType = "arraybuffer";
        this.socket.onopen = function (e) {
            console.log("[socket] connected!");
            window.addEventListener("mousemove", function (event) {
                // wait 8 ticks
                if (socket.eventTicks > 0) {
                    socket.eventTicks--;
                    return;
                }
                socket.eventTicks = 8;
                let packet = new Packet(MOUSE_MOVE_ID, new Uint8Array(5));
                // set two bytes for x and two bytes for y
                let buffer = new ArrayBuffer(5);
                let data = new DataView(buffer);
                data.setUint8(0, MOUSE_MOVE_ID);
                // Instead of sending the raw x and y values, we send a value between 0 and 65535
                data.setUint16(1, event.clientX / window.innerWidth * 65535);
                data.setUint16(3, event.clientY / window.innerHeight * 65535);
                packet.data = new Uint8Array(buffer);
                socket.socket.send(packet.data.buffer);
            });
        }
        this.socket.onclose = function (e) {
            console.log("[socket] disconnected!");
            if (e.wasClean) {
                console.log(`[socket] closed cleanly, code=${e.code} reason=${e.reason}`);
            } else {
                console.log('[socket] connection died, attempting to reconnect in 5 seconds...');
                setTimeout(() => {
                    socket = new Socket();
                }, 5000);
            }
        }
        this.socket.onmessage = function (e) {
            let arrayBuffer = e.data;
            let byteArray = new Uint8Array(arrayBuffer);
            let packet = new Packet(byteArray[0], byteArray.slice(1));
            let socketId = 0;
            if (packet.id & CLIENT_PACKET_TYPE) {
                // copy the first 2 bytes of the packet data
                socketId = new DataView(packet.data.buffer).getUint16(0);
            }
            switch (packet.id) {
                case WELCOME_PACKET_ID:
                    addPlayer(socketId);
                    break;
                case CYA_PACKET_ID:
                    removePlayer(socketId);
                    break;
                case MOUSE_MOVE_ID:
                    let data = new DataView(packet.data.buffer);
                    // x and y are two uint16 values between 0 and 65535
                    let x = data.getUint16(2 + 0) / 65535 * window.innerWidth;
                    let y = data.getUint16(2 + 2) / 65535 * window.innerHeight;
                    movePlayer(socketId, x, y);
                    break;
            }
        }
    }
}

function addPlayer(socketId) {
    let player = document.createElement("div");
    player.id = socketId;
    player.style.width = "24px";
    player.style.height = "24px";
    player.style.borderRadius = "50%";
    player.style.backgroundColor = "#" + Math.floor(Math.random() * 16777215).toString(16);
    player.style.position = "absolute";
    player.style.left = "0px";
    player.style.top = "0px";
    document.body.appendChild(player);
}

function removePlayer(socketId) {
    // remove the player from the dom
    let player = document.getElementById(socketId);
    if (!player) return;
    document.body.removeChild(player);
}

function movePlayer(socketId, x, y) {
    let player = document.getElementById(socketId);
    if (!player) return;
    player.style.left = x + "px";
    player.style.top = y + "px";
}

const WELCOME_PACKET_ID = 0b01000000;
const CYA_PACKET_ID = 0b01000001;
const MOUSE_MOVE_ID = 0b01000010;
const CLIENT_PACKET_TYPE = 0b01000000;
let socket = new Socket();

document.addEventListener("readystatechange", async function (event) {
    if (document.readyState !== "complete") return;

    let swindle: Swindle = new Swindle();
    let tiles = Array<Tile>();

    const default_grass = await swindle.loadSprite("default_grass", "assets/grasses/grass_0.png");
    const top_border = await swindle.loadSprite("top_border", "assets/borders/top_border.png");
    const bottom_border = await swindle.loadSprite("bottom_border", "assets/borders/bottom_border.png");
    const left_border = await swindle.loadSprite("left_border", "assets/borders/left_border.png");
    const right_border = await swindle.loadSprite("right_border", "assets/borders/right_border.png");
    const top_left_corner = await swindle.loadSprite("top_left_corner", "assets/borders/top_left_border.png");
    const top_right_corner = await swindle.loadSprite("top_right_corner", "assets/borders/top_right_border.png");
    const bottom_left_corner = await swindle.loadSprite("bottom_left_corner", "assets/borders/bottom_left_border.png");
    const bottom_right_corner = await swindle.loadSprite("bottom_right_corner", "assets/borders/bottom_right_border.png");
    const small_sapling = await swindle.loadSprite("small_sapling", "assets/objects/plants_3.png");

    swindle.init(document.getElementById("canvas") as HTMLCanvasElement);

    for (let y = 0; y < 3; y++) {
        for (let x = 0; x < 3; x++) {
            if (x === 0 && y === 0) {
                tiles.push(new Tile(top_left_corner));
                continue ;
            }
            if (x === 0 && y === 2) {
                tiles.push(new Tile(bottom_left_corner));
                continue ;
            }
            if (x === 2 && y === 0) {
                tiles.push(new Tile(top_right_corner));
                continue ;
            }
            if (x === 2 && y === 2) {
                tiles.push(new Tile(bottom_right_corner));
                continue ;
            }
            if (x > 0 && x < 2 && y === 0) {
                tiles.push(new Tile(top_border));
                continue ;
            }
            if (x > 0 && x < 2 && y === 2) {
                tiles.push(new Tile(bottom_border));
                continue ;
            }
            if (x === 0 && y > 0 && y < 2) {
                tiles.push(new Tile(left_border));
                continue ;
            }
            if (x === 2 && y > 0 && y < 2) {
                tiles.push(new Tile(right_border));
                continue ;
            }
            tiles.push(new Tile(default_grass));
        }
    }

    const emptyTile = new Tile();
    let tileMap: Tilemap = new Tilemap(3, 3, tiles);
    let layers: Array<Layer> = new Array<Layer>();

    // To prove that layers are working, we're gonna add a small sapling at 1, 1
    let saplingTileMapTiles = Array<Tile>();
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(new Tile(small_sapling));
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    saplingTileMapTiles.push(emptyTile);
    let saplingTileMap: Tilemap = new Tilemap(3, 3, saplingTileMapTiles);
    let saplingLayer: Layer = new Layer(3, 3, saplingTileMap);
    layers.push(new Layer(3, 3, tileMap));
    layers.push(saplingLayer);
    swindle.render(layers);
});