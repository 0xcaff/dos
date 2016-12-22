import React, { Component } from 'react';
import JoinView from './JoinView';
import PlayView from './PlayView';
import SpectatorView from './SpectatorView';
import Players from './Players';
import SocketStatus from './SocketStatus';
import { dos } from './proto';

var WEBSOCKET_PATH = null;
if (process.env.NODE_ENV === 'production') {
  const isSecure = window.location.protocol === 'https:';
  WEBSOCKET_PATH = `${isSecure ? 'wss' : 'ws'}://${window.location.host}/socket`;
} else {
  WEBSOCKET_PATH = 'ws://drone.lan:8080';
}

// TODO: Implement score board
class App extends Component {
  state = {
    view: (window.location.pathname.slice(1) || 'join'),
    players: [], // {name: string, active: boolean}
    cards: [],
    discard: null,
    error: '',
    name: '',
    hasDrawn: false,
    hasPlayed: false,
  }

  constructor(props) {
    super(props);

    this.setName = this.setName.bind(this);
    this.startGame = this.startGame.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.playCard = this.playCard.bind(this);
    this.drawCard = this.drawCard.bind(this);
    this.turnDone = this.turnDone.bind(this);
    this.handleSocketChange = this.handleSocketChange.bind(this);

    // Open Connection
    this.socket = new WebSocket(WEBSOCKET_PATH);
    this.socket.binaryType = 'arraybuffer';
    window.onunload = () => this.socket.close();

    this.socket.addEventListener('message', this.handleMessage);
    this.socket.addEventListener('close', this.handleSocketChange);
    this.socket.addEventListener('open', this.handleSocketChange);

    // eslint-disable-next-line react/no-direct-mutation-state
    this.state.connectionStatus = this.socket.readyState;
  }

  handleSocketChange() {
    this.setState({
      connectionStatus: this.socket.readyState,
    });
  }

  setName(name) {
    encodeAndSend(
      this.socket,
      dos.MessageType.READY,
      dos.ReadyMessage.encode({name: name}),
    );

    const messageHandler = (event) => {
      const data = new Uint8Array(event.data);
      const envelope = dos.Envelope.decode(data);
      this.socket.removeEventListener('message', messageHandler);

      if (envelope.type === dos.MessageType.SUCCESS) {
        this.setState({
          name: name,
          view: 'lobby',
        });
      } else if (envelope.type === dos.MessageType.ERROR) {
        const errorMessage = dos.ErrorMessage.decode(envelope.contents);

        this.setState({
          error: errorMessage.reason,
        });
      }
    }

    this.socket.addEventListener('message', messageHandler);
  }

  startGame() {
    encodeAndSend(this.socket, dos.MessageType.START);
  }

  playCard(card, wildColor) {
    encodeAndSend(
      this.socket,
      dos.MessageType.PLAY,
      dos.PlayMessage.encode({id: card.id, color: wildColor}),
    );
    this.setState({hasPlayed: true});
  }

  drawCard() {
    encodeAndSend(this.socket, dos.MessageType.DRAW);
    this.setState({hasDrawn: true});
  }

  turnDone() {
    encodeAndSend(this.socket, dos.MessageType.DONE);
  }

  handleMessage(event) {
	const data = new Uint8Array(event.data);
    const envelope = dos.Envelope.decode(data);

    if (envelope.type === dos.MessageType.PLAYERS) {
      const playersMessage = dos.PlayersMessage.decode(envelope.contents);
      console.log(playersMessage);

      if (playersMessage.additions.length > 0) {
        this.setState({
          players: this.state.players.concat(
            playersMessage.additions.map(name => ({
              name: name,
              active: false,
            }))
          ),
        });
      }

      if (playersMessage.deletions.length > 0) {
        const deletions = new Set(playersMessage.deletions);
        this.setState({
          players: this.state.players
            .filter(player => !deletions.has(player.name)),
        });
      }

    } else if (envelope.type === dos.MessageType.CARDS) {
      const cardsMessage = dos.CardsChangedMessage.decode(envelope.contents);
      console.log(cardsMessage);

      if (cardsMessage.additions.length > 0) {
        this.setState({
          cards: this.state.cards.concat(cardsMessage.additions),
        });
      }

      if (cardsMessage.deletions.length > 0) {
        const deletions = new Set(cardsMessage.deletions);
        this.setState({
          cards: this.state.cards.filter(card => !deletions.has(card.id)),
        });
      }
      this.navigateTo('play');
    } else if (envelope.type === dos.MessageType.TURN) {
      const turnMessage = dos.TurnMessage.decode(envelope.contents);
      console.log(turnMessage);

      this.setState({
        players: this.state.players.map(player => {
          if (player.name === turnMessage.player) {
            player.active = true;
            if (player.name === this.state.name) {
              this.setState({isActive: true});
            }
          } else {
            player.active = false;
          }

          return player;
        }),
        discard: turnMessage.lastPlayed,
        hasDrawn: false,
        hasPlayed: false,
      });
    }
  }

  render() {
    let view;
    if (this.state.view === 'join') {
      view = <JoinView
               socket={this.socket}
               connectionStatus={this.state.connectionStatus}
               setName={this.setName}
               error={this.state.error} />

    } else if (this.state.view === 'lobby') {
      view = (<div className='flex-center'>
        <Players
          players={this.state.players} />
      </div>);

    } else if (this.state.view === 'play') {
      view = <PlayView
               cards={this.state.cards}
               players={this.state.players}
               name={this.state.name}
               discard={this.state.discard}
               playCard={this.playCard} 
               drawCard={this.drawCard}
               turnDone={this.turnDone}
               hasDrawn={this.state.hasDrawn}
               hasPlayed={this.state.hasPlayed} />

    } else if (this.state.view === 'spectate') {
      view = <SpectatorView
               socket={this.socket}
               discard={this.state.discard}
               players={this.state.players}
               startGame={this.startGame}
               connectionStatus={this.state.connectionStatus} />

    } else if (this.state.view === 'scores') {
      // TODO: Implement
    }

    return (
      <div>
        <SocketStatus readyState={this.state.connectionStatus} />
        { view }
      </div>
    );
  }

  navigateTo(destination) {
    this.setState({view: destination});
  }
}

function encodeAndSend(socket, type, message) {
  let encoded = new Uint8Array([]);

  if (message) {
    encoded = message.finish();
  }

  const envelope = dos.Envelope.encode({
    type: type,
    contents: encoded,
  }).finish();

  socket.send(envelope);
}

export default App;
