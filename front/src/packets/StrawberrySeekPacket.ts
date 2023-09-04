import { APacket } from "./APacket";

export class StrawberrySeekPacket extends APacket {
    position: number;

    constructor(data: DataView) {
        super(data);
        this.position = data.getUint32(this.offset);
        this.offset += 4;
        // EOF
    }

    render() {}

    update() {}

    renderOrUpdate() {
        throw new Error("Method not implemented.");
    }
}