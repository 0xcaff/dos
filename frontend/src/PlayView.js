import React, { Component } from 'react';
import Cards from './Cards';
import Card, { CanCover } from './Card';
import Players from './Players';
import { dos } from './proto';
import './PlayView.css';

// TODO: PositionFixed looks wierd up stuff in ff mobile
// TODO: Remove padding which hides cards at edges.

class PlayView extends Component {
  state = {
    seekingColor: false,
    card: null,
  };

  constructor(props) {
    super(props);

    this.playCard = this.playCard.bind(this);
    this.playWildCard = this.playWildCard.bind(this);
    this.beforePlayCard = this.beforePlayCard.bind(this);
  }

  playWildCard(color) {
    this.setState({seekingColor: false});
    this.props.playCard(this.state.card, color);
    this.props.turnDone();
  }

  playCard(card, event) {
    if (card.color === dos.CardColor.BLACK) {
      this.setState({seekingColor: true, card: card});
      return
    }

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
        <Cards
          cards={this.props.cards}
          onSwipe={this.playCard}
          onBeforeSwipe={this.beforePlayCard} />

        <Players
          players={this.props.players} />

        <Card
          card={this.props.discard} />

        <div className='actions'>
          <div className='container'>
            { button }

            <div className='status-dots'>
              {this.props.players.map(player => (
                <span
                  key={player.name}
                  title={player.name}
                  className={[
                    player.active ? 'active' : undefined,
                    player.name === this.props.name ? 'me' : undefined,
                 ].filter(className => !!className).join(' ')} />
              ))}
            </div>
          </div>
        </div>

        {this.state.seekingColor && <div className='color-picker strech'>
          <div
            onClick={() => this.playWildCard(dos.CardColor.YELLOW)}
            className='top left yellow' />

          <div
            onClick={() => this.playWildCard(dos.CardColor.GREEN)}
            className='top right green' />

          <div
            onClick={() => this.playWildCard(dos.CardColor.BLUE)}
            className='bottom left blue' />

          <div
            onClick={() => this.playWildCard(dos.CardColor.RED)}
            className='bottom right red' />
        </div>}
      </div>
    );
  }
}

export default PlayView;
