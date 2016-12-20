import React, { Component } from 'react';
import Cards from './Cards';
import Card, { CanCover } from './Card';
import Players from './Players';
import './PlayView.css';

class PlayView extends Component {
  constructor(props) {
    super(props);

    this.playCard = this.playCard.bind(this);
    this.beforePlayCard = this.beforePlayCard.bind(this);
  }

  playCard(card, event) {
    if (CanCover(this.props.discard, card)) {
      this.props.playCard(card);
      this.props.turnDone();
    } else {
      event.preventDefault();
    }
  }

  beforePlayCard(event) {
    const active = this.props.players.find(player => player.active);
    let isActive = active && active.name === this.props.name;

    if (!isActive) {
      event.preventDefault();
    }
  }

  render() {
    const active = this.props.players.find(player => player.active);
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
          onSwipe={this.playCard}
          onBeforeSwipe={this.beforePlayCard} />

        <Players
          players={this.props.players} />

        <Card
          card={this.props.discard} />
      </div>
    );
  }
}

export default PlayView;
