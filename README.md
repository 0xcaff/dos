dos
===

A multi-device clone of [the card game UNO][uno]. Built with React, Websockets
and Go.

Building
--------

Building from source requires [protoc][protoc], a javascript runtime, and the
[rice tool][rice].

    $ go get -d github.com/caffinatedmonkey/dos
    $ cd $GOPATH/src/github.com/caffinatedmonkey/dos
    $ go generate
    $ cd frontend
    $ yarn
    $ yarn run protobuf
    $ yarn run build
    $ cd ..
    $ rice embed-go
    $ go install

[uno]: https://en.wikipedia.org/wiki/Uno_(card_game)
[protoc]: https://github.com/google/protobuf/blob/master/src/README.md
[rice]: https://github.com/GeertJohan/go.rice
