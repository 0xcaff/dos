import React, { Component } from 'react';
import InputView from './InputView';
import PlayView from './PlayView';
import SpectatorView from './SpectatorView';
import Players from './Players';
import SocketStatus from './SocketStatus';
import { dos } from './proto';

const isSecure = window.location.protocol === 'https:';
const WEBSOCKET_PATH = `${isSecure ? 'wss' : 'ws'}://${window.location.host}/socket`;

class App extends Component {
  state = {
    players: [], // {name: string, active: boolean}
    cards: [],
    discard: null,
    error: '',
    name: '',
    hasDrawn: false,
    hasPlayed: false,
    started: false,
  }

  constructor(props) {
    super(props);

    // Primitive Routing
    let path = window.location.pathname.slice(1);
    if (!['join', 'spectate'].includes(path)) {
      path = 'join';
      history.replaceState(null, null, "/join");
    }

    // eslint-disable-next-line react/no-direct-mutation-state
    this.state.view = path;

    this.setName = this.setName.bind(this);
    this.setSession = this.setSession.bind(this);
    this.startGame = this.startGame.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.playCard = this.playCard.bind(this);
    this.drawCard = this.drawCard.bind(this);
    this.turnDone = this.turnDone.bind(this);
    this.handleSocketChange = this.handleSocketChange.bind(this);
    this.handleClose = this.handleClose.bind(this);

    // Open Connection
    this.socket = new WebSocket(WEBSOCKET_PATH);
    this.socket.binaryType = 'arraybuffer';
    window.onunload = () => this.socket.close();

    // Send handshake message
    this.socket.onopen = () => {
      if (this.state.view === 'join') {
        // TODO: Bug in protobufjs causes empty messages to crash.
        this.socket.send(new Uint8Array([]));
      } else if (this.state.view === 'spectate') {
        const handshake = dos.HandshakeMessage.encode({
          type: dos.ClientType.SPECTATOR
        }).finish();
        this.socket.send(handshake);
      }
    }

    this.socket.addEventListener('message', this.handleMessage);
    this.socket.addEventListener('close', this.handleSocketChange);
    this.socket.addEventListener('close', this.handleClose);
    this.socket.addEventListener('open', this.handleSocketChange);

    // eslint-disable-next-line react/no-direct-mutation-state
    this.state.connectionStatus = this.socket.readyState;
  }

  handleClose(event) {
    console.log(event);
    if (event.code === 1000 && event.reason === "won!") {
      this.navigateTo('done');
    }
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
          error: '',
        });
      }
    }

    this.socket.addEventListener('message', messageHandler);
  }

  setSession(session) {
    encodeAndSend(
      this.socket,
      dos.MessageType.SESSION,
      dos.SessionMessage.encode({session: session}),
    );

    const messageHandler = (event) => {
      const data = new Uint8Array(event.data);
      const envelope = dos.Envelope.decode(data);
      this.socket.removeEventListener('message', messageHandler);

      if (envelope.type === dos.MessageType.SUCCESS) {
        this.setState({
          session: session,
          error: '',
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
        started: true,
      });
    } else if (envelope.type === dos.MessageType.SESSION) {
      const sessionMessage = dos.SessionMessage.decode(envelope.contents);
      console.log(sessionMessage);

      this.setState({session: sessionMessage.session});
    } else if (envelope.type === dos.MessageType.ERROR) {
      const errorMessage = dos.ErrorMessage.decode(envelope.contents);

      this.setState({error: convertError(errorMessage.reason)});
    }
  }

  render() {
    let view;
    if (this.state.view === 'join') {
      const disconnected = this.state.connectionStatus !== 1;

      if (!this.state.session) {
        view = <InputView
          // Required to prevent reconciliation.
          // https://facebook.github.io/react/docs/reconciliation.html
          key='session'
          placeholder='Game PIN'
          onSubmit={this.setSession}
          error={this.state.error}
          type='tel'
          disabled={disconnected}>
            <p style={{fontSize: 'x-small'}}>
              Get a Game PIN by visiting <b>{ `${window.location.host}/spectate` }</b> on a shared screen.
            </p>
          </InputView>
      } else {
        view = <InputView
          key='name'
          placeholder='Name'
          onSubmit={this.setName}
          error={this.state.error}
          type='text'
          disabled={disconnected} />
      }
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
               discard={this.state.discard}
               players={this.state.players}
               session={this.state.session}
               started={this.state.started}
               startGame={this.startGame}
               connectionStatus={this.state.connectionStatus} />

    } else if (this.state.view === 'done') {
      view = (<div className='flex-center'>
        <h1>Game Done</h1>
      </div>);
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

function convertError(error) {
  if (error === dos.ErrorReason.INVALIDGAME) {
    return "That game doesn't exist";
  } else if (error === dos.ErrorReason.INVALIDNAME) {
    return "That name is already taken";
  } else if (error === dos.ErrorReaspon.GAMESTARTED) {
    return "That game is already started.";
  } else {
    return '';
  }
}

export default App;
