import React, { Component } from 'react';
import './Card.css'

class Card extends Component {
  constructor(props) {
    super(props);

    this.onSwipe = this.onSwipe.bind(this);
  }

  componentDidMount() {
    this.element.addEventListener('slip:swipe', this.onSwipe);
  }

  onSwipe(event) {
    console.log(event);
  }

  render() {
    let oval = null;
    if (this.props.type === 'wild') {
      oval = (<div className='oval'>
        <div className='top left' />
        <div className='top right' />
        <div className='bottom left' />
        <div className='bottom right' />
      </div>);
    } else {
      oval = <div className='oval' />
    }

    // TODO: Handle +4, +2, Skip and Reverse
    let bigNumber = null;
    let smallLNum = null;
    let smallRNum = null;
    if (this.props.number !== undefined && this.props.number !== null) {
      bigNumber = <div className='number big'>{ this.props.number }</div>
      smallLNum = <div className='number small left'>{ this.props.number }</div>
      smallRNum = <div className='number small right'>{ this.props.number }</div>
    }

    return (
      <div ref={(element) => this.element = element}
           className={['card', this.props.color].join(' ')} >
        <div className='background'>
          { oval }
          { bigNumber }
          { smallLNum }
          { smallRNum }
        </div>
      </div>
    );
  }
}

export default Card;

