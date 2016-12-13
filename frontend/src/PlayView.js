import React, { Component } from 'react';
import Cards from './Cards';

const players = [
  { name: 'Alice' },
  { name: 'Bob' },
  { name: 'Chris' },
  { name: 'Austin' },
  { name: 'Brandon' },
  { name: 'Zach' },
]

const cards = [
  { number: 1, type: 'normal', color: 'red' },
  { number: 1, type: 'normal', color: 'yellow' },
  { number: 1, type: 'normal', color: 'green' },
  { number: 1, type: 'normal', color: 'blue' },
  { type: 'wild', color: 'black' },
  { type: 'quaddraw', color: 'black' },
  { type: 'reverse', color: 'blue' },
];

// TODO: Last played card
class PlayView extends Component {
  render() {
    return (
      <div>
        <Cards cards={cards} />
        { players.map(player =>
            <div key={player.name}>{ player.name }</div> )}
      </div>
    );
  }
}

export default PlayView;
