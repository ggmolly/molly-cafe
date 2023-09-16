import { APacket } from "./APacket";

export class LeitnerUpdatePacket extends APacket {
    update() {
        const element: HTMLHeadingElement = document.getElementById(this.name) as HTMLHeadingElement;
        if (!element) { return; }
        let completed: number = this.raw.getUint32(this.offset);
        this.offset += 4;
        let total: number = this.raw.getUint32(this.offset);
        element.getElementsByTagName('span')[0].innerText = `${completed}/${total}`
        element.getElementsByTagName('span')[1].innerText = `${((completed / total) * 100) | 0}%`;
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