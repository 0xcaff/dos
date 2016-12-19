import React, { Component } from 'react';
import Cards from './Cards';
import Card from './Card';
import Players from './Players';
import './PlayView.css';

class PlayView extends Component {
  render() {
    let active = this.props.players.find(player => player.active);
    let isActive = active && active.name === this.props.name;

    let button = null;
    if (this.props.hasDrawn || this.props.hasPlayed) {
      button = <button
                 onClick={this.props.turnDone}
                 disabled={!isActive}>Done</button>
    } else {
      button = <button
                 onClick={this.props.drawCard}
                 disabled={!isActive}>Draw</button>
    }

    return (
      <div className='play-view'>
        { button }
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
