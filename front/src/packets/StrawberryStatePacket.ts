import { APacket } from "./APacket";

export class StrawberryStatePacket extends APacket {
    playing: boolean;

    constructor(data: DataView) {
        super(data);
        this.playing = data.getUint8(this.offset) == 1;
        this.offset += 1;
        // EOF
    }

    render() {}

    update() {
        // Stop / start spinning the disc
        if (this.playing) {
            document.getElementById("strawberry-disc")!!.classList.add("spinning");
        } else {
            document.getElementById("strawberry-disc")!!.classList.remove("spinning");
        }
    }

    renderOrUpdate() {
        this.update();
    }
}