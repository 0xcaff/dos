System.register([], function(exports_1) {
    var Serializable;
    function prop(one, two) {
        if (typeof one == 'string') {
            var serializedkey = one;
            return function (target, propertykey) {
                return serialize(target, serializedkey, propertykey);
            };
        }
        else if (typeof one == 'object') {
            var target = one;
            var propertykey = two;
            serialize(target, propertykey, propertykey);
        }
    }
    exports_1("prop", prop);
    function serialize(target, serializedkey, propertykey) {
        if (!target.hasOwnProperty('__serialization__'))
            target['__serialization__'] = {};
        target['__serialization__'][serializedkey] = propertykey;
    }
    function removeFirst(es, o) {
        for (var i = 0; i < es.length; i++) {
            var s = es[i];
            if (s.equals(o)) {
                es.splice(i, 1);
                break;
            }
        }
    }
    exports_1("removeFirst", removeFirst);
    return {
        setters:[],
        execute: function() {
            Serializable = (function () {
                function Serializable(o) {
                    this.unmarshall(o);
                }
                Serializable.prototype.unmarshall = function (o) {
                    if (o === void 0) { o = {}; }
                    var serial = this['__serialization__'];
                    var keys = Object.keys(serial);
                    for (var i = 0; i < keys.length; i++) {
                        var serializedKey = keys[i];
                        var prettyKey = serial[serializedKey];
                        var lo = o;
                        var keySegments = serializedKey.split('.');
                        for (var j = 0; j < keySegments.length && lo !== undefined; j++) {
                            lo = lo[keySegments[j]];
                        }
                        if (lo !== undefined)
                            this[prettyKey] = lo;
                    }
                };
                Serializable.prototype.marshall = function () {
                    var result = {};
                    var serial = this['__serialization__'];
                    var keys = Object.keys(serial);
                    for (var i = 0; i < keys.length; i++) {
                        var serializedKey = keys[i];
                        var prettyKey = serial[serializedKey];
                        var value = this[prettyKey];
                        if (value !== undefined)
                            recurse(serializedKey.split('.'), 0, result, value);
                    }
                    return result;
                    function recurse(keyPath, index, result, value) {
                        var path = keyPath[index];
                        if (index == keyPath.length - 1) {
                            result[path] = value;
                        }
                        else {
                            var rp = result[path];
                            if (!rp) {
                                rp = {};
                                result[path] = rp;
                            }
                            recurse(keyPath, index + 1, rp, value);
                        }
                    }
                };
                Serializable.prototype.equals = function (o) {
                    var keys = Object.keys(this['__serialization__']);
                    for (var i = 0; i < keys.length; i++) {
                        var prettyKey = this['__serialization__'][keys[i]];
                        var serializedKey = keys[i];
                        if (this.hasOwnProperty(prettyKey)
                            &&
                                o.hasOwnProperty(serializedKey)
                            &&
                                this[prettyKey] == o[serializedKey])
                            continue;
                        else
                            return false;
                    }
                    return true;
                };
                return Serializable;
            })();
            exports_1("Serializable", Serializable);
        }
    }
});
//# sourceMappingURL=serializable.js.map