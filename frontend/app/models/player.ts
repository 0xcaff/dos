import {Serializable, prop} from './serializable';
import {Card} from './card';

export class Player extends Serializable {
  hand: Card[];
  @prop name: string;

  unmarshall(o: Object) {
    super.unmarshall(o)

    if (o['hand']) {
      this.hand = o['hand'].map(c => new Card(c));
    }
  }
}

