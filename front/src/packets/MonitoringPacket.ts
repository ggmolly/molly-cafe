import { APacket } from "./APacket";
import { DataType } from "./DataTypes";

export class MonitoringPacket extends APacket {
    interpretState(): Array<string> {
        switch (this.data) {
            case 0x00:
                return ["DEAD", "red"]
            case 0x01:
                return ["UNHEALTHY", "yellow"]
            case 0x02:
                return ["OK", "green"]
            default:
                return ["?", "blue"]
        }
    }

    update() {
        let element: HTMLElement | null = document.getElementById("m-" + this.id.toString());
        if (!element) {
            throw new Error('Element not found');
        }
        let span: HTMLSpanElement = this.getColoredSpan(); 
        let elSpan: HTMLSpanElement = element.getElementsByTagName('span')[0];
        // replace old span with new one
        element.replaceChild(span, elSpan);
    }

    getColoredSpan(): HTMLSpanElement {
        let span: HTMLSpanElement = document.createElement('span');
        switch (this.datatype) {
            case DataType.UINT8: // state
                let state: Array<string> = this.interpretState();
                span.innerText = state[0];
                span.classList.add(state[1]);
                break;
            case DataType.UINT32: // count
                span.innerText = this.data.toString();
                span.classList.add('fuchsia')
                break;
            case DataType.PERCENTAGE: // percentage
                span.innerText = this.data.toFixed(2) + '%';
                span.classList.add('fuchsia')
                break;
            case DataType.TEMPERATURE: // temperature
                span.innerText = this.data.toFixed(2) + '°C';
                // green if < 50, yellow if < 70, red otherwise
                if (this.data < 50) {
                    span.classList.add('green');
                } else if (this.data < 70) {
                    span.classList.add('yellow');
                } else {
                    span.classList.add('red');
                }
                break;
            case DataType.LOAD_USAGE: // load / usage
                span.innerText = this.data.toFixed(2) + '%';
                // green if < 50, yellow if < 70, red otherwise
                if (this.data < 50) {
                    span.classList.add('green');
                } else if (this.data < 70) {
                    span.classList.add('yellow');
                } else {
                    span.classList.add('red');
                }
                break;
        }
        return span;
    }

    newLine(): HTMLHeadingElement {
        let h4: HTMLHeadingElement = document.createElement('h4');
        h4.id = "m-" + this.id.toString();
        // pad with enough '.' to have a string of length 30
        let name: string = this.name;
        let padLength = 30 - name.length;
        name += '.'.repeat(padLength > 0 ? padLength : 0);
        h4.innerText = name;
        // add a span for the unit
        let span: HTMLSpanElement = this.getColoredSpan();
        // add a space between the h4's text and the span
        h4.appendChild(document.createTextNode(' ['));
        h4.appendChild(span);
        h4.appendChild(document.createTextNode(']'));
        return h4;
    }

    render() {
        let element: HTMLElement | null = null;
        switch (this.category) {
            case 0:
                element = document.getElementById('services');
                break;
            case 1:
                element = document.getElementById('hard-resources');
                break;
            case 2:
                element = document.getElementById('soft-resources');
                break;
            case 3:
                element = document.getElementById('misc');
                break;
        }
        if (!element) {
            throw new Error('Element not found');
        }
        element.appendChild(this.newLine());
        // Sort elements by id
        let children: Array<HTMLElement> = Array.from(element.children) as Array<HTMLElement>;
        children.sort((a, b) => {
            let idA: number = parseInt(a.id.split('-')[1]);
            let idB: number = parseInt(b.id.split('-')[1]);
            return idA - idB;
        });
        // Reorder elements
        for (let child of children) {
            element.appendChild(child);
        }
    }

    renderOrUpdate() {
        try {
            this.update();
        } catch (error) {
            this.render();
        }
    }
}