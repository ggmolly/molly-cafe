const websocket = new WebSocket('ws://localhost:50154/ws');
websocket.binaryType = 'arraybuffer';

let connectedClients: number = 0;

enum DataType {
    UINT8 = 0x00,
    UINT32 = 0x01,
    PERCENTAGE = 0x02,
    TEMPERATURE = 0x03,
    LOAD_USAGE = 0x04,
}

function interpretState(state: number): Array<string> {
    switch (state) {
        case 0x00:
            return ["DEAD", "red"]
        case 0x01:
            return ["UNHEALTHY", "yellow"]
        case 0x02:
            return ["OK", "green"]
        default:
            return ["?", "blue"]
    }
}

function getColoredSpan(packet: Packet): HTMLSpanElement {
    let span: HTMLSpanElement = document.createElement('span');
    switch (packet.datatype) {
        case DataType.UINT8: // state
            let state: Array<string> = interpretState(packet.data);
            span.innerText = state[0];
            span.classList.add(state[1]);
            break;
        case DataType.UINT32: // count
            span.innerText = packet.data.toString();
            span.classList.add('fuchsia')
            break;
        case DataType.PERCENTAGE: // percentage
            span.innerText = packet.data.toFixed(2) + '%';
            span.classList.add('fuchsia')
            break;
        case DataType.TEMPERATURE: // temperature
            span.innerText = packet.data.toFixed(2) + '°C';
            // green if < 50, yellow if < 70, red otherwise
            if (packet.data < 50) {
                span.classList.add('green');
            } else if (packet.data < 70) {
                span.classList.add('yellow');
            } else {
                span.classList.add('red');
            }
            break;
        case DataType.LOAD_USAGE: // load / usage
            span.innerText = packet.data.toFixed(2) + '%';
            // same rules as temperature
            if (packet.data < 50) {
                span.classList.add('green');
            } else if (packet.data < 60) {
                span.classList.add('yellow');
            } else if (packet.data < 80) {
                span.classList.add('orange');
            } else {
                span.classList.add('red');
            }
            break;
    }
    span.title = 'Last updated: ' + new Date().toLocaleString();
    return span;
}

function newLine(packet: Packet): HTMLHeadingElement {
    let h3: HTMLHeadingElement = document.createElement('h3');
    h3.id = "m-" + packet.id.toString();
    // pad with enough '.' to have a string of length 30
    let name: string = packet.name;
    let padLength = 30 - name.length;
    name += '.'.repeat(padLength > 0 ? padLength : 0);
    h3.innerText = name;
    // add a span for the unit
    let span: HTMLSpanElement = getColoredSpan(packet);
    // add a space between the h3's text and the span
    h3.appendChild(document.createTextNode(' ['));
    h3.appendChild(span);
    h3.appendChild(document.createTextNode(']'));
    return h3;
}

class Packet {
    category: number;
    id: number;
    datatype: DataType;
    name: string;
    data: number;

    constructor(data: DataView) {
        this.category = data.getUint8(0);
        this.id = data.getUint16(1);
        this.datatype = data.getUint8(3);
        let nameLength: number = data.getUint16(4);
        // decode name
        this.name = '';
        for (let i = 0; i < nameLength; i++) {
            this.name += String.fromCharCode(data.getUint8(6 + i));
        }
        let offset = 6 + nameLength;
        switch (this.datatype) {
            case DataType.UINT8:
                this.data = data.getUint8(offset);
                break;
            case DataType.UINT32:
                this.data = data.getUint32(offset);
                break;
            case DataType.PERCENTAGE:
            case DataType.TEMPERATURE:
            case DataType.LOAD_USAGE:
                this.data = data.getFloat32(offset);
                break;
            default:
                throw new Error('Unknown datatype');
        }
    }

    update() {
        let element: HTMLElement | null = document.getElementById("m-" + this.id.toString());
        if (!element) {
            throw new Error('Element not found');
        }
        let span: HTMLSpanElement = getColoredSpan(this); 
        let elSpan: HTMLSpanElement = element.getElementsByTagName('span')[0];
        // replace old span with new one
        element.replaceChild(span, elSpan);
    }

    render() {
        let element: HTMLElement | null = null;
        switch (this.category) {
            case 0:
                element = document.getElementById('services');
                break;
            case 1:
                element = document.getElementById('hard-resources');
                break;
            case 2:
                element = document.getElementById('soft-resources');
                break;
            case 3:
                element = document.getElementById('misc');
                break;
        }
        if (!element) {
            throw new Error('Element not found');
        }
        element.appendChild(newLine(this));
        // Sort elements by id
        let children: Array<HTMLElement> = Array.from(element.children) as Array<HTMLElement>;
        children.sort((a, b) => {
            let idA: number = parseInt(a.id.split('-')[1]);
            let idB: number = parseInt(b.id.split('-')[1]);
            return idA - idB;
        });
        // Reorder elements
        for (let i = 0; i < children.length; i++) {
            element.appendChild(children[i]);
        }
    }

    renderOrUpdate() {
        try {
            this.update();
        } catch (error) {
            this.render();
        }
    }
}

websocket.onopen = () => {
    
};

websocket.onmessage = (event) => {
    // get as DataView
    const data = new DataView(event.data);
    // check if first byte is 0xFE or 0xFF
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
            // create packet
            let packet = new Packet(data);
            packet.renderOrUpdate();
            break;
    }
}
