System.register(['./serializable'], function(exports_1) {
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
    var serializable_1;
    var Card;
    return {
        setters:[
            function (serializable_1_1) {
                serializable_1 = serializable_1_1;
            }],
        execute: function() {
            Card = (function (_super) {
                __extends(Card, _super);
                function Card(o) {
                    if (o === void 0) { o = {}; }
                    _super.call(this, o);
                }
                Card.prototype.toString = function () {
                    return this.kind + " " + this.color + " " + this.digit;
                };
                __decorate([
                    serializable_1.prop, 
                    __metadata('design:type', String)
                ], Card.prototype, "color", void 0);
                __decorate([
                    serializable_1.prop('typ'), 
                    __metadata('design:type', String)
                ], Card.prototype, "kind", void 0);
                __decorate([
                    serializable_1.prop('number'), 
                    __metadata('design:type', Number)
                ], Card.prototype, "digit", void 0);
                return Card;
            })(serializable_1.Serializable);
            exports_1("Card", Card);
        }
    }
});
//# sourceMappingURL=card.js.map