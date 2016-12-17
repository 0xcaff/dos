import React, { Component } from 'react';
import Cards from './Cards';
import Players from './Players';
import './PlayView.css';

// TODO: Last played card
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
      </div>
    );
  }
}

export default PlayView;
