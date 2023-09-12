import { Rect } from "./interfaces/rect.interface";

export { };

declare global {
    interface Window {
        interval: any;
        progress: number;
        tableRect: Rect;
    }
}