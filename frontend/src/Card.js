import React, { Component } from 'react';
import './Card.css'
import { dos } from './proto';

const NumbersToColors = ((object) => {
  let newObject = {};

  for (let key in object) {
    if (object.hasOwnProperty(key)) {
      let value = object[key];
      key = key.toLowerCase();

      newObject[value] = key;
    }
  }
  return newObject;

})(dos.CardColor);

class Card extends Component {
  constructor(props) {
    super(props);

    if (props.onSwipe) {
      this.onSwipe = props.onSwipe.bind(this);
    }
  }

  componentDidMount() {
    if (this.props.onSwipe) {
      this.element.addEventListener('slip:swipe', this.onSwipe);
    }
  }

  render() {
    if (!this.props.card) {
      return null;
    }

    let oval = null;
    if (this.props.card.type === dos.CardType.WILD) {
      oval = (
        <div className='oval'>
          <div className='top left' />
          <div className='top right' />
          <div className='bottom left' />
          <div className='bottom right' />
        </div>
      );
    } else {
      oval = <div className='oval' />
    }

    let heroSymbol = null;
    if (this.props.card.type === dos.CardType.DOUBLEDRAW) {
      heroSymbol = (
        <div className='inner-cards'>
          <div className='inner-card' />
          <div className='inner-card' />
        </div>
      );
    } else if (this.props.card.type === dos.CardType.QUADDRAW) {
      // TODO: Handle
    } else if (this.props.card.type === dos.CardType.SKIP) {
      // TODO: Handle
      // TODO: Handle corners
    } else if (this.props.card.type === dos.CardType.REVERSE) {
      // TODO: Handle
      // TODO: Handle Corners
    }

    let bigNumber = null;
    let smallLNum = null;
    let smallRNum = null;

    if (this.props.card.type === dos.CardType.QUADDRAW) {
      let number = '+4';
      smallLNum = <div className='number small left'>{ number }</div>
      smallRNum = <div className='number small right'>{ number }</div>
    } else if (this.props.card.type === dos.CardType.DOUBLEDRAW) {
      let number = '+2';
      smallLNum = <div className='number small left'>{ number }</div>
      smallRNum = <div className='number small right'>{ number }</div>
    } else if (this.props.card.number !== -1) {
      let number = this.props.card.number;
      bigNumber = <div className='number big'>{ number }</div>
      smallLNum = <div className='number small left'>{ number }</div>
      smallRNum = <div className='number small right'>{ number }</div>
    }

    return (
      <div ref={(element) => this.element = element}
           className={['card', NumbersToColors[this.props.card.color]].join(' ')} >
        <div className='background'>
          { oval }
          { heroSymbol }

          { bigNumber }
          { smallLNum }
          { smallRNum }
        </div>
      </div>
    );
  }
}

export default Card;

