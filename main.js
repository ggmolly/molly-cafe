var socket = new WebSocket("ws://localhost:3000/ws");
socket.binaryType = "arraybuffer";

const WELCOME_PACKET_ID = 0b01000000;
const CYA_PACKET_ID = 0b01000001;
const MOUSE_MOVE_ID = 0b01000010;

const CLIENT_PACKET_TYPE = 0b01000000;

let eventTicks = 8;

class Packet {
    constructor(id, data) {
        this.id = id;
        this.data = data;
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
    if (!player) return ;
    document.body.removeChild(player);
}

function movePlayer(socketId, x, y) {
    let player = document.getElementById(socketId);
    if (!player) return ;
    player.style.left = x + "px";
    player.style.top = y + "px";
}

socket.onopen = function (event) {
    window.addEventListener("mousemove", function (event) {
        // wait 8 ticks
        lastX = event.clientX;
        lastY = event.clientY;
        if (eventTicks > 0) {
            eventTicks--;
            return ;
        }
        eventTicks = 8;
        lastX = event.clientX;
        lastY = event.clientY;
        let packet = new Packet(MOUSE_MOVE_ID, new Uint8Array(5));
        // set two bytes for x and two bytes for y
        let buffer = new ArrayBuffer(5);
        let data = new DataView(buffer);
        data.setUint8(0, MOUSE_MOVE_ID);
        // Instead of sending the raw x and y values, we send a value between 0 and 65535
        data.setUint16(1, event.clientX / window.innerWidth * 65535);
        data.setUint16(3, event.clientY / window.innerHeight * 65535);
        packet.data = new Uint8Array(buffer);
        socket.send(packet.data);
    });
}


socket.onmessage = function (event) {
    let arrayBuffer = event.data;
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