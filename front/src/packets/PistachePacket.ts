import { APacket } from "./APacket";

export class PistachePacket extends APacket {
    creationDate: Date;
    href: string;

    constructor(data: DataView) {
        super(data);
        // Href
        const hrefLength = data.getUint16(this.offset);
        this.offset += 2;
        this.href = new TextDecoder().decode(data.buffer.slice(this.offset, this.offset + hrefLength));
        this.offset += hrefLength;
        // Creation date
        const creationDate: number = data.getUint32(this.offset);
        this.creationDate = new Date(creationDate * 1000);
        this.offset += 4;
        // EOF
    }

    ordinal(n: number): string {
        const suffix = ["th", "st", "nd", "rd"];
        const idx = n % 100;
        return n + (suffix[(idx - 20) % 10] || suffix[idx] || suffix[0]);
    }

    formatDate(): string {
        const days: string[] = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
        const months: string[] = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
        const day: string = days[this.creationDate.getDay()];
        const date: string = this.ordinal(this.creationDate.getDate());
        const month: string = months[this.creationDate.getMonth()];
        const year: number = this.creationDate.getFullYear();
        return day + ", " + date + " " + month + " " + year;
    }

    dateSpan(): HTMLSpanElement {
        const dateSpan: HTMLSpanElement = document.createElement('span');
        dateSpan.innerText = this.formatDate();
        return dateSpan;
    }

    createTitleLink() {
        const titleLink: HTMLAnchorElement = document.createElement('a');
        titleLink.href = "/pistache/" + this.href;
        titleLink.innerText = "➜ " + this.name;
        titleLink.classList.add("green");
        return titleLink;
    }

    update() {
        let listElement: HTMLElement | HTMLLIElement | null = document.getElementById("pb-" + this.id);
        if (!listElement) {
            throw new Error("Element not found");
        }
        listElement = listElement as HTMLLIElement;
        listElement.children[0].replaceWith(this.createTitleLink());
        listElement.children[1].replaceWith(this.dateSpan());
    }

    render() {
        const list: HTMLUListElement = document.getElementById("blog-posts")!! as HTMLUListElement;
        const listElement: HTMLLIElement = document.createElement('li');
        listElement.id = "pb-" + this.id;
        listElement.appendChild(this.createTitleLink());
        listElement.appendChild(document.createTextNode(" — "));
        listElement.appendChild(this.dateSpan());
        listElement.setAttribute("data-timestamp", this.creationDate.getTime().toString());
        list.appendChild(listElement);
    }

    sort() {
        const list: HTMLUListElement = document.getElementById("blog-posts")!! as HTMLUListElement;
        const items: HTMLLIElement[] = Array.from(list.children) as HTMLLIElement[];
        items.sort((a: HTMLLIElement, b: HTMLLIElement) => {
            const aTimestamp: number = parseInt(a.getAttribute("data-timestamp")!!);
            const bTimestamp: number = parseInt(b.getAttribute("data-timestamp")!!);
            return bTimestamp - aTimestamp;
        });
        items.forEach((item: HTMLLIElement) => {
            list.appendChild(item);
        });
    }

    renderOrUpdate() {
        try {
            this.update();
        } catch (error) {
            this.render();
        }
        this.sort();
    }
}