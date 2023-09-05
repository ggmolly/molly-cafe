import { APacket } from "./APacket";

export class StrawberryPacket extends APacket {
    title: string;
    artists: Array<string>;
    cover: string;
    length: number; // in microseconds

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

        // Progress
        window.progress = data.getUint32(this.offset) - 1000000; // -1s to compensate the forced timeTick call
        this.offset += 4;
        // EOF

        // Replace interval for time ticking
        if (window.interval) {
            clearInterval(window.interval);
            window.interval = null;
        }
        window.length = this.length;
        window.interval = setInterval(() => {
            this.timeTick();
        }, 1000);
        this.timeTick();
    }

    timeTick() {
        if (document.getElementById("strawberry-disc")!!.classList.contains("spinning")) {
            window.progress += 1000000;
        }
        if (window.progress > this.length) {
            window.progress = this.length;
        }
        document.getElementById("song-time")!!.innerText = this.formatTime(window.progress) + " / " + this.formatTime(this.length);
    }

    formatTime(timeUs: number) {
        let time = Math.floor(timeUs / 1000000);
        let minutes = Math.floor(time / 60);
        let seconds = time % 60;
        return minutes + ":" + (seconds < 10 ? "0" : "") + seconds;
    }

    render() {}

    update() {
        let titleElement: HTMLElement = document.getElementById("song-title") as HTMLElement;
        let artistElement: HTMLElement = document.getElementById("song-artist") as HTMLElement;
        document.querySelector("#strawberry-disc > img")!!.setAttribute("src", this.cover);
        titleElement.innerText = this.title;
        artistElement.innerText = this.artists.join(", ");
        if (window.length < 0) {
            return
        }
    }

    renderOrUpdate() {
        this.update();
    }
}