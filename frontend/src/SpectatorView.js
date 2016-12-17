import React, { Component } from 'react';
import Players from './Players';
import Card from './Players';
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
    return (
      <div>
        <button onClick={this.props.startGame}>Start Game</button>

        <Players
          players={this.props.players} />

        <Card
          card={this.props.discard} />

      </div>
    );
  }
}

export default SpectatorView;
