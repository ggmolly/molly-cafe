import { getAssets, loadAssetsType } from "../../assets";
import { Point } from "../../interfaces/point.interface";
import { Velocity } from "../../interfaces/velocity.interface";
import { AssetType } from "../../types";
import { ADrawable } from "../bases/ADrawable";
import { AMovable } from "../bases/AMovable";

const N_RAINDROPS = 8;

function checkPosition(pos: Point, rect: DOMRect): Point {
    if (pos.x + window.scrollX - rect.width < 10) { // If the raindrop is too close to the left edge
        pos.x += rect.width - 5;
    } else if (pos.x + window.scrollX + rect.width > window.innerWidth - 10) { // If the raindrop is too close to the right edge
        pos.x -= rect.width + 5;
    }
    return pos;
}

function randomCard(): HTMLElement {
    let cards: Array<HTMLElement> = Array.from(document.querySelectorAll("div.card:last-child,div.side-card:last-child"));
    return cards[Math.floor(Math.random() * cards.length)];
}

class ResidualRaindrop extends AMovable {
    private _parentDom: HTMLElement;
    private _waitingForDiabling: boolean = false;
    constructor(
        sprite: HTMLImageElement,
        context: CanvasRenderingContext2D,
        initialVelocity: Velocity,
    ) {
        let parentCard: HTMLElement = randomCard();
        // Get rect of parent card
        let rect: DOMRect = parentCard.getBoundingClientRect();
        let pos = {
            x: rect.x + window.scrollX + rect.width * Math.random(),
            y: rect.y + window.scrollY + rect.height - sprite.height,
        };
        pos = checkPosition(pos, rect);
        super(sprite, context, pos, "ResidualRaindrop", initialVelocity);
        this._parentDom = parentCard;
        this.enabled = this.conditionsMet();
    }

    tick() {
        if (this.conditionsMet() && !this.enabled) {
            this.resetPosition();
            this.enable();
        } else if (!this.conditionsMet() && this.enabled) {
            this.disable();
        }
        if (!this.enabled) return;
        if (this.pos.y > this.context.canvas.height) {
            return this.resetPosition();
        }
        let rect: DOMRect = this._parentDom.getBoundingClientRect();
        if (this.pos.y > rect.y + window.scrollY + rect.height - this.sprite.height / (Math.random() * 2 + 2)) {
            this.velocity.y = Math.random() * 2.5 + 0.5;
        } else {
            this.velocity.y = Math.random() * 0.014 + 0.003;
        }
    }

    resetPosition() {
        if (this._waitingForDiabling) {
            this.enabled = false;
            this._waitingForDiabling = false;
            return;
        }
        // Pick another parent card
        this._parentDom = randomCard();
        // Randomize position
        let rect: DOMRect = this._parentDom.getBoundingClientRect();
        this.pos.x = rect.x + window.scrollX + rect.width * Math.random();
        this.pos.y = rect.y + window.scrollY + rect.height - (this.sprite.height * 1.30);
        this.pos = checkPosition(this.pos, rect);
    }

    public disable(): void {
        this._waitingForDiabling = true;
    }

    conditionsMet(): boolean {
        // Return true if it is raining or it rained in the last hour
        return window.s_Weather.currentTime !== 0 &&
            (
                window.s_Weather.rainIntensity > 0 ||
                (window.s_Weather.lastRainTime < 1800 && window.s_Weather.lastRainTime > 0)
            );
    }
}

export async function residualRaindropInit(ctx: CanvasRenderingContext2D): Promise<Array<ADrawable>> {
    let raindrops: Array<ResidualRaindrop> = new Array<ResidualRaindrop>();
    return loadAssetsType(AssetType.RAINDROP).then(() => {
        let raindropSprite = getAssets(AssetType.RAINDROP);
        for (let i = 0; i < N_RAINDROPS; i++) {
            raindrops.push(
                new ResidualRaindrop(
                    raindropSprite[0],
                    ctx,
                    { x: 0, y: Math.random() * 0.003 + 0.002 }
                )
            );
        }
        return raindrops;
    });
}