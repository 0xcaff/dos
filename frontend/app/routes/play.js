System.register(['angular2/core', '../services/player'], function(exports_1) {
    var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
        var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
        if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
        else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
        return c > 3 && r && Object.defineProperty(target, key, r), r;
    };
    var __metadata = (this && this.__metadata) || function (k, v) {
        if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
    };
    var core_1, player_1;
    var PlayRoute;
    return {
        setters:[
            function (core_1_1) {
                core_1 = core_1_1;
            },
            function (player_1_1) {
                player_1 = player_1_1;
            }],
        execute: function() {
            PlayRoute = (function () {
                function PlayRoute(ps) {
                    this.ps = ps;
                }
                PlayRoute.prototype.ngOnInit = function () {
                    var name = prompt("Please enter a name to play as");
                    if (name) {
                        this.ps.sendMessage('start', { name: name });
                        this.ps.name = name;
                    }
                };
                PlayRoute = __decorate([
                    core_1.Component({
                        templateUrl: 'templates/routes/play.html',
                    }), 
                    __metadata('design:paramtypes', [player_1.PlayerService])
                ], PlayRoute);
                return PlayRoute;
            })();
            exports_1("PlayRoute", PlayRoute);
        }
    }
});
//# sourceMappingURL=play.js.map