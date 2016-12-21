import React, { Component } from 'react';
import Players from './Players';
import Card from './Card';
import { dos } from './proto';
import './SpectatorView.css'

class SpectatorView extends Component {
  constructor(props) {
    super(props);

    this.startGame = this.props.startGame.bind(this);

    const handshake = dos.HandshakeMessage.encode({
      type: dos.ClientType.SPECTATOR
    }).finish();
    this.props.socket.onopen = () => this.props.socket.send(handshake);
  }

  render() {
    return (
      <div className='spectator-view'>
        <div className='players'>
          <Players
            players={this.props.players} />

          <button
            onClick={this.props.startGame}
            disabled={this.props.connectionStatus !== 1}>Start Game!</button>
        </div>

        <div className='discard'>
          <Card
            card={this.props.discard} />
        </div>
      </div>
    );
  }
}

export default SpectatorView;
