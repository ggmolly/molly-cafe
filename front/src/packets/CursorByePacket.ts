import { APacket } from "./APacket";

export class CursorByePacket extends APacket {
    update() {}

    render() {}

    renderOrUpdate() {
        document.getElementById("cursor-" + this.data)?.remove();
    }
}