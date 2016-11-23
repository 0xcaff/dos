import {Serializable, prop} from './serializable';

export class Card extends Serializable {
  @prop color: string;
  @prop('typ') kind: string;
  @prop('number') digit: number;

  constructor(o: Object = {}) {
    super(o);
  }

  toString(): string {
    return `${this.kind} ${this.color} ${this.digit}`
  }
}

