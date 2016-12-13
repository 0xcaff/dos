import React, { Component } from 'react';

class LobbyView extends Component {
  render() {
    return (
      <ul>
        { this.props.players.map(player =>
            <li key={player}>{ player }</li> )}
      </ul>
    );
  }
}

export default LobbyView;
