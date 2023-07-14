var socket = new WebSocket("ws://localhost:3000/ws");
socket.binaryType = "arraybuffer";

const WELCOME_PACKET_ID = 0b01000000;
const CYA_PACKET_ID = 0b01000001;
const MOUSE_MOVE_ID = 0b01000010;

const CLIENT_PACKET_TYPE = 0b01000000;

let lastX, lastY;

class Packet {
    constructor(id, data) {
        this.id = id;
        this.data = data;
    }
}

function addPlayer(uuid) {
    let player = document.createElement("div");
    player.id = uuid;
    player.style.width = "24px";
    player.style.height = "24px";
    player.style.borderRadius = "50%";
    player.style.backgroundColor = "#" + Math.floor(Math.random() * 16777215).toString(16);
    player.style.position = "absolute";
    player.style.left = "0px";
    player.style.top = "0px";
    document.body.appendChild(player);
}

function removePlayer(uuid) {
    // remove the player from the dom
    let player = document.getElementById(uuid);
    if (!player) return ;
    document.body.removeChild(player);
}

function movePlayer(uuid, x, y) {
    let player = document.getElementById(uuid);
    if (!player) return ;
    player.style.left = x * window.innerWidth + "px";
    player.style.top = y * window.innerHeight + "px";
}

socket.onopen = function (event) {
    window.addEventListener("mousemove", function (event) {
        // must move at least 2px to send a packet
        if (Math.abs(event.clientX - lastX) < 2 && Math.abs(event.clientY - lastY) < 2) return ;
        lastX = event.clientX;
        lastY = event.clientY;
        let packet = new Packet(MOUSE_MOVE_ID, new Uint8Array(5));
        // set two bytes for x and two bytes for y
        let buffer = new ArrayBuffer(9);
        let data = new DataView(buffer);
        data.setUint8(0, MOUSE_MOVE_ID);
        data.setFloat32(1, event.clientX / window.innerWidth);
        data.setFloat32(5, event.clientY / window.innerHeight);
        packet.data = new Uint8Array(buffer);
        socket.send(packet.data);
    });
}


socket.onmessage = function (event) {
    let arrayBuffer = event.data;
    let byteArray = new Uint8Array(arrayBuffer);
    let packet = new Packet(byteArray[0], byteArray.slice(1));
    let uuid = "";
    if (packet.id & CLIENT_PACKET_TYPE) {
        // copy the first 36 bytes into the uuid
        for (let i = 0; i < 36; i++) {
            uuid += String.fromCharCode(packet.data[i]);
        }
    }
    switch (packet.id) {
        case WELCOME_PACKET_ID:
            addPlayer(uuid);
            break;
        case CYA_PACKET_ID:
            removePlayer(uuid);
            break;
        case MOUSE_MOVE_ID:
            let data = new DataView(packet.data.buffer);
            // x and y are two float32 values
            let x = data.getFloat32(36 + 0);
            let y = data.getFloat32(36 + 4);
            movePlayer(uuid, x, y);
            break;
    }
}