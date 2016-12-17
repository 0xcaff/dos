import React, { Component } from 'react';
import Cards from './Cards';
import Card from './Cards';
import Players from './Players';
import './PlayView.css';

// TODO: Draw/Done could be one button depending on the state. After drawing,
// you can't draw again and after playing, you must be done.
class PlayView extends Component {
  render() {
    return (
      <div className='play-view'>
        <div className='buttons'>
          <button onClick={this.props.drawCard}>Draw</button>
          <button onClick={this.props.turnDone}>Done</button>
        </div>

        <Cards
          cards={this.props.cards}
          onSwipe={this.props.playCard} />

        <Players
          players={this.props.players} />

        <Card
          card={this.props.discard} />
      </div>
    );
  }
}

export default PlayView;
