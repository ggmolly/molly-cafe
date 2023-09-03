import { MonitoringPacket } from "./packets/MonitoringPacket";
import { PistachePacket } from "./packets/PistachePacket";
import { SchoolProjectPacket } from "./packets/SchoolProjectPacket";

const websocket = new WebSocket('ws://localhost:50154/ws');
websocket.binaryType = 'arraybuffer';

let connectedClients: number = 0;

websocket.onopen = () => {
    
};

const packetTypes: Record<number, any> = {
    0x00: MonitoringPacket,
    0x01: SchoolProjectPacket,
    0x02: PistachePacket,
}

websocket.onmessage = (event) => {
    const data = new DataView(event.data);
    switch (data.getUint8(0)) {
        case 0xFE:
            connectedClients--;
            document.getElementById('connected-count')!!.innerText = connectedClients.toString();
            break;
        case 0xFF:
            connectedClients++;
            document.getElementById('connected-count')!!.innerText = connectedClients.toString();
            break;
        case 0xFD:
            connectedClients = data.getUint32(1);
            document.getElementById('connected-count')!!.innerText = connectedClients.toString();
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
}
