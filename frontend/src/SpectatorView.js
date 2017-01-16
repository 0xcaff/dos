import React, { Component } from 'react';
import Players from './Players';
import Card from './Card';
import './SpectatorView.css'

class SpectatorView extends Component {
  constructor(props) {
    super(props);

    this.startGame = this.props.startGame.bind(this);
  }

  render() {
    return (
      <div className='spectator-view'>
        { !this.props.started &&
          <div className='session'>
            <h1>Game PIN</h1>
            <h1>{this.props.session}</h1>

            <p>
              Join the game by visiting <b>{ `${window.location.host}/join` }</b>.
            </p>
          </div>
        }

        <div className='players'>
          <Players
            players={this.props.players} />

          { !this.props.started &&
            <button
              onClick={this.props.startGame}
              disabled={
                this.props.connectionStatus !== 1 || this.props.players.length < 2
              }>Start Game!</button>
          }
        </div>

        { this.props.started &&
          <div className='discard'>
            <Card
              card={this.props.discard} />
          </div>
        }
      </div>
    );
  }
}

export default SpectatorView;
