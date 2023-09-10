export class MouseActionPacket {
    percentageX: number;
    percentageY: number;
    clientId: number;

    constructor(data: DataView) {
        this.percentageX = data.getFloat32(0);
        this.percentageY = data.getFloat32(4);
        this.clientId = data.getUint32(8);
    }

    interpolate(from: number[], to: number[], steps: number): void {
        let cursor: HTMLImageElement = document.getElementById("cursor-" + this.clientId) as HTMLImageElement;
        if (cursor === null) {
            throw new Error("Cursor not found");
        }
        let currentStep = 0;
        let x: number, y: number;
        const stepX = (to[0] - from[0]) / steps;
        const stepY = (to[1] - from[1]) / steps;
        const animate = () => {
            x = from[0] + stepX * currentStep;
            y = from[1] + stepY * currentStep;
            cursor.style.left = x.toString() + "px";
            cursor.style.top = y.toString() + "px";
            currentStep++;
            if (currentStep <= steps) {
                requestAnimationFrame(animate);
            }
        };
        requestAnimationFrame(animate);
    }

    render() {
        let cursor: HTMLImageElement = document.createElement("img");
        cursor.src = "/assets/cursor.png";
        cursor.id = "cursor-" + this.clientId;
        cursor.style.position = "absolute";
        cursor.style.left = (this.percentageX * window.innerWidth - window.scrollX).toString() + "px";
        cursor.style.top = (this.percentageY * window.innerHeight - window.scrollY).toString() + "px";
        document.getElementById("cursors")!!.appendChild(cursor);
    }

    update() {
        let cursor: HTMLImageElement = document.getElementById("cursor-" + this.clientId) as HTMLImageElement;
        if (cursor === null) {
            throw new Error("Cursor not found");
        }
        let x: number = (this.percentageX - window.scrollX) * window.innerWidth;
        let y: number = (this.percentageY - window.scrollY) * window.innerHeight;
        this.interpolate([parseInt(cursor.style.left), parseInt(cursor.style.top)], [x, y], 8);
    }

    renderOrUpdate() {
        try {
            this.update();
        } catch (e) {
            this.render();
        }
    }
}