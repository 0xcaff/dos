System.register(['angular2/core', './remote', '../models/card'], function(exports_1) {
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
    var core_1, remote_1, card_1;
    var PlayerService, map;
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
            }],
        execute: function() {
            PlayerService = (function (_super) {
                __extends(PlayerService, _super);
                function PlayerService(url) {
                    _super.call(this, url);
                    this.hand = new Array();
                }
                PlayerService.prototype.playCard = function (i) {
                    _super.prototype.sendMessage.call(this, "play", { 'card': i });
                };
                PlayerService.prototype.drawCard = function () {
                    _super.prototype.sendMessage.call(this, 'draw');
                };
                PlayerService.prototype.receiveMessage = function (o) {
                    console.log(o);
                    switch (o['type']) {
                        case 'init':
                            this.hand = o['hand'].map(function (e) { return new card_1.Card(e); });
                            break;
                        case 'update':
                            if (o['for'] == 'turn') {
                                // if (this.lastCard)
                                //   this.lastCard.unmarshall(o['what']);
                                // else
                                this.lastCard = new card_1.Card(o['what']);
                                this.currentActivePlayer = o['active'];
                            }
                            break;
                        case "addition":
                        case "deletion":
                            var m = map[o['for']];
                            _super.prototype.processMessage.call(this, o, this[m[0]], m[1]);
                            break;
                        case "error":
                            console.error(o["message"]);
                            break;
                        case "end":
                            this.winner = o['winner'];
                            break;
                    }
                };
                Object.defineProperty(PlayerService.prototype, "isActive", {
                    get: function () {
                        return this.name == this.currentActivePlayer;
                    },
                    enumerable: true,
                    configurable: true
                });
                Object.defineProperty(PlayerService.prototype, "isWinner", {
                    get: function () {
                        return this.name == this.winner;
                    },
                    enumerable: true,
                    configurable: true
                });
                PlayerService = __decorate([
                    core_1.Injectable(), 
                    __metadata('design:paramtypes', [String])
                ], PlayerService);
                return PlayerService;
            })(remote_1.RemoteService);
            exports_1("PlayerService", PlayerService);
            map = {
                "hand": ["hand", card_1.Card],
                "turn": ["lastCard"],
            };
        }
    }
});
//# sourceMappingURL=player.js.map