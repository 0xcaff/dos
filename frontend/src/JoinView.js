import React, { Component } from 'react';

class JoinView extends Component {
  render() {
    return (
      <div>
        <input
          type='text'
          placeholder='Enter name'
          ref={(input) => this.input = input} />

        <button
          onClick={() => this.props.setName(this.input.value)}>
            Join Game
        </button>
      </div>
    );
  }
}

export default JoinView;
