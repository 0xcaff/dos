import {Injectable} from 'angular2/core';
import {RemoteService} from './remote';
import {Card} from '../models/card';
import {Player} from '../models/player';

@Injectable()
export class SpectatorService extends RemoteService {
  deck: Card[] = new Array<Card>();
  players: Player[] = new Array<Player>();
  lastCard: Card;
  isStarted: boolean;
  winner: string;

  constructor(url: string) {
    super(url);
  }

  receiveMessage(o: Object) {
    console.log(o);

    switch(o['type']) {
      case "init":
        this.deck = o['deck'].map(elem => new Card(elem));
        this.players = o['players'].map(elem => new Player(elem));
        this.lastCard = new Card(o['lastCard']);
        break;

      case "start":
        this.isStarted = true;
        break;

      case "update":
        switch(o['for']) {
          case "card":
            this.lastCard = new Card(o['what'])
            break;
        }
        break;

      case "addition":
      case "deletion":
        switch (o['for']) {
          case "player":
            var hand = this.players.filter(p => p.name == o['name'])[0].hand;
            super.processMessage(o, hand, Card);
            break;

          default:
            var m = map[o['for']];
            super.processMessage(o, this[m[0]], m[1]);
            break;
        }
        break;

      case "end":
        this.winner = o['winner'];
        break;
    }
  }

  startGame() {
    super.sendRaw("start")
  }
}

var map = {
  "deck": ["deck", Card],
  "players": ["players", Player],
}

