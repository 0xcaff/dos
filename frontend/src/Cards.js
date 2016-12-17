import React, { Component } from 'react';
import Card from './Card';
import Slip from 'slipjs';

class Cards extends Component {
  constructor(props) {
    super(props);

    this.onSwipe = function(event) {
      const card = this.props.card;
      props.onSwipe(card);
    }
  }

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
            key={card.id}
            onSwipe={this.onSwipe} />
        )}
      </div>
    );
  }
}

export default Cards;
