export class MouseActionPacket {
    percentageX: number;
    percentageY: number;
    clientId: number;

    constructor(data: DataView) {
        this.percentageX = data.getFloat32(0);
        this.percentageY = data.getFloat32(4);
        this.clientId = data.getUint32(8);
    }

    render() {
        let cursor: HTMLImageElement = document.createElement("img");
        cursor.src = "/assets/cursor.png";
        cursor.id = "cursor-" + this.clientId;
        cursor.style.position = "fixed";
        // remove the scroll offset
        cursor.style.left = (this.percentageX * window.innerWidth - window.scrollX).toString() + "px";
        cursor.style.top = (this.percentageY * window.innerHeight - window.scrollY).toString() + "px";
        document.getElementById("cursors")!!.appendChild(cursor);
    }

    update() {
        let cursor: HTMLImageElement = document.getElementById("cursor-" + this.clientId) as HTMLImageElement;
        if (cursor === null) {
            throw new Error("Cursor not found");
        }
        // restore the scroll offset
        cursor.style.left = (this.percentageX * window.innerWidth - window.scrollX).toString() + "px";
        cursor.style.top = (this.percentageY * window.innerHeight - window.scrollY).toString() + "px";
    }

    renderOrUpdate() {
        try {
            this.update();
        } catch (e) {
            this.render();
        }
    }
}