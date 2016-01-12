This is the spec for the websocket backend.

Endpoints:
* /ws/play:
  - Server to Client
    - {'type': 'turn'}
      This message is sent to client when it is their turn. From this point,
they can draw a card and must play a card.

    - {'type': 'start', 'hand': Card[], 'players': Player[]}
      This message is sent to all clients after the game is started by the
spectator. It contains the information about the player's hand and other
players' names in the game.

    - {'type': 'updateCard', 'card': Card}
      This message is sent to all clients after the card on the DiscardPile is
changed.

  - Client to Server
    - {'type': 'start', 'name': string}
      This message is used to register a new client with the server.

    - {'type': 'drawCard'}
      This message is sent when the player requests to draw a card.

    - {'type': 'playCard', 'card': Card}
      This message is sent when the player plays a card. Once a card is played,
if it can be played, the next player is given a turn.

* /ws/spectate:
  - Server to Client:
