import React, { Component } from 'react';
import './JoinView.css'

class JoinView extends Component {
  constructor(props) {
    super(props);

    this.handleSubmit = this.handleSubmit.bind(this);
  }

  render() {
    return (
      <div className='flex-center'>
        <div className='join-view'>
          <h1>Dos</h1>

          <form onSubmit={this.handleSubmit}>
            <input
              type='text'
              placeholder='Name'
              ref={(input) => this.input = input} />

            <button>Join Game</button>
          </form>
        </div>
      </div>
    );
  }

  handleSubmit(event) {
    event.stopPropagation();
    event.preventDefault();

    this.props.setName(this.input.value);
  }
}

export default JoinView;
