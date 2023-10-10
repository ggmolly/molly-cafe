import { Rect } from "./interfaces/rect.interface";
import { Weather } from "./interfaces/weather.interface";
import { ADrawable } from "./sirius/objects/bases/ADrawable";

export { };

declare global {
    interface Window {
        interval: any;
        progress: number;
        tableRect: Rect;
        s_Weather: Weather;
        s_Objects: ADrawable[];
    }
}