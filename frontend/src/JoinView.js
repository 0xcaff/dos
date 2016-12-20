import React, { Component } from 'react';
import './JoinView.css'

class JoinView extends Component {
  state = {
    name: '',
  }

  constructor(props) {
    super(props);

    // Send player handshake message
    // TODO: Bug in protobufjs causes empty messages to crash.
    props.socket.onopen = () => props.socket.send(new Uint8Array([]));

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleInput = this.handleInput.bind(this);
  }

  handleSubmit(event) {
    event.stopPropagation();
    event.preventDefault();

    this.props.setName(this.state.name);
  }

  handleInput(event) {
    this.setState({name: event.target.value});
  }

  render() {
    const disabled = this.props.connectionStatus !== 1;

    return (
      <div className='flex-center'>
        <div className='join-view'>
          <h1>Dos</h1>

          <form onSubmit={this.handleSubmit}>
            <input
              type='text'
              placeholder='Name'
              onChange={this.handleInput}
              disabled={disabled} />

            <button
              disabled={this.state.name.length === 0 || disabled}>
                Join Game
            </button>

          </form>

          { this.props.error && <div className='error'>
            { this.props.error }
          </div> }
        </div>
      </div>
    );
  }
}

export default JoinView;
