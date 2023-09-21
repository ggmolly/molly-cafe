import { Rect } from "./interfaces/rect.interface";
import { ADrawable } from "./sirius/objects/bases/ADrawable";

export { };

declare global {
    interface Window {
        interval: any;
        progress: number;
        tableRect: Rect;
        windSpeed: number;
        s_Objects: ADrawable[];
    }
}