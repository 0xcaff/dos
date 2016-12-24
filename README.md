dos
===

A multi-device clone of [the card game UNO][uno]. Built with React, Websockets
and Go.

Demo
----
[![Thumbnail][thumb]][video]

Running
-------
Download the latest [release] and run it. There are no runtime dependencies.

    $ dos

Navigate to `/spectator` on the shared screen. The shared screen will show
player information and the last played card. It's the top screen in the demo.
Navigate to `/` on player devices. Once all players are ready, press `Start
Game` on the shared screen. The server needs to be restarted between games.

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
[release]: https://github.com/caffinatedmonkey/dos/releases
[video]: https://www.youtube.com/watch?v=0eZ_SirmF2c
[thumb]: https://caffinatedmonkey.github.io/dos/thumb.png
