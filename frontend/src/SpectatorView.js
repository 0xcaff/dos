import React, { Component } from 'react';
import Players from './Players';
import Card from './Card';
import { dos } from './proto';

class SpectatorView extends Component {
  constructor(props) {
    super(props);

    this.startGame = this.props.startGame.bind(this);

    const handshake = dos.HandshakeMessage.encode({
      type: dos.ClientType.SPECTATOR
    }).finish();
    this.props.socket.onopen = () => this.props.socket.send(handshake);
  }

  // TODO: Make this look nice.
  render() {
    return (
      <div>
        <div
          style={{
            'display': 'flex',
            'justifyContent': 'space-around',
          }}>
          <div
            style={{
              'alignSelf': 'flex-start',
            }}>
            <Players
              players={this.props.players} />

            <button
              onClick={this.props.startGame}
              disabled={this.props.connectionStatus !== 1}>Start Game</button>
          </div>

          <div
            style={{
              'alignSelf': 'center',
              'width': '50vmin',
            }}>
            <Card
              card={this.props.discard} />
          </div>
        </div>
      </div>
    );
  }
}

export default SpectatorView;
