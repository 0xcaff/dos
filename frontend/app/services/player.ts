import {Injectable} from 'angular2/core';
import {RemoteService} from './remote';

import {Card} from '../models/card';

@Injectable()
export class PlayerService extends RemoteService {
  hand: Card[] = new Array<Card>();
  lastCard: Card;
  name: string;
  currentActivePlayer: string;
  winner: string;

  constructor(url: string) {
    super(url);
  }

  playCard(i: number) {
    super.sendMessage("play", {'card': i});
  }

  drawCard() {
    super.sendMessage('draw');
  }

  receiveMessage(o: Object) {
    console.log(o);
    switch (o['type']) {
      case 'init':
        this.hand = o['hand'].map(e => new Card(e));
        break;

      case 'update':
        if (o['for'] == 'turn') {
          // if (this.lastCard)
          //   this.lastCard.unmarshall(o['what']);
          // else
          this.lastCard = new Card(o['what']);
          this.currentActivePlayer = o['active'];
        }
        break;

      case "addition":
      case "deletion":
        var m = map[o['for']];
        super.processMessage(o, this[m[0]], m[1]);
        break;

      case "error":
        console.error(o["message"]);
        break;

      case "end":
        this.winner = o['winner'];
        break;
    }
  }

  get isActive(): boolean {
    return this.name == this.currentActivePlayer;
  }

  get isWinner(): boolean {
    return this.name == this.winner;
  }
}

var map = {
  "hand": ["hand", Card],
  "turn": ["lastCard"],
}

