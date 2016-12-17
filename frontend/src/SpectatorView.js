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

  render() {
    return (
      <div>
        <button onClick={this.props.startGame}>Start Game</button>

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
          </div>

          <div
            style={{
              'alignSelf': 'center',
              'width': '10em',
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
