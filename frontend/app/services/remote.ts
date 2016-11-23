import {Serializable, removeFirst} from '../models/serializable';

export class RemoteService {
  ws: WebSocket;

  constructor(url: string) {
    this.ws = new WebSocket(url)
    this.ws.addEventListener("message", m => this.receiveMessage(JSON.parse(m.data)))
  }

  sendMessage(typ?: string, msg: Object = {}) {
    if (typ)
      msg = Object.assign({'type': typ}, msg);
    this.ws.send(JSON.stringify(msg));
  }

  sendRaw(msg: string) {
    this.ws.send(msg);
  }

  receiveMessage(o: Object) { }

  processMessage(o: Object, fr: Serializable[], what: any) {
    switch (o['type']) {
      case "addition":
        fr.push(new what(o['what']));
        break;

      case "deletion":
        removeFirst(fr, o['what']);
        break;
    }
  }
}

function isEmpty(o: Object): boolean {
  for (var p in o)
    if (o.hasOwnProperty(p))
      return false;

  return true;
}

