import { SubscribeActionTypes } from "./types";

export abstract class AAction {
    /**
     * The ID of the DOM element to add the event listener to
     */
    abstract readonly domID?: string;

    /**
     * The type of the action
     */
    abstract readonly actionType: SubscribeActionTypes;

    /**
     * The key to use in local storage
     */
    abstract readonly localStorageKey: string;

    /**
     * Write the state of the action to local storage
     * @param state The state to write to local storage
     */
    abstract writeLocalStorage(state: boolean): void;

    /**
     * Read the state of the action from local storage
     * 
     * @returns The state of the action, true if it is not in local storage
     */
    abstract readLocalStorage(): boolean;

    /**
     * Add the event listener to the document (if any)
     */
    abstract addEventListener(): void;

    /**
     * Remove the event listener from the document (if any)
     */
    abstract removeEventListener(): void;

    /**
     * Function called when the custom event is fired
     * @param event The event that was fired
     */
    abstract onEvent(event: Event): void;

    /**
     * Set the websocket to use for sending data, called when the websocket has reconnected or connected for the first time
     * @param websocket The websocket to use
     */
    abstract setWebsocket(websocket: WebSocket): void;
}