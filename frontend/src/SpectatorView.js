import React, { Component } from 'react';
import { dos } from './proto';

class SpectatorView extends Component {
  constructor(props) {
    super(props);

    this.startGame = this.startGame.bind(this);
    const handshake = dos.HandshakeMessage.encode({
      type: dos.ClientType.SPECTATOR
    }).finish();
    this.props.socket.onopen = () => this.props.socket.send(handshake);
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
