import { APacket } from "./APacket";

let localTimeInterval: any | undefined = undefined;
const formatter: Intl.DateTimeFormat = new Intl.DateTimeFormat(undefined, {
    timeZone: "Europe/Paris",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false
});

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
        window.s_Weather.timeToSunrise = data.getUint32(this.offset) * 1000;
        this.offset += 4;
        // Byte 15-19 = timeToSunset (unix time)
        window.s_Weather.timeToSunset = data.getUint32(this.offset) * 1000;
        this.offset += 4;
        // Byte 20-24 = currentTime (unix time)
        window.s_Weather.currentTime = data.getUint32(this.offset) * 1000;
        this.offset += 4;
        // Byte 25-26 = lastRainTime (seconds)
        window.s_Weather.lastRainTime = data.getInt16(this.offset);
        this.offset += 2;
        // Byte 27 = currentCondition (string length)
        let currentConditionLength: number = data.getUint8(this.offset++);
        // Byte 28-... = currentCondition (string)
        window.s_Weather.currentCondition = new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + currentConditionLength));
        if (window.s_Weather.onCloudinessChange !== undefined) {
            window.s_Weather.onCloudinessChange(window.s_Weather.cloudiness);
        }
        if (window.s_Weather.onRainIntensityChange !== undefined) {
            window.s_Weather.onRainIntensityChange(window.s_Weather.rainIntensity);
        }
        if (window.s_Weather.onWindSpeedChange !== undefined) {
            window.s_Weather.onWindSpeedChange(window.s_Weather.windSpeed);
        }
    }

    timeFormatting(time: number): string {
        // hh:mm
        let hours: string = Math.floor(time / 3600).toString();
        let minutes: string = Math.floor((time % 3600) / 60).toString();
        return hours.padStart(2, "0") + ":" + minutes.padStart(2, "0");
    }

    updateLocalTime() {
        // Render the time using Intl.DateTimeFormat
        (document.querySelector("#w-current-time > span.value") as HTMLElement).innerText = formatter.format(new Date());
        window.s_Weather.currentTime += 1000;
    }

    update() {
        // Update Sunrise/Sunset
        const sunrise = new Date(window.s_Weather.timeToSunrise);
        const sunset = new Date(window.s_Weather.timeToSunset);
        (document.querySelector("#w-sunset > span.value") as HTMLElement).innerText = this.timeFormatting(sunset.getHours() * 3600 + sunset.getMinutes() * 60);
        (document.querySelector("#w-sunrise > span.value") as HTMLElement).innerText = this.timeFormatting(sunrise.getHours() * 3600 + sunrise.getMinutes() * 60);
        this.updateLocalTime();

        // Set interval to update local time
        if (localTimeInterval !== undefined) {
            clearInterval(localTimeInterval);
        }
        localTimeInterval = setInterval(this.updateLocalTime, 1000);

        // Update cloudiness
        (document.querySelector("#w-cloudiness > span.value") as HTMLElement).innerText = window.s_Weather.cloudiness.toString() + "%";

        // Update temperature
        (document.querySelector("#w-temperature > span.value") as HTMLElement).innerText = window.s_Weather.feelsLike.toString() + "°C";

        // Update humidity
        (document.querySelector("#w-humidity > span.value") as HTMLElement).innerText = window.s_Weather.humidity.toString() + "%";

        // Update current condition
        (document.querySelector("#w-current-cond") as HTMLElement).innerText = window.s_Weather.currentCondition;

        // Update wind speed (m/s)
        (document.querySelector("#w-wind-speed > span.value") as HTMLElement).innerText = window.s_Weather.windSpeed.toString() + "m/s";
    }

    render() {}

    renderOrUpdate() {
        this.update();
    }
}