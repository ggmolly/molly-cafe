import { APacket } from "./APacket";

const MASK_WIP      = 0b00000001;
const MASK_GRADING  = 0b00000010;

export class SchoolProjectPacket extends APacket {
    href:        string;
    description: string;
    wip:         boolean;
    grading:      boolean;
    grade:       number;

    constructor(data: DataView) {
        super(data);
        const stateMask = this.raw.getUint8(this.offset++);
        this.wip = (stateMask & MASK_WIP) != 0;
        this.grading = (stateMask & MASK_GRADING) != 0;
        this.grade = this.raw.getUint8(this.offset++);
        const descriptionLength = this.raw.getUint8(this.offset++);
        this.description = new TextDecoder().decode(this.raw.buffer.slice(this.offset, this.offset + descriptionLength));
        this.offset += descriptionLength;
        const hrefLength = this.raw.getUint8(this.offset++);
        this.href = new TextDecoder().decode(this.raw.buffer.slice(this.offset, this.offset + hrefLength));
        this.offset += hrefLength;
    }

    // TODO: add a way to differenciate between github / blog posts
    constructNameCell(): HTMLTableCellElement {
        const nameCell: HTMLTableCellElement = document.createElement('td');
        if (this.href) {
            const link: HTMLAnchorElement = document.createElement('a');
            link.href = this.href;
            link.innerText = this.name;
            nameCell.appendChild(link);
        } else {
            nameCell.innerText = this.name;
        }
        return nameCell;
    }

    constructDescriptionCell(): HTMLTableCellElement {
        const descriptionCell: HTMLTableCellElement = document.createElement('td');
        descriptionCell.innerText = this.description;
        return descriptionCell;
    }

    constructGradeCell(): HTMLTableCellElement {
        const gradeCell: HTMLTableCellElement = document.createElement('td');
        if (this.wip) {
            gradeCell.innerText = "WIP";
            gradeCell.classList.add('grade-wip');
            return gradeCell;
        } else if (this.grading) {
            gradeCell.innerText = "Not graded yet";
            gradeCell.classList.add('grade-wip');
            return gradeCell;
        } else {
            gradeCell.innerText = this.grade.toString();
        }

        // Grade color
        if (this.grade < 100) {
            gradeCell.classList.add('grade-fail');
        } else if (this.grade < 125) {
            gradeCell.classList.add('grade-ok');
        } else if (this.grade == 125) {
            gradeCell.classList.add('grade-max-bonus');
        }

        return gradeCell;
    }

    update() {
        let row: HTMLElement | HTMLTableRowElement | null = document.getElementById("sp-" + this.id.toString());
        if (!row) {
            throw new Error('Element not found');
        }
        row = row as HTMLTableRowElement;
        row.children[0].replaceWith(this.constructNameCell());
        row.children[1].replaceWith(this.constructDescriptionCell());
        row.children[2].replaceWith(this.constructGradeCell());
    }

    render() {
        const table: HTMLTableElement = document.getElementById("school-projects")!! as HTMLTableElement;
        const row: HTMLTableRowElement = document.createElement('tr');
        row.id = "sp-" + this.id.toString();
        row.appendChild(this.constructNameCell());
        row.appendChild(this.constructDescriptionCell());
        row.appendChild(this.constructGradeCell());
        table.appendChild(row);
    }

    sort() {
        const table: HTMLTableElement = document.getElementById("school-projects")!! as HTMLTableElement;
        const rows: HTMLCollectionOf<HTMLTableRowElement> = table.rows;
        const rowsArray: HTMLTableRowElement[] = Array.from(rows);
        rowsArray.sort((a, b) => {
            const aName = a.children[0].innerHTML;
            const bName = b.children[0].innerHTML;
            return aName.localeCompare(bName);
        });
        rowsArray.forEach((row) => {
            table.appendChild(row);
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