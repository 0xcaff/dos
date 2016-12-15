import React, { Component } from 'react';
import Cards from './Cards';

// TODO: Last played card
// TODO: Player list duplicated
class PlayView extends Component {
  render() {
    return (
      <div>
        <Cards cards={this.props.cards} />

        <h1>Players</h1>
        <ul>
          {this.props.players.map(player =>
            <li key={player}>{player}</li>)
          }
        </ul>
      </div>
    );
  }
}

export default PlayView;
