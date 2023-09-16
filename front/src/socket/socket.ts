import { CursorByePacket } from "../packets/CursorByePacket";
import { LeitnerUpdatePacket } from "../packets/LeitnerUpdatePacket";
import { MonitoringPacket } from "../packets/MonitoringPacket";
import { MouseActionPacket } from "../packets/MouseActionPacket";
import { PistachePacket } from "../packets/PistachePacket";
import { SchoolProjectPacket } from "../packets/SchoolProjectPacket";
import { SleepTrackingPacket } from "../packets/SleepTrackingPacket";
import { StrawberryPacket } from "../packets/StrawberryPacket";
import { StrawberrySeekPacket } from "../packets/StrawberrySeekPacket";
import { StrawberryStatePacket } from "../packets/StrawberryStatePacket";
import { AAction } from "./actions/AAction";
import { MouseMoveAction } from "./actions/MouseMoveAction";
import { SubscribeActionTypes } from "./actions/types";

const SUB_HEADER = 0x00;
const UNSUB_HEADER = 0x01;
const packetTypes: Record<number, any> = {
    0x00: MonitoringPacket,
    0x01: SchoolProjectPacket,
    0x02: PistachePacket,
    0x03: StrawberryPacket,
    0x04: StrawberrySeekPacket,
    0x05: StrawberryStatePacket,
    0x07: CursorByePacket,
    0x08: SleepTrackingPacket,
    0x09: LeitnerUpdatePacket,

    // Special case for mouse actions
    0x06: MouseActionPacket,
};

// These sessions get (re)populated when the websocket (re)connects to the server, so we need to clear them
let RECONNECT_CLEANUP: Record<string, string> = {
    "blog-posts": "",
    "school-projects": "",
    "services": "",
    "hard-resources": "",
    "soft-resources": "",
    "misc": "",
}

export class CafeSocket {
    websocket: WebSocket;
    connectedClients: number;
    subscriptions: Set<SubscribeActionTypes>;
    actions: Record<SubscribeActionTypes, AAction>;

    subscribe(actionType: SubscribeActionTypes) {
        this.actions[actionType].writeLocalStorage(true);
        this.websocket.send(new Uint8Array([SUB_HEADER, actionType]));
        this.subscriptions.add(actionType);
        this.actions[actionType].addEventListener();
    }

    unsubscribe(actionType: SubscribeActionTypes) {
        this.actions[actionType].writeLocalStorage(false);
        this.websocket.send(new Uint8Array([UNSUB_HEADER, actionType]));
        this.subscriptions.delete(actionType);
        this.actions[actionType].removeEventListener();
    }

    isSubscribed(actionType: SubscribeActionTypes): boolean {
        return this.subscriptions.has(actionType);
    }

    setupSubscriptions() {
        this.subscriptions.clear();
        // Loop through all the actions and subscribe to them if they are in local storage
        for (let actionType: SubscribeActionTypes = 0; actionType < Object.keys(this.actions).length; actionType++) {
            if (this.actions[actionType].readLocalStorage()) {
                this.subscribe(actionType);
            }
        }
    }

    setHandlers() {
        this.websocket.onmessage = (event) => {
            const data = new DataView(event.data);
            switch (data.getUint8(0)) {
                case 0xFE:
                    this.connectedClients--;
                    document.getElementById('connected-count')!!.innerText = this.connectedClients.toString();
                    break;
                case 0xFF:
                    this.connectedClients++;
                    document.getElementById('connected-count')!!.innerText = this.connectedClients.toString();
                    break;
                case 0xFD:
                    this.connectedClients = data.getUint32(1);
                    document.getElementById('connected-count')!!.innerText = this.connectedClients.toString();
                    break;
                case 0xFC: // DOMPopPacket
                    document.getElementById(new TextDecoder().decode(event.data.slice(1)))!!.remove();
                    break;
                default:
                    const target: number = new Uint8Array(event.data)[0];
                    const buff: DataView = new DataView(event.data.slice(1));
                    new packetTypes[target](buff).renderOrUpdate();
                    break;
            }
        };

        this.websocket.onclose = (event) => {
            if (event.wasClean) {
                console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
            } else {
                console.log('[close] Connection died, reconnecting in 1 second');
                setTimeout(() => {
                    this.websocket = new WebSocket('ws://localhost:50154/ws');
                    this.websocket.binaryType = 'arraybuffer';
                    this.setHandlers();
                }, 1000);
            }
        };

        this.websocket.onopen = (error) => {
            // Loop through all the actions and set the websocket
            for (let actionType: SubscribeActionTypes = 0; actionType < Object.keys(this.actions).length; actionType++) {
                this.actions[actionType].setWebsocket(this.websocket);
                this.actions[actionType].removeEventListener(); // Remove all event listeners if the socket reconnects
                document.getElementById(this.actions[actionType].domID!!)!!.addEventListener("change", (event) => {
                    if ((event.target as HTMLInputElement).checked) {
                        this.subscribe(actionType);
                    } else {
                        this.unsubscribe(actionType);
                    }
                });
            }
            // Clear sections that are populated by the websocket
            for (let section in RECONNECT_CLEANUP) {
                document.getElementById(section)!!.innerHTML = RECONNECT_CLEANUP[section];
            }
            this.setupSubscriptions();
        };
    }

    constructor() {
        // Keep a copy of each original sections, so we can restore them when the websocket reconnects
        for (let section in RECONNECT_CLEANUP) {
            RECONNECT_CLEANUP[section] = document.getElementById(section)!!.innerHTML;
        }
        this.websocket = new WebSocket('ws://localhost:50154/ws');
        this.websocket.binaryType = 'arraybuffer';
        this.connectedClients = 0;
        this.subscriptions = new Set<SubscribeActionTypes>();
        this.setHandlers();
        this.actions = {
            [SubscribeActionTypes.S_CURSOR]: new MouseMoveAction(),
        };

        // For each actions, check if the user is subscribed to it, and if so, check the checkbox
        for (let actionType: SubscribeActionTypes = 0; actionType < Object.keys(this.actions).length; actionType++) {
            if (this.actions[actionType].readLocalStorage()) {
                document.getElementById(this.actions[actionType].domID!!)!!.setAttribute("checked", "");
            } else {
                document.getElementById(this.actions[actionType].domID!!)!!.removeAttribute("checked");
            }
        }
    }
}