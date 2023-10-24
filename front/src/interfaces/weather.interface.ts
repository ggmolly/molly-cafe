export interface Weather {
    /**
     * Wind speed normalized to 0-255 range to avoid floats
     */
    windSpeed: number;

    /**
     * Rain intensity normalized to 0-255 range to avoid floats
     */
    rainIntensity: number;

    /**
     * Cloudiness from 0 to 100
     */
    cloudiness: number;

    /**
     * Temperature in celsius
     */
    temperature: number;
    feelsLike: number;

    /**
     * Humidity in percent
     */
    humidity: number;

    /**
     * Time to sunrise in unix time
     */
    timeToSunrise: number;

    /**
     * Time to sunset in unix time
     */
    timeToSunset: number;

    /**
     * Time of day in unix time
     */
    currentTime: number;

    /**
     * Weather condition
     */
    currentCondition: string;

    /**
     * Last time it rained in seconds
     */
    lastRainTime: number;

    // --- Intervals ---

    /**
     * Interval for incrementing 
     */
    lastRainTimeUpdate: number | undefined;

    // --- Callback functions ---

    /**
     * Callback for cloudiness change
     */
    onCloudinessChange: (newCloudiness: number) => void | undefined;

    /**
     * Callback for rain intensity change
     */
    onRainIntensityChange: (newRainIntensity: number) => void | undefined;

    /**
     * Callback for wind speed change
     */
    onWindSpeedChange: (newWindSpeed: number) => void | undefined;
}