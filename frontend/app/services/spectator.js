System.register(['angular2/core', './remote', '../models/card', '../models/player'], function(exports_1) {
    var __extends = (this && this.__extends) || function (d, b) {
        for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, remote_1, card_1, player_1;
    var SpectatorService, map;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (remote_1_1) {
                remote_1 = remote_1_1;
            },
            function (card_1_1) {
                card_1 = card_1_1;
            },
            function (player_1_1) {
                player_1 = player_1_1;
            }],
        execute: function() {
            SpectatorService = (function (_super) {
                __extends(SpectatorService, _super);
                function SpectatorService(url) {
                    _super.call(this, url);
                    this.deck = new Array();
                    this.players = new Array();
                }
                SpectatorService.prototype.receiveMessage = function (o) {
                    console.log(o);
                    switch (o['type']) {
                        case "init":
                            this.deck = o['deck'].map(function (elem) { return new card_1.Card(elem); });
                            this.players = o['players'].map(function (elem) { return new player_1.Player(elem); });
                            this.lastCard = new card_1.Card(o['lastCard']);
                            break;
                        case "start":
                            this.isStarted = true;
                            break;
                        case "update":
                            switch (o['for']) {
                                case "card":
                                    this.lastCard = new card_1.Card(o['what']);
                                    break;
                            }
                            break;
                        case "addition":
                        case "deletion":
                            switch (o['for']) {
                                case "player":
                                    var hand = this.players.filter(function (p) { return p.name == o['name']; })[0].hand;
                                    _super.prototype.processMessage.call(this, o, hand, card_1.Card);
                                    break;
                                default:
                                    var m = map[o['for']];
                                    _super.prototype.processMessage.call(this, o, this[m[0]], m[1]);
                                    break;
                            }
                            break;
                        case "end":
                            this.winner = o['winner'];
                            break;
                    }
                };
                SpectatorService.prototype.startGame = function () {
                    _super.prototype.sendRaw.call(this, "start");
                };
                SpectatorService = __decorate([
                    core_1.Injectable(), 
                    __metadata('design:paramtypes', [String])
                ], SpectatorService);
                return SpectatorService;
            })(remote_1.RemoteService);
            exports_1("SpectatorService", SpectatorService);
            map = {
                "deck": ["deck", card_1.Card],
                "players": ["players", player_1.Player],
            };
        }
    }
});
//# sourceMappingURL=spectator.js.map