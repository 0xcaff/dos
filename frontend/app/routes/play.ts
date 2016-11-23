import {Component} from 'angular2/core';
import {PlayerService} from '../services/player';

@Component({
  templateUrl: 'templates/routes/play.html',
})
export class PlayRoute {
  constructor(public ps: PlayerService) {}

  ngOnInit() {
    var name = prompt("Please enter a name to play as");
    if (name) {
      this.ps.sendMessage('start', {name: name});
      this.ps.name = name;
    }
  }
}

