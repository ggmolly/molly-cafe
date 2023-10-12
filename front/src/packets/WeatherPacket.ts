import { APacket } from "./APacket";

export class WeatherPacket extends APacket {
    constructor(data: DataView) {
        super(data);
        // Byte 1 = windSpeed (0-255)
        window.s_Weather.windSpeed = data.getUint8(this.offset++);
        // Byte 2 = rainIntensity (0-255)
        window.s_Weather.rainIntensity = data.getUint8(this.offset++);
        // Byte 3 = cloudiness (0-100)
        window.s_Weather.cloudiness = data.getUint8(this.offset++);
        // Byte 4 = temperature (real part)
        window.s_Weather.temperature = data.getUint8(this.offset++);
        // Byte 5 = temperature (fractional part)
        window.s_Weather.temperature += data.getUint8(this.offset++) / 100;
        // Byte 6 = feelsLike (real part)
        window.s_Weather.feelsLike = data.getUint8(this.offset++);
        // Byte 7 = feelsLike (fractional part)
        window.s_Weather.feelsLike += data.getUint8(this.offset++) / 100;
        // Byte 8 = humidity (real part)
        window.s_Weather.humidity = data.getUint8(this.offset++);
        // Byte 9 = humidity (fractional part)
        window.s_Weather.humidity += data.getUint8(this.offset++) / 100;
        // Byte 10-14 = timeToSunrise (unix time)
        window.s_Weather.timeToSunrise = data.getUint32(this.offset);
        this.offset += 4;
        // Byte 15-19 = timeToSunset (unix time)
        window.s_Weather.timeToSunset = data.getUint32(this.offset);
        this.offset += 4;
        // Byte 20-24 = currentTime (unix time)
        window.s_Weather.currentTime = data.getUint32(this.offset);
        this.offset += 4;
        // Byte 25 = currentCondition (string length)
        let currentConditionLength: number = data.getUint8(this.offset++);
        // Byte 26-... = currentCondition (string)
        window.s_Weather.currentCondition = new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + currentConditionLength));
        if (window.s_Weather.onCloudinessChange !== undefined) {
            window.s_Weather.onCloudinessChange(window.s_Weather.cloudiness);
        }
        if (window.s_Weather.onRainIntensityChange !== undefined) {
            window.s_Weather.onRainIntensityChange(window.s_Weather.rainIntensity);
        }
    }

    update() {}

    render() {}

    renderOrUpdate() {
    }
}