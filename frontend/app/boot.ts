import {ROUTER_PROVIDERS} from 'angular2/router';
import {bootstrap} from 'angular2/platform/browser';
import {provide} from 'angular2/core';

import {AppComponent} from './app';
import {PlayerService} from './services/player';
import {SpectatorService} from './services/spectator';

var baseURL = getWebsocketURL();
bootstrap(AppComponent, [
  ROUTER_PROVIDERS,
  provide(PlayerService, {
    useFactory:
      () => new PlayerService(baseURL + "../ws/play"),
  }),
  provide(SpectatorService, {
    useFactory:
      () => new SpectatorService(baseURL + "../ws/spectate"),
  }),
]);

function getWebsocketURL(): string {
  var baseURL: string = '';
  if (window.location.protocol === "https:")
    baseURL += "wss:";
  else
    baseURL += "ws:";
  baseURL += "//" + window.location.host + window.location.pathname + "/";
  return baseURL;
}

