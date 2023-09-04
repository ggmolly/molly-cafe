import { APacket } from "./APacket";

export class StrawberryPacket extends APacket {
    title: string;
    artists: Array<string>;
    cover: string;
    length: number; // in microseconds
    interval: any; // used to update the time
    progress: number = 0; // in microseconds

    constructor(data: DataView) {
        super(data);
        // Title parsing
        const titleLength = data.getUint16(this.offset);
        this.offset += 2;
        this.title = new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + titleLength));
        this.offset += titleLength;
        
        // Artists parsing
        const artistsLength = data.getUint16(this.offset);
        this.offset += 2;
        this.artists = new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + artistsLength)).split(",");
        this.offset += artistsLength;

        // Cover length
        const coverLength = data.getUint32(this.offset);
        this.offset += 4;
        // base64 decode cover
        this.cover = "data:image/jpg;base64," + new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + coverLength));
        this.offset += coverLength;

        // Length
        this.length = data.getUint32(this.offset);
        this.offset += 4;
        // EOF

        this.timeTick();
        this.interval = setInterval(() => {
            this.timeTick();
        }, 1000);
    }

    // TODO: Store the progress in the HTML element, to properly manage the seek packet
    timeTick() {
        if (!document.getElementById("strawberry-disc")!!.classList.contains("spinning")) {
            return;
        }
        this.progress += 1000000;
        if (this.progress > this.length) {
            this.progress = this.length;
        }
        document.getElementById("song-time")!!.innerText = this.formatTime(this.progress) + " / " + this.formatTime(this.length);
    }

    formatTime(timeUs: number) {
        let time = Math.floor(timeUs / 1000000);
        let minutes = Math.floor(time / 60);
        let seconds = time % 60;
        return minutes + ":" + (seconds < 10 ? "0" : "") + seconds;
    }

    render() {}

    update() {
        document.querySelector("#strawberry-disc > img")!!.setAttribute("src", this.cover);
        document.getElementById("song-title")!!.innerText = this.title;
        document.getElementById("song-artist")!!.innerText = this.artists.join(", ");
    }

    renderOrUpdate() {
        this.update();
    }
}