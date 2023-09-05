import { APacket } from "./APacket";

export class StrawberrySeekPacket extends APacket {
    position: number;

    constructor(data: DataView) {
        super(data);
        this.position = data.getUint32(this.offset);
        this.offset += 4;
        // EOF
    }

    formatTime(timeUs: number) {
        let time = Math.floor(timeUs / 1000000);
        let minutes = Math.floor(time / 60);
        let seconds = time % 60;
        return minutes + ":" + (seconds < 10 ? "0" : "") + seconds;
    }

    render() {}

    update() {
        window.progress = this.position;
        let element: HTMLElement | null = document.getElementById("song-time");
        if (element == null) {
            return;
        }
        if (element.innerText.includes("?")) {
            return;
        }
        document.getElementById("song-time")!!.innerText = this.formatTime(window.progress) + " / " + this.formatTime(window.length);
    }

    renderOrUpdate() {
        this.update();
    }
}