import { APacket } from "./APacket";

const GOAL_SECONDS: number = 8 * 60 * 60; // 8 hours
const COLORS: string[] = [
    "green",
    "yellow",
    "red",
];
export class SleepTrackingPacket extends APacket {
    private convertSeconds(seconds: number): string {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds - (hours * 3600)) / 60);
        return `${hours}h${minutes.toString().padStart(2, '0')}`;
    }

    update() {
        const sleepTrackDOM: HTMLParagraphElement = document.getElementById('sleep-tracking')!! as HTMLParagraphElement;
        sleepTrackDOM.innerText = `${this.convertSeconds(this.data)} / ${this.convertSeconds(GOAL_SECONDS)}`
        // Remove all colors
        COLORS.forEach((color) => {
            sleepTrackDOM.classList.remove(color);
        });
        const percentage = this.data / GOAL_SECONDS;
        if (percentage >= 1.0) {
            sleepTrackDOM.classList.add(COLORS[0]);
        } else if (percentage > 0.925) {
            sleepTrackDOM.classList.add(COLORS[1]);
        } else {
            sleepTrackDOM.classList.add(COLORS[2]);
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