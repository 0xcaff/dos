import React, { Component } from 'react';
import JoinView from './JoinView';
import PlayView from './PlayView';
import SpectatorView from './SpectatorView';
import Players from './Players';
import SocketStatus from './SocketStatus';
import Card from './Card';
import { dos } from './proto';

// TODO: Report played cards right away
// TODO: SVGify Cards
// TODO: Handle black card color selection.
// TODO: Implement score board
// TODO: Change Done/Draw button size.
// TODO: Horizontal Card Layout like an acutal hand
class App extends Component {
  state = {
    view: (window.location.pathname.slice(1) || 'join'),
    players: [], // {name: string, active: boolean}
    cards: [],
    discard: null,
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

    // Open Connection
    this.socket = new WebSocket(`ws://drone.lan:8080/socket`);
    this.socket.binaryType = 'arraybuffer';
    window.onunload = () => this.socket.close();

    this.socket.addEventListener('message', this.handleMessage);
  }

  setName(name) {
    if (!name) {
      throw new Error("Invalid Name");
    }

    // Send handshake player handshake message
    // TODO: Bug in protobufjs causes empty messages to crash.
    this.socket.send(new Uint8Array([]));

    encodeAndSend(
      this.socket,
      dos.MessageType.READY,
      dos.ReadyMessage.encode({name: name}),
    );

    this.setState({
      name: name,
      view: 'lobby',
    });
  }

  startGame() {
    encodeAndSend(this.socket, dos.MessageType.START);
  }

  playCard(card) {
    encodeAndSend(
      this.socket,
      dos.MessageType.PLAY,
      dos.PlayMessage.encode({id: card.id}),
    );

    // TODO: errors from playing
    this.setState({hasPlayed: true});
  }

  // TODO: Handle one event per turn
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

      if (playersMessage.initial.length > 0) {
        this.setState({
          players: playersMessage.initial.map(
            name => ({name: name, active: false})
          ),
        });
      } else if (playersMessage.addition) {
        this.setState({
          players: this.state.players.concat([{
            name: playersMessage.addition,
            active: false
          }]),
        });
      } else if (playersMessage.deletion) {
        this.setState({
          players: this.state.players
            .filter(player => player.name !== playersMessage.deletion),
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
               setName={this.setName} />

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
               hasDrawn={this.state.hasDrawn}
               hasPlayed={this.state.hasPlayed}
               playCard={this.playCard} 
               drawCard={this.drawCard}
               turnDone={this.turnDone} />

    } else if (this.state.view === 'spectate') {
      view = <SpectatorView
               socket={this.socket}
               discard={this.state.discard}
               players={this.state.players}
               startGame={this.startGame} />

    } else if (this.state.view === 'scores') {
      // TODO: Implement
    } else {
      view = <Card
               card={this.state.discard} />
    }

    return (
      <div>
        <SocketStatus socket={this.socket} />
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
