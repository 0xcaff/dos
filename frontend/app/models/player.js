System.register(['./serializable', './card'], function(exports_1) {
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
    var serializable_1, card_1;
    var Player;
    return {
        setters:[
            function (serializable_1_1) {
                serializable_1 = serializable_1_1;
            },
            function (card_1_1) {
                card_1 = card_1_1;
            }],
        execute: function() {
            Player = (function (_super) {
                __extends(Player, _super);
                function Player() {
                    _super.apply(this, arguments);
                }
                Player.prototype.unmarshall = function (o) {
                    _super.prototype.unmarshall.call(this, o);
                    if (o['hand']) {
                        this.hand = o['hand'].map(function (c) { return new card_1.Card(c); });
                    }
                };
                __decorate([
                    serializable_1.prop, 
                    __metadata('design:type', String)
                ], Player.prototype, "name", void 0);
                return Player;
            })(serializable_1.Serializable);
            exports_1("Player", Player);
        }
    }
});
//# sourceMappingURL=player.js.map