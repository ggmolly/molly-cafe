export interface Weather {
    /**
     * Wind speed normalized to 0-100 range to avoid floats
     */
    windSpeed: number;

    /**
     * Rain intensity normalized to 0-100 range to avoid floats
     */
    rainIntensity: number;
    timeOfDay: number;
    temperature: number;
}