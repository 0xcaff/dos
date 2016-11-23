export function prop(serializedkey: string): PropertyDecorator;
export function prop(target: Object, propertykey: string): void;

export function prop(one: string|Object, two?: string): PropertyDecorator {
  if (typeof one == 'string') {
    var serializedkey: string = <string>one;

    return (target: Object, propertykey: string) =>
      serialize(target, serializedkey, propertykey);
  } else if (typeof one == 'object') {
    var target: Object = <Object>one;
    var propertykey: string = <string>two;

    serialize(target, propertykey, propertykey);
  }
}

function serialize(target: Object, serializedkey: string, propertykey: string): void {
  if (!target.hasOwnProperty('__serialization__'))
    target['__serialization__'] = {};

  target['__serialization__'][serializedkey] = propertykey;
}

export class Serializable {
  __serialization__: Object;

  constructor(o: Object) {
    this.unmarshall(o);
  }

  unmarshall(o: Object = {}) {
    var serial = this['__serialization__'];
    var keys = Object.keys(serial);

    for (var i = 0; i < keys.length; i++) {
      var serializedKey = keys[i];
      var prettyKey = serial[serializedKey];

      var lo: any = o;
      var keySegments = serializedKey.split('.');
      for (var j = 0; j < keySegments.length && lo !== undefined; j++) {
        lo = lo[keySegments[j]];
      }

      if (lo !== undefined)
        this[prettyKey] = lo;
    }
  }

  marshall(): Object {
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

    function recurse(keyPath: string[], index: number, result: Object, value: any) {
      var path = keyPath[index];
      if (index == keyPath.length - 1) {
        result[path] = value;
      } else {
        var rp = result[path];
        if (!rp) {
          rp = {};
          result[path] = rp;
        }
        recurse(keyPath, index + 1, rp, value);
      }
    }
  }

  equals(o: Object): boolean {
    var keys = Object.keys(this['__serialization__']);
    for (var i = 0; i < keys.length; i++) {
      var prettyKey = this['__serialization__'][keys[i]];
      var serializedKey = keys[i];

      if (
        this.hasOwnProperty(prettyKey)
          &&
        o.hasOwnProperty(serializedKey)
          &&
        this[prettyKey] == o[serializedKey]
      )
        continue
      else
        return false;
    }
    return true;
  }
}

export function removeFirst(es: Serializable[], o: Object): void {
  for (var i = 0; i < es.length; i++) {
    var s = es[i];
    if (s.equals(o)) {
      es.splice(i, 1)
      break;
    }
  }
}

