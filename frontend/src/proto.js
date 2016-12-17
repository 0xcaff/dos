"use strict"; // eslint-disable-line strict

var $protobuf = require("protobufjs/runtime");

// Lazily resolved type references
var $lazyTypes = [];

// Exported root namespace
var $root = {};

/** @alias dos */
$root.dos = (function() {

    /** @alias dos */
    var dos = {};

    /** @alias dos.Card */
    dos.Card = (function() {

        /**
         * Constructs a new Card.
         * @exports dos.Card
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function Card(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.Card.prototype */
        var $prototype = Card.prototype;

        /**
         * Card id.
         * @name dos.Card#id
         * @type {number}
         */
        $prototype["id"] = 0;

        /**
         * Card number.
         * @name dos.Card#number
         * @type {number}
         */
        $prototype["number"] = 0;

        /**
         * Card type.
         * @name dos.Card#type
         * @type {number}
         */
        $prototype["type"] = 0;

        /**
         * Card color.
         * @name dos.Card#color
         * @type {number}
         */
        $prototype["color"] = 0;

        /**
         * Encodes the specified Card.
         * @function
         * @param {dos.Card|Object} message Card or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        Card.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,"dos.CardType","dos.CardColor"]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['id']!==undefined&&m['id']!==0)
                    w.tag(1,0).int32(m['id'])
                if(m['number']!==undefined&&m['number']!==0)
                    w.tag(2,0).int32(m['number'])
                if(m['type']!==undefined&&m['type']!==0)
                    w.tag(3,0).uint32(m['type'])
                if(m['color']!==undefined&&m['color']!==0)
                    w.tag(4,0).uint32(m['color'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified Card, length delimited.
         * @param {dos.Card|Object} message Card or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        Card.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a Card from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.Card} Card
         */
        Card.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,"dos.CardType","dos.CardColor"]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.Card
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['id']=r.int32()
                            break
                        case 2:
                            m['number']=r.int32()
                            break
                        case 3:
                            m['type']=r.uint32()
                            break
                        case 4:
                            m['color']=r.uint32()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a Card from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.Card} Card
         */
        Card.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a Card.
         * @param {dos.Card|Object} message Card or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        Card.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,"dos.CardType","dos.CardColor"]);
            return function verify(m) {
                if(m['id']!==undefined){
                    if(!util.isInteger(m['id']))
                        return"invalid value for field .dos.Card.id (integer expected)"
                }
                if(m['number']!==undefined){
                    if(!util.isInteger(m['number']))
                        return"invalid value for field .dos.Card.number (integer expected)"
                }
                if(m['type']!==undefined){
                    switch(m['type']){
                        default:
                            return"invalid value for field .dos.Card.type (enum value expected)"
                        case 0:
                        case 1:
                        case 2:
                        case 3:
                        case 4:
                        case 5:
                            break
                    }
                }
                if(m['color']!==undefined){
                    switch(m['color']){
                        default:
                            return"invalid value for field .dos.Card.color (enum value expected)"
                        case 0:
                        case 1:
                        case 2:
                        case 3:
                        case 4:
                            break
                    }
                }
                return null
            }
            /* eslint-enable */
        })();

        return Card;
    })();

    /**
     * CardType values.
     * @exports dos.CardType
     * @type {Object.<string,number>}
     */
    dos.CardType = {

        NORMAL: 0,
        SKIP: 1,
        DOUBLEDRAW: 2,
        REVERSE: 3,
        WILD: 4,
        QUADDRAW: 5
    };

    /**
     * CardColor values.
     * @exports dos.CardColor
     * @type {Object.<string,number>}
     */
    dos.CardColor = {

        RED: 0,
        YELLOW: 1,
        GREEN: 2,
        BLUE: 3,
        BLACK: 4
    };

    /** @alias dos.CardsChangedMessage */
    dos.CardsChangedMessage = (function() {

        /**
         * Constructs a new CardsChangedMessage.
         * @exports dos.CardsChangedMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function CardsChangedMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.CardsChangedMessage.prototype */
        var $prototype = CardsChangedMessage.prototype;

        /**
         * CardsChangedMessage additions.
         * @name dos.CardsChangedMessage#additions
         * @type {Array.<dos.Card>}
         */
        $prototype["additions"] = $protobuf.util.emptyArray;

        /**
         * CardsChangedMessage deletions.
         * @name dos.CardsChangedMessage#deletions
         * @type {Array.<number>}
         */
        $prototype["deletions"] = $protobuf.util.emptyArray;

        /**
         * Encodes the specified CardsChangedMessage.
         * @function
         * @param {dos.CardsChangedMessage|Object} message CardsChangedMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        CardsChangedMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.Card",null]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['additions'])
                    for(var i=0;i<m['additions'].length;++i)
                    types[0].encode(m['additions'][i],w.tag(1,2).fork()).ldelim()
                if(m['deletions']&&m['deletions'].length){
                    w.fork()
                    for(var i=0;i<m['deletions'].length;++i)
                        w.int32(m['deletions'][i])
                    w.ldelim(2)
                }
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified CardsChangedMessage, length delimited.
         * @param {dos.CardsChangedMessage|Object} message CardsChangedMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        CardsChangedMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a CardsChangedMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.CardsChangedMessage} CardsChangedMessage
         */
        CardsChangedMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.Card",null]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.CardsChangedMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['additions']&&m['additions'].length?m['additions']:m['additions']=[]
                            m['additions'][m['additions'].length]=types[0].decode(r,r.uint32())
                            break
                        case 2:
                            m['deletions']&&m['deletions'].length?m['deletions']:m['deletions']=[]
                            if(t.wireType===2){
                                var e=r.uint32()+r.pos
                                while(r.pos<e)
                                    m['deletions'][m['deletions'].length]=r.int32()
                            }else
                                m['deletions'][m['deletions'].length]=r.int32()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a CardsChangedMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.CardsChangedMessage} CardsChangedMessage
         */
        CardsChangedMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a CardsChangedMessage.
         * @param {dos.CardsChangedMessage|Object} message CardsChangedMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        CardsChangedMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.Card",null]);
            return function verify(m) {
                if(m['additions']!==undefined){
                    if(!Array.isArray(m['additions']))
                        return"invalid value for field .dos.CardsChangedMessage.additions (array expected)"
                    for(var i=0;i<m['additions'].length;++i){
                        var r;
                        if(r=types[0].verify(m['additions'][i]))
                            return r
                    }
                }
                if(m['deletions']!==undefined){
                    if(!Array.isArray(m['deletions']))
                        return"invalid value for field .dos.CardsChangedMessage.deletions (array expected)"
                    for(var i=0;i<m['deletions'].length;++i){
                        if(!util.isInteger(m['deletions'][i]))
                            return"invalid value for field .dos.CardsChangedMessage.deletions (integer[] expected)"
                    }
                }
                return null
            }
            /* eslint-enable */
        })();

        return CardsChangedMessage;
    })();

    /** @alias dos.PlayMessage */
    dos.PlayMessage = (function() {

        /**
         * Constructs a new PlayMessage.
         * @exports dos.PlayMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function PlayMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.PlayMessage.prototype */
        var $prototype = PlayMessage.prototype;

        /**
         * PlayMessage id.
         * @name dos.PlayMessage#id
         * @type {number}
         */
        $prototype["id"] = 0;

        /**
         * PlayMessage color.
         * @name dos.PlayMessage#color
         * @type {number}
         */
        $prototype["color"] = 0;

        /**
         * Encodes the specified PlayMessage.
         * @function
         * @param {dos.PlayMessage|Object} message PlayMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        PlayMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.CardColor"]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['id']!==undefined&&m['id']!==0)
                    w.tag(2,0).int32(m['id'])
                if(m['color']!==undefined&&m['color']!==0)
                    w.tag(3,0).uint32(m['color'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified PlayMessage, length delimited.
         * @param {dos.PlayMessage|Object} message PlayMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        PlayMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a PlayMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.PlayMessage} PlayMessage
         */
        PlayMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.CardColor"]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.PlayMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 2:
                            m['id']=r.int32()
                            break
                        case 3:
                            m['color']=r.uint32()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a PlayMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.PlayMessage} PlayMessage
         */
        PlayMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a PlayMessage.
         * @param {dos.PlayMessage|Object} message PlayMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        PlayMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.CardColor"]);
            return function verify(m) {
                if(m['id']!==undefined){
                    if(!util.isInteger(m['id']))
                        return"invalid value for field .dos.PlayMessage.id (integer expected)"
                }
                if(m['color']!==undefined){
                    switch(m['color']){
                        default:
                            return"invalid value for field .dos.PlayMessage.color (enum value expected)"
                        case 0:
                        case 1:
                        case 2:
                        case 3:
                        case 4:
                            break
                    }
                }
                return null
            }
            /* eslint-enable */
        })();

        return PlayMessage;
    })();

    /** @alias dos.Envelope */
    dos.Envelope = (function() {

        /**
         * Constructs a new Envelope.
         * @exports dos.Envelope
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function Envelope(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.Envelope.prototype */
        var $prototype = Envelope.prototype;

        /**
         * Envelope type.
         * @name dos.Envelope#type
         * @type {number}
         */
        $prototype["type"] = 0;

        /**
         * Envelope contents.
         * @name dos.Envelope#contents
         * @type {Uint8Array}
         */
        $prototype["contents"] = $protobuf.util.emptyArray;

        /**
         * Encodes the specified Envelope.
         * @function
         * @param {dos.Envelope|Object} message Envelope or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        Envelope.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.MessageType",null]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['type']!==undefined&&m['type']!==0)
                    w.tag(1,0).uint32(m['type'])
                if(m['contents']!==undefined&&m['contents']!==[])
                    w.tag(2,2).bytes(m['contents'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified Envelope, length delimited.
         * @param {dos.Envelope|Object} message Envelope or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        Envelope.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a Envelope from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.Envelope} Envelope
         */
        Envelope.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.MessageType",null]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.Envelope
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['type']=r.uint32()
                            break
                        case 2:
                            m['contents']=r.bytes()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a Envelope from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.Envelope} Envelope
         */
        Envelope.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a Envelope.
         * @param {dos.Envelope|Object} message Envelope or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        Envelope.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.MessageType",null]);
            return function verify(m) {
                if(m['type']!==undefined){
                    switch(m['type']){
                        default:
                            return"invalid value for field .dos.Envelope.type (enum value expected)"
                        case 0:
                        case 1:
                        case 2:
                        case 3:
                        case 4:
                        case 5:
                        case 6:
                        case 7:
                            break
                    }
                }
                if(m['contents']!==undefined){
                    if(!(m['contents']&&typeof m['contents'].length==='number'||util.isString(m['contents'])))
                        return"invalid value for field .dos.Envelope.contents (buffer expected)"
                }
                return null
            }
            /* eslint-enable */
        })();

        return Envelope;
    })();

    /**
     * MessageType values.
     * @exports dos.MessageType
     * @type {Object.<string,number>}
     */
    dos.MessageType = {

        TURN: 0,
        PLAYERS: 1,
        CARDS: 2,
        DRAW: 3,
        PLAY: 4,
        DONE: 5,
        READY: 6,
        START: 7
    };

    /** @alias dos.TurnMessage */
    dos.TurnMessage = (function() {

        /**
         * Constructs a new TurnMessage.
         * @exports dos.TurnMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function TurnMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.TurnMessage.prototype */
        var $prototype = TurnMessage.prototype;

        /**
         * TurnMessage player.
         * @name dos.TurnMessage#player
         * @type {string}
         */
        $prototype["player"] = "";

        /**
         * TurnMessage lastPlayed.
         * @name dos.TurnMessage#lastPlayed
         * @type {dos.Card}
         */
        $prototype["lastPlayed"] = null;

        /**
         * Encodes the specified TurnMessage.
         * @function
         * @param {dos.TurnMessage|Object} message TurnMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        TurnMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.Card"]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['player']!==undefined&&m['player']!=="")
                    w.tag(1,2).string(m['player'])
                if(m['lastPlayed']!==undefined&&m['lastPlayed']!==null)
                    types[1].encode(m['lastPlayed'],w.fork()).len&&w.ldelim(2)||w.reset()
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified TurnMessage, length delimited.
         * @param {dos.TurnMessage|Object} message TurnMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        TurnMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a TurnMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.TurnMessage} TurnMessage
         */
        TurnMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.Card"]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.TurnMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['player']=r.string()
                            break
                        case 2:
                            m['lastPlayed']=types[1].decode(r,r.uint32())
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a TurnMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.TurnMessage} TurnMessage
         */
        TurnMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a TurnMessage.
         * @param {dos.TurnMessage|Object} message TurnMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        TurnMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,"dos.Card"]);
            return function verify(m) {
                if(m['player']!==undefined){
                    if(!util.isString(m['player']))
                        return"invalid value for field .dos.TurnMessage.player (string expected)"
                }
                if(m['lastPlayed']!==undefined){
                    var r;
                    if(r=types[1].verify(m['lastPlayed']))
                        return r
                }
                return null
            }
            /* eslint-enable */
        })();

        return TurnMessage;
    })();

    /** @alias dos.HandshakeMessage */
    dos.HandshakeMessage = (function() {

        /**
         * Constructs a new HandshakeMessage.
         * @exports dos.HandshakeMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function HandshakeMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.HandshakeMessage.prototype */
        var $prototype = HandshakeMessage.prototype;

        /**
         * HandshakeMessage type.
         * @name dos.HandshakeMessage#type
         * @type {number}
         */
        $prototype["type"] = 0;

        /**
         * Encodes the specified HandshakeMessage.
         * @function
         * @param {dos.HandshakeMessage|Object} message HandshakeMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        HandshakeMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.ClientType"]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['type']!==undefined&&m['type']!==0)
                    w.tag(1,0).uint32(m['type'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified HandshakeMessage, length delimited.
         * @param {dos.HandshakeMessage|Object} message HandshakeMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        HandshakeMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a HandshakeMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.HandshakeMessage} HandshakeMessage
         */
        HandshakeMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.ClientType"]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.HandshakeMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['type']=r.uint32()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a HandshakeMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.HandshakeMessage} HandshakeMessage
         */
        HandshakeMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a HandshakeMessage.
         * @param {dos.HandshakeMessage|Object} message HandshakeMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        HandshakeMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = ["dos.ClientType"]);
            return function verify(m) {
                if(m['type']!==undefined){
                    switch(m['type']){
                        default:
                            return"invalid value for field .dos.HandshakeMessage.type (enum value expected)"
                        case 0:
                        case 1:
                            break
                    }
                }
                return null
            }
            /* eslint-enable */
        })();

        return HandshakeMessage;
    })();

    /**
     * ClientType values.
     * @exports dos.ClientType
     * @type {Object.<string,number>}
     */
    dos.ClientType = {

        PLAYER: 0,
        SPECTATOR: 1
    };

    /** @alias dos.ReadyMessage */
    dos.ReadyMessage = (function() {

        /**
         * Constructs a new ReadyMessage.
         * @exports dos.ReadyMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function ReadyMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.ReadyMessage.prototype */
        var $prototype = ReadyMessage.prototype;

        /**
         * ReadyMessage name.
         * @name dos.ReadyMessage#name
         * @type {string}
         */
        $prototype["name"] = "";

        /**
         * Encodes the specified ReadyMessage.
         * @function
         * @param {dos.ReadyMessage|Object} message ReadyMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        ReadyMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['name']!==undefined&&m['name']!=="")
                    w.tag(1,2).string(m['name'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified ReadyMessage, length delimited.
         * @param {dos.ReadyMessage|Object} message ReadyMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        ReadyMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ReadyMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.ReadyMessage} ReadyMessage
         */
        ReadyMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.ReadyMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['name']=r.string()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a ReadyMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.ReadyMessage} ReadyMessage
         */
        ReadyMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a ReadyMessage.
         * @param {dos.ReadyMessage|Object} message ReadyMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        ReadyMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null]);
            return function verify(m) {
                if(m['name']!==undefined){
                    if(!util.isString(m['name']))
                        return"invalid value for field .dos.ReadyMessage.name (string expected)"
                }
                return null
            }
            /* eslint-enable */
        })();

        return ReadyMessage;
    })();

    /** @alias dos.PlayersMessage */
    dos.PlayersMessage = (function() {

        /**
         * Constructs a new PlayersMessage.
         * @exports dos.PlayersMessage
         * @constructor
         * @param {Object} [properties] Properties to set
         */
        function PlayersMessage(properties) {
            if (properties) {
                var keys = Object.keys(properties);
                for (var i = 0; i < keys.length; ++i)
                    this[keys[i]] = properties[keys[i]];
            }
        }

        /** @alias dos.PlayersMessage.prototype */
        var $prototype = PlayersMessage.prototype;

        /**
         * PlayersMessage initial.
         * @name dos.PlayersMessage#initial
         * @type {Array.<string>}
         */
        $prototype["initial"] = $protobuf.util.emptyArray;

        /**
         * PlayersMessage addition.
         * @name dos.PlayersMessage#addition
         * @type {string}
         */
        $prototype["addition"] = "";

        /**
         * PlayersMessage deletion.
         * @name dos.PlayersMessage#deletion
         * @type {string}
         */
        $prototype["deletion"] = "";

        /**
         * Encodes the specified PlayersMessage.
         * @function
         * @param {dos.PlayersMessage|Object} message PlayersMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        PlayersMessage.encode = (function() {
            /* eslint-disable */
            var Writer = $protobuf.Writer;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,null]);
            return function encode(m, w) {
                w||(w=Writer.create())
                if(m['initial'])
                    for(var i=0;i<m['initial'].length;++i)
                    w.tag(1,2).string(m['initial'][i])
                if(m['addition']!==undefined&&m['addition']!=="")
                    w.tag(2,2).string(m['addition'])
                if(m['deletion']!==undefined&&m['deletion']!=="")
                    w.tag(3,2).string(m['deletion'])
                return w
            }
            /* eslint-enable */
        })();

        /**
         * Encodes the specified PlayersMessage, length delimited.
         * @param {dos.PlayersMessage|Object} message PlayersMessage or plain object to encode
         * @param {Writer} [writer] Writer to encode to
         * @returns {Writer} Writer
         */
        PlayersMessage.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a PlayersMessage from the specified reader or buffer.
         * @function
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {dos.PlayersMessage} PlayersMessage
         */
        PlayersMessage.decode = (function() {
            /* eslint-disable */
            var Reader = $protobuf.Reader;
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,null]);
            return function decode(r, l) {
                r instanceof Reader||(r=Reader.create(r))
                var c=l===undefined?r.len:r.pos+l,m=new $root.dos.PlayersMessage
                while(r.pos<c){
                    var t=r.tag()
                    switch(t.id){
                        case 1:
                            m['initial']&&m['initial'].length?m['initial']:m['initial']=[]
                            m['initial'][m['initial'].length]=r.string()
                            break
                        case 2:
                            m['addition']=r.string()
                            break
                        case 3:
                            m['deletion']=r.string()
                            break
                        default:
                            r.skipType(t.wireType)
                            break
                    }
                }
                return m
            }
            /* eslint-enable */
        })();

        /**
         * Decodes a PlayersMessage from the specified reader or buffer, length delimited.
         * @param {Reader|Uint8Array} readerOrBuffer Reader or buffer to decode from
         * @returns {dos.PlayersMessage} PlayersMessage
         */
        PlayersMessage.decodeDelimited = function decodeDelimited(readerOrBuffer) {
            readerOrBuffer = readerOrBuffer instanceof $protobuf.Reader ? readerOrBuffer : $protobuf.Reader(readerOrBuffer);
            return this.decode(readerOrBuffer, readerOrBuffer.uint32());
        };

        /**
         * Verifies a PlayersMessage.
         * @param {dos.PlayersMessage|Object} message PlayersMessage or plain object to verify
         * @returns {?string} `null` if valid, otherwise the reason why it is not
         */
        PlayersMessage.verify = (function() {
            /* eslint-disable */
            var util = $protobuf.util;
            var types; $lazyTypes.push(types = [null,null,null]);
            return function verify(m) {
                if(m['initial']!==undefined){
                    if(!Array.isArray(m['initial']))
                        return"invalid value for field .dos.PlayersMessage.initial (array expected)"
                    for(var i=0;i<m['initial'].length;++i){
                        if(!util.isString(m['initial'][i]))
                            return"invalid value for field .dos.PlayersMessage.initial (string[] expected)"
                    }
                }
                if(m['addition']!==undefined){
                    if(!util.isString(m['addition']))
                        return"invalid value for field .dos.PlayersMessage.addition (string expected)"
                }
                if(m['deletion']!==undefined){
                    if(!util.isString(m['deletion']))
                        return"invalid value for field .dos.PlayersMessage.deletion (string expected)"
                }
                return null
            }
            /* eslint-enable */
        })();

        return PlayersMessage;
    })();

    return dos;
})();

// Resolve lazy types
$lazyTypes.forEach(function(types) {
    types.forEach(function(path, i) {
        if (!path)
            return;
        path = path.split('.');
        var ptr = $root;
        while (path.length)
            ptr = ptr[path.shift()];
        types[i] = ptr;
    });
});

$protobuf.roots = {};
$protobuf.roots["default"] = $root;

module.exports = $root;
