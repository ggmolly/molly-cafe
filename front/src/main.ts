import { MonitoringPacket } from "./packets/MonitoringPacket";

const websocket = new WebSocket('ws://localhost:50154/ws');
websocket.binaryType = 'arraybuffer';

let connectedClients: number = 0;

websocket.onopen = () => {
    
};

const packetTypes: Record<number, any> = {
    0x00: MonitoringPacket,
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
        default:
            const target: number = new Uint8Array(event.data)[0];
            const buff: DataView = new DataView(event.data.slice(1));
            new packetTypes[target](buff).renderOrUpdate();
            break;
    }
}
