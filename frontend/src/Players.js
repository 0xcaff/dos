import React, { Component } from 'react';
import './Players.css'

// TODO: Some way to show which direction the game is going in.
// TODO: Player Scores
class Players extends Component {
  render() {
    return (
      <div>
        <h1>Players</h1>
        <ul>
          {this.props.players.map(player =>
            <li
              key={player.name}
              className={player.active && 'active'}>
                {player.name}
            </li>
          )}
        </ul>
      </div>
    );
  }
}

export default Players;

