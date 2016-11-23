import {Component} from 'angular2/core';
import {Router, RouteConfig, ROUTER_DIRECTIVES} from 'angular2/router';

import {SpectateRoute} from './routes/spectate';
import {PlayRoute} from './routes/play';

@Component({
  selector: 'dos-app',
  templateUrl: 'templates/index.html',
  directives: [ROUTER_DIRECTIVES],
})
@RouteConfig([
  {
    path: '/spectate',
    name: 'Spectate',
    component: SpectateRoute,
    useAsDefault: true,
  },
  {
    path: '/play',
    name: 'Play',
    component: PlayRoute,
  },
])

export class AppComponent {
  constructor() {}
}

