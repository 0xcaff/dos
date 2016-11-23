System.register(['angular2/router', 'angular2/platform/browser', 'angular2/core', './app', './services/player', './services/spectator'], function(exports_1) {
    var router_1, browser_1, core_1, app_1, player_1, spectator_1;
    var baseURL;
    function getWebsocketURL() {
        var baseURL = '';
        if (window.location.protocol === "https:")
            baseURL += "wss:";
        else
            baseURL += "ws:";
        baseURL += "//" + window.location.host + window.location.pathname + "/";
        return baseURL;
    }
    return {
        setters:[
            function (router_1_1) {
                router_1 = router_1_1;
            },
            function (browser_1_1) {
                browser_1 = browser_1_1;
            },
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (app_1_1) {
                app_1 = app_1_1;
            },
            function (player_1_1) {
                player_1 = player_1_1;
            },
            function (spectator_1_1) {
                spectator_1 = spectator_1_1;
            }],
        execute: function() {
            baseURL = getWebsocketURL();
            browser_1.bootstrap(app_1.AppComponent, [
                router_1.ROUTER_PROVIDERS,
                core_1.provide(player_1.PlayerService, {
                    useFactory: function () { return new player_1.PlayerService(baseURL + "../ws/play"); },
                }),
                core_1.provide(spectator_1.SpectatorService, {
                    useFactory: function () { return new spectator_1.SpectatorService(baseURL + "../ws/spectate"); },
                }),
            ]);
        }
    }
});
//# sourceMappingURL=boot.js.map