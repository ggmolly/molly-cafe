import { AAction } from "./AAction";
import { SubscribeActionTypes } from "./types";

let lastPacketAt: number = 0;
const MIN_DELAY: number = 13;
export class MouseMoveAction extends AAction {
    readonly domID?: string = "show-cursors";
    readonly actionType: SubscribeActionTypes = SubscribeActionTypes.S_CURSOR;
    readonly localStorageKey: string = "cursor_sub";
    websocket?: WebSocket;

    writeLocalStorage(state: boolean): void {
        if (!state) {
            document.getElementById("cursors")!!.innerHTML = "";
        }
        localStorage.setItem(this.localStorageKey, state ? 'true' : 'false');
    }

    readLocalStorage(): boolean {
        return localStorage.getItem(this.localStorageKey) === 'true' || localStorage.getItem(this.localStorageKey) === null;
    }

    addEventListener(): void {
        this.onEvent = this.onEvent.bind(this); // Re-bind the function to the class, so the websocket can be accessed when changed
        document.addEventListener("mousemove", this.onEvent);
    }

    removeEventListener(): void {
        document.removeEventListener("mousemove", this.onEvent);
    }

    onEvent(event: MouseEvent): void {
        // Get the x and y position of the mouse (the DOM can be scrolled!!!)
        let x: number = event.clientX + window.scrollX;
        let y: number = event.clientY + window.scrollY;
        x = x / window.innerWidth;
        y = y / window.innerHeight;
        // Only send the data if it has been 100ms since the last packet
        if (Date.now() - lastPacketAt < MIN_DELAY) {
            return;
        }
        // Send the raw data to test quickly
        let buff: ArrayBuffer = new ArrayBuffer(9);
        let data: DataView = new DataView(buff);
        // Set first byte to 0x06 (cursor pos)
        data.setUint8(0, 0x06);
        // Set the next 4 bytes to the x position
        data.setFloat32(1, x);
        // Set the next 4 bytes to the y position
        data.setFloat32(5, y);
        if (this.websocket !== undefined && this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(buff);
        }
        lastPacketAt = Date.now();
    }

    setWebsocket(websocket: WebSocket): void {
        this.websocket = websocket;
    }

    constructor() {
        super();
    }
}