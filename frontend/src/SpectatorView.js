import React, { Component } from 'react';
import { dos } from './proto';

class SpectatorView extends Component {
  constructor(props) {
    super(props);

    this.handleMessage = this.handleMessage.bind(this);
    this.startGame = this.startGame.bind(this);

    const socket = props.socket;

    const handshake = dos.HandshakeMessage.encode({
      type: dos.ClientType.SPECTATOR
    }).finish();
    socket.onopen = () => socket.send(handshake);

    socket.addEventListener('message', this.handleMessage);
  }

  handleMessage(event) {
	const data = new Uint8Array(event.data);
    const envelope = dos.Envelope.decode(data);

    // TODO: Handle messages
    // Turn Message
    // Amount of Cards In Each Players Hand
  }

  render() {
    return (<button onClick={this.startGame}>Start Game</button>);
  }

  startGame() {
    const message = dos.Envelope.encode({type: dos.MessageType.START}).finish();
    this.props.socket.send(message);
  }
}

export default SpectatorView;
