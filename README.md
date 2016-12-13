dos
===

A multi-device, pluggable clone of [the card game UNO][uno].

Generate Protobuf Files
-----------------------

    protoc --go_out=. proto/*.proto
    cd frontend && yarn run protobuf

[uno]: https://en.wikipedia.org/wiki/Uno_(card_game)
