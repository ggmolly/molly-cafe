import { DataType } from "./DataTypes";

export abstract class APacket {
    category: number;
    id: number;
    datatype: DataType;
    name: string;
    data: number;
    raw: DataView;
    offset: number = 0;

    constructor(data: DataView) {
        this.category = data.getUint8(0);
        this.id = data.getUint16(1);
        this.datatype = data.getUint8(3);
        let nameLength: number = data.getUint16(4);
        // decode name
        this.name = '';
        this.raw = data;
        for (let i = 0; i < nameLength; i++) {
            this.name += String.fromCharCode(data.getUint8(6 + i));
        }
        this.offset = 6 + nameLength;
        switch (this.datatype) {
            case DataType.UINT8:
                this.data = data.getUint8(this.offset);
                this.offset++;
                break;
            case DataType.UINT32:
                this.data = data.getUint32(this.offset);
                this.offset++;
                break;
            case DataType.PERCENTAGE:
            case DataType.TEMPERATURE:
            case DataType.LOAD_USAGE:
                this.data = data.getFloat32(this.offset);
                this.offset++;
                break;
            case DataType.SPECIAL:
                this.data = 0;
                break;
            default:
                throw new Error('Unknown datatype');
        }
    }

    abstract update(): void;
    abstract render(): void;
    abstract renderOrUpdate(): void;
}