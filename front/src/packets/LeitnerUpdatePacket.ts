import { APacket } from "./APacket";

const LEITNER_PROGRESS_UPDATE = 0x00;
const LEITNER_STREAK_UPDATE = 0x01;

export class LeitnerUpdatePacket extends APacket {
    update() {
        let element: HTMLElement | HTMLHeadingElement | HTMLSpanElement | null = null;
        switch (this.category) {
            case LEITNER_PROGRESS_UPDATE:
                element = document.getElementById(this.name) as HTMLHeadingElement;
                if (!element) { return; }
                let completed: number = this.raw.getUint32(this.offset);
                this.offset += 4;
                let total: number = this.raw.getUint32(this.offset);
                element.getElementsByTagName('span')[0].innerText = `${completed}/${total}`
                element.getElementsByTagName('span')[1].innerText = `${((completed / total) * 100) | 0}%`;
                break;
            case LEITNER_STREAK_UPDATE:
                element = document.getElementById("L_streak") as HTMLSpanElement;
                if (!element) { return; }
                element.innerText = this.data.toString();
        }
    }

    render() {}

    renderOrUpdate() {
        try {
            this.update();
        } catch (error) {
            this.render();
        }
    }
}