import React, { Component } from 'react';
import './Cards.css';
import Card from './Card';
import Slip from 'slipjs';

class Cards extends Component {
  componentDidMount() {
    this.slip = new Slip(this.element);
  }

  componentWillUnmount() {
    this.slip.detach();
  }

  render() {
    return (
      <div ref={(element) => this.element = element} className='cards'>
        {this.props.cards.map(card =>
          <Card
            card={card}
            key={card.id} />
        )}
      </div>
    );
  }
}

export default Cards;
