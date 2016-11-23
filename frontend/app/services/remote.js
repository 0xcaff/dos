System.register(['../models/serializable'], function(exports_1) {
    var serializable_1;
    var RemoteService;
    function isEmpty(o) {
        for (var p in o)
            if (o.hasOwnProperty(p))
                return false;
        return true;
    }
    return {
        setters:[
            function (serializable_1_1) {
                serializable_1 = serializable_1_1;
            }],
        execute: function() {
            RemoteService = (function () {
                function RemoteService(url) {
                    var _this = this;
                    this.ws = new WebSocket(url);
                    this.ws.addEventListener("message", function (m) { return _this.receiveMessage(JSON.parse(m.data)); });
                }
                RemoteService.prototype.sendMessage = function (typ, msg) {
                    if (msg === void 0) { msg = {}; }
                    if (typ)
                        msg = Object.assign({ 'type': typ }, msg);
                    this.ws.send(JSON.stringify(msg));
                };
                RemoteService.prototype.sendRaw = function (msg) {
                    this.ws.send(msg);
                };
                RemoteService.prototype.receiveMessage = function (o) { };
                RemoteService.prototype.processMessage = function (o, fr, what) {
                    switch (o['type']) {
                        case "addition":
                            fr.push(new what(o['what']));
                            break;
                        case "deletion":
                            serializable_1.removeFirst(fr, o['what']);
                            break;
                    }
                };
                return RemoteService;
            })();
            exports_1("RemoteService", RemoteService);
        }
    }
});
//# sourceMappingURL=remote.js.map