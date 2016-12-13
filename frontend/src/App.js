import React, { Component } from 'react';
import JoinView from './JoinView';
import PlayView from './PlayView';
import LobbyView from './LobbyView';
import SpectatorView from './SpectatorView';
import { dos } from './proto';

class App extends Component {
  state = {
    view: (window.location.pathname.slice(1) || 'join'),
    players: [],
  }

  constructor(props) {
    super(props);

    this.setName = this.setName.bind(this);
    this.handleMessage = this.handleMessage.bind(this);

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

    const encoded = dos.ReadyMessage.encode({name: name}).finish();
	const envelope = dos.Envelope.encode({
      type: dos.MessageType.READY,
      contents: encoded,
    }).finish();
    this.socket.send(envelope);

    this.navigateTo('lobby');
  }

  handleMessage(event) {
	const data = new Uint8Array(event.data);
    const envelope = dos.Envelope.decode(data);

    if (envelope.type === dos.MessageType.PLAYERS) {
      const playersMessage = dos.PlayersMessage.decode(envelope.contents);

      if (playersMessage.initial.length > 0) {
        this.setState({
          players: playersMessage.initial,
        });
      } else if (playersMessage.addition) {
        this.setState({
          players: this.state.players.concat(playersMessage.addition),
        });
      } else if (playersMessage.deletion) {
        this.setState({
          players: this.state.players
            .filter(player => player !== playersMessage.deletion),
        });
      }
    }
  }

  render() {
    if (this.state.view === 'join') {
      return <JoinView
               setName={this.setName} />
    } else if (this.state.view === 'lobby') {
      return <LobbyView
               players={this.state.players} />
    } else if (this.state.view === 'play') {
      return <PlayView />
    } else if (this.state.view === 'spectate') {
      return <SpectatorView
               socket={this.socket} />
    }
  }

  navigateTo(destination) {
    this.setState({view: destination});
  }
}

export default App;
