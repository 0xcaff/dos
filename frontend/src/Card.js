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

const ReverseImage = (
<svg
  className='reverse strech'
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

const SkipImage = (
<svg
  width='10'
  height='10'
  viewBox='0 0 15 15'
  className='skip strech'>

  <circle
    cx='50%'
    cy='50%'
    r='40%'
    strokeWidth='2'
    stroke='inherit'
    fill='none' />
 
  <line
    x1='10%'
    y1='50%'
    x2='90%'
    y2='50%'
    strokeWidth='2'
    stroke='inherit'
    transform='rotate(-45 7.5 7.5)' />
</svg>
);

class Card extends Component {
  constructor(props) {
    super(props);

    if (props.onSwipe) {
      this.onSwipe = props.onSwipe.bind(this);
    }

    if (props.onBeforeSwipe) {
      this.onBeforeSwipe = props.onBeforeSwipe.bind(this);
    }
  }

  componentDidMount() {
    if (this.props.onSwipe) {
      this.element.addEventListener('slip:swipe', this.onSwipe);
    }

    if (this.props.onBeforeSwipe) {
      this.element.addEventListener('slip:beforeswipe', this.onBeforeSwipe);
    }
  }

  render() {
    if (!this.props.card) {
      return null;
    }

    let oval = null;
    if (this.props.card.type === dos.CardType.WILD) {
      oval = (
        <div className='wild oval'>
          <div className='yellow top left' />
          <div className='green top right' />
          <div className='blue bottom left' />
          <div className='red bottom right' />
        </div>
      );
    } else {
      oval = <div className='oval' />
    }

    let heroSymbol = null;
    let corners = null;
    let cornerClasses = ['corner'];
    if (this.props.card.type === dos.CardType.DOUBLEDRAW) {
      corners = '+2';
      heroSymbol = (
        <div className='strech inner-cards'>
          <div className='inner-card' />
          <div className='inner-card' />
        </div>
      );
    } else if (this.props.card.type === dos.CardType.QUADDRAW) {
      corners = '+4';
      heroSymbol = (
        <div className='strech quad-cards'>
          <div className='left red' />
          <div className='top blue' />
          <div className='right green' />
          <div className='bottom yellow' />
        </div>
      );
    } else if (this.props.card.type === dos.CardType.SKIP) {
      let skip = SkipImage;
      corners = skip;
      heroSymbol = <div className='hero-skip strech'>{ skip }</div>
    } else if (this.props.card.type === dos.CardType.REVERSE) {
      let reverse = ReverseImage;

      corners = reverse;
      heroSymbol = <div className='hero-reverse strech'>{ ReverseImage }</div>
    } else if (this.props.card.type === dos.CardType.NORMAL) {
      let number = this.props.card.number;
      corners = number;
      heroSymbol = <div className='number big'>{ number }</div>
    } else if (this.props.card.type === dos.CardType.WILD) {
      corners = oval;
      cornerClasses.push('wild');
    }

    return (
      <div ref={(element) => this.element = element}
           className={['card', NumbersToColors[this.props.card.color]].join(' ')} >
        <div className='strech background'>
          { oval }
          { heroSymbol }
          <div className={['left'].concat(cornerClasses).join(' ')}>{ corners }</div>
          <div className={['right'].concat(cornerClasses).join(' ')}>{ corners }</div>
        </div>
      </div>
    );
  }
}

export default Card;

