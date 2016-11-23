import {Component} from 'angular2/core';
import {SpectatorService} from '../services/spectator';

@Component({
  templateUrl: 'templates/routes/spectate.html',
})
export class SpectateRoute {
  constructor(public ss: SpectatorService) { }
}


