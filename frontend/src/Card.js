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

const ReverseImage =
(<svg
  className='reverse'
  height='3'
  width='5'
  viewBox='-2 0 6 6'
  xmlns='http://www.w3.org/2000/svg'>

  <path
    d='M0 3 h4 v1 l2 -2 l-2 -2 v1 h-2 q-2,0 -2,2' />

  <path
    d='M0 3 h4 v1 l2 -2 l-2 -2 v1 h-2 q-2,0 -2,2'
    transform='rotate(180, 3, 2) translate(4, -2)' />
</svg>);


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
    let corners = null;
    if (this.props.card.type === dos.CardType.DOUBLEDRAW) {
      corners = '+2';
      heroSymbol = (
        <div className='inner-cards'>
          <div className='inner-card' />
          <div className='inner-card' />
        </div>
      );
    } else if (this.props.card.type === dos.CardType.QUADDRAW) {
      corners = '+4';
      // TODO: Handle
    } else if (this.props.card.type === dos.CardType.SKIP) {
      // TODO: Impl
      // let skip = (
      // <div className='ban outer'>
      //   <div className='ban inner'></div>
      //   <div className='bar inner'></div>
      // </div>);

      // corners = skip;
      // heroSymbol = skip;
    } else if (this.props.card.type === dos.CardType.REVERSE) {
      let reverse = ReverseImage;

      corners = reverse;
      heroSymbol = <div className='hero-reverse'>{ReverseImage}</div>

    } else if (this.props.card.type === dos.CardType.NORMAL) {
      let number = this.props.card.number;
      corners = number;
      heroSymbol = <div className='number big'>{ number }</div>
    }

    let leftCorner = null;
    let rightCorner = null;
    if (this.props.card.type === dos.CardType.NORMAL) {
      leftCorner = <div className='corner number small left'>{ corners }</div>
      rightCorner = <div className='corner number small right'>{ corners }</div>
    } else {
      leftCorner = <div className='corner small left'>{ corners }</div>
      rightCorner = <div className='corner small right'>{ corners }</div>
    }

    return (
      <div ref={(element) => this.element = element}
           className={['card', NumbersToColors[this.props.card.color]].join(' ')} >
        <div className='background'>
          { oval }
          { heroSymbol }
          { leftCorner }
          { rightCorner }
        </div>
      </div>
    );
  }
}

export default Card;

